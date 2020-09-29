package main

import (
	"fmt"
	"github.com/op/go-logging"
	"math"
	"math/rand"
	"strconv"
	"time"

	"domino/benchmark/config"
	"domino/benchmark/workload"
	//"domino/clientlib"
	"domino/common"
)

var logger = logging.MustGetLogger("benchmark-client")
var bool2int = map[bool]int{true: 1, false: 0}

func main() {
	// Parses commad-line args
	ParseArgs()

	// Debug mode
	common.ConfigLogger(IsDebug)

	// Loads benchmark configurations
	LoadBenchmarkConfig()

	// Inits a client instance

	var client Client

	switch ClientType {
	//case config.CLIENT_DO:
	//	c := clientlib.NewClient(ClientId, DcId, ConfigFile)
	//	if IsExec {
	//		client = &DoExecClient{FpClient: c}
	//	} else {
	//		client = &DoCommitClient{FpClient: c}
	//	}
	//case config.CLIENT_HYBRID:
	//	client = NewHybridClient(ClientId, DcId, ConfigFile)
	case config.CLIENT_FP:
		client = NewFastPaxosClient(ClientId, DcId, ConfigFile, ReplicaFile)
	case config.CLIENT_DYNAMIC:
		dc := NewDynamicClient(ClientId, DcId, ConfigFile, ReplicaFile, TargetDcId)
		client = dc
	case config.CLIENT_EPAXOS:
		client = NewEpaxosClient(DcId, TargetDcId, ConfigFile, ReplicaFile, true, false)
	case config.CLIENT_MENCIUS:
		client = NewEpaxosClient(DcId, TargetDcId, ConfigFile, ReplicaFile, true, false)
	case config.CLIENT_PAXOS:
		client = NewEpaxosClient(DcId, TargetDcId, ConfigFile, ReplicaFile, false, false)
	//case config.CLIENT_GPAXOS:
	//	client = NewEpaxosClient(DcId, TargetDcId, ConfigFile, ReplicaFile, false, true)
	default:
		logger.Fatalf("Invalid client type = %s", ClientType)
	}

	if IsOpenLoop {
		olExecute(client)
	} else {
		// Close loop
		execute(client)
	}

	client.Close()
}

// Executes the txns in the give list in one by one (close-loop)
func execute(client Client) {

	txnStatTable := make(map[string]*TxnStat) // mappings from txnId to its txn stat
	var execCount int64

	var durationPerTxn time.Duration // expected execution duration per transaction
	var timeLeg time.Duration        // time leg to achieve the target rate
	if TxnTargetRate > 0 {
		durationPerTxn = time.Duration(int64(time.Second) / int64(TxnTargetRate))
		timeLeg = time.Duration(0)
	}

	printTxnExecStatTag() // print comments for understanding the output statistics

	txn := Workload.GenTxn() // first transaction
	txnCount := 1

	benchmarkStartTime := time.Now() // benchmark starting time
	benchmarkExecutionTime := time.Since(benchmarkStartTime)

	// Starts transaction execution
	for benchmarkExecutionTime < Duration || (Duration <= 0 && txnCount <= TxnTotalNum) {
		if _, exists := txnStatTable[txn.TxnId]; !exists {
			txnStatTable[txn.TxnId] = NewTxnStat(txn.TxnId, txn.ReadKeys, txn.WriteData)
			execCount = 0
		}

		// Executes the transaction
		execStat := executeTxn(client, txn, strconv.FormatInt(execCount, 10))
		execCount++

		// Records the execution statistics
		txnStat := txnStatTable[txn.TxnId]
		txnStat.ExecHistory = append(txnStat.ExecHistory, execStat)

		if !execStat.IsAccept {
			if isRetry, waitTime := isRetryTxn(execCount); isRetry {
				// Retries the transaction after waiting for a certain amount of time

				logger.Debug(
					"txnStatId = ", txnStat.TxnStatId, // may not be the txnId in Carousel's execution
					"rejected. Execution time = ", execStat.ExecutionDuration(),
					"Waiting for", waitTime, "to retry the transaction.")

				time.Sleep(waitTime) // Wait for a while to retry the transaction

				logger.Debug("txnStatId = ", txnStat.TxnStatId, " starts to retry.")

				continue
			}
		}

		// Outputs txn execution information
		latency := printTxnExecStat(execStat, txnStat, benchmarkStartTime)

		// Releases memory
		delete(txnStatTable, txn.TxnId)

		if TxnTargetRate > 0 {
			// Tries to keep sending the target number of transactions per second
			timeLeg = tryToMaintainTxnTargetRate(latency, durationPerTxn, timeLeg)

		}

		benchmarkExecutionTime = time.Since(benchmarkStartTime)
		txn = Workload.GenTxn() // next transaction
		txnCount++
	}
}

func tryToMaintainTxnTargetRate(latency, durationPerTxn, timeLeg time.Duration) time.Duration {
	if latency < durationPerTxn {
		expectWait := durationPerTxn - latency

		if timeLeg <= expectWait {
			expectWait -= timeLeg
			timeLeg = 0
		} else {
			timeLeg -= expectWait
			expectWait = 0
		}

		if expectWait > 0 {
			time.Sleep(expectWait)
		}
	} else {
		leg := latency - durationPerTxn
		if tmp := timeLeg + leg; tmp > timeLeg {
			// Overflow
			timeLeg = tmp
		}
	}

	return timeLeg
}

// Prints the transaction execution statistics
func printTxnExecStat(
	execStat *TxnExecutionStat,
	txnStat *TxnStat,
	benchmarkStartTime time.Time,
) time.Duration {
	firstExecStat := txnStat.ExecHistory[0]

	// Transaction completion time from the first time starting the transaction until
	// the time the transaction is accepted or rejected after retrying the max number of times.
	latency := execStat.EndTime.Sub(*(firstExecStat.BeginTime))

	// txnId, isCommit, latency, abort num, 1st-run start time
	fmt.Println(
		//txnStat.TxnStatId, ",", // request id
		bool2int[execStat.IsUseFp], ",", // true: use Fast Paxos, otherwise false
		bool2int[execStat.IsAccept], ",", // true: accepted, false: rejected
		bool2int[execStat.IsFast], ",", // true: fast-path decision, false: slow-path decision
		int64(latency), ",", // unit: ns
		//len(txnStat.ExecHistory)-1, ",", // # of retried
		// elapsed time (ns) that the txn first executes from the benchmark starts
		int64(firstExecStat.BeginTime.Sub(benchmarkStartTime)),
		//",",
		//firstExecStat.BeginTime.UnixNano(), // time in ns the transaction first starts
	)

	return latency
}

func printTxnExecStatTag() {
	fmt.Println("## isUseFastPaxos, isAccept, isFast, latency(ns), " +
		"elapsed time (ns) since start")
}

// Returns true and waiting time if retrying the txn
func isRetryTxn(execNum int64) (bool, time.Duration) {
	if RetryMax >= 0 && execNum > RetryMax {
		// Exceeded the max number of retries, do not retry any more
		return false, 0
	}

	waitTime := RetryInterval

	if RetryMode == config.RETRY_MODE_EXP {
		//exponential backoff
		abortNum := execNum
		randomFactor := rand.Int63n(int64(math.Exp2(float64(abortNum))))
		if randomFactor > RetrySlotMax {
			randomFactor = RetrySlotMax
		}
		waitTime = RetryInterval * time.Duration(randomFactor)
	}
	return true, time.Duration(waitTime)
}

// Executes one txn
func executeTxn(
	client Client,
	txn *workload.Txn,
	txnExecId string,
) *TxnExecutionStat {
	startTime := time.Now()
	isUseFp, isAccept, isFast, _ := client.ExecTxn(txn.ReadKeys, txn.WriteData)
	//logger.Infof("TxnExecId = %s, isAccept = %t, isFast = %t, ret = %s", txnExecId, isAccept, isFast, ret)
	endTime := time.Now()
	execStat := &TxnExecutionStat{
		txnExecId,
		&startTime,
		&endTime,
		isUseFp,
		isAccept,
		isFast,
	}
	return execStat
}
