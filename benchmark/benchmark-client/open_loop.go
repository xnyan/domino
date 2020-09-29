package main

import (
	"fmt"
	"sync"
	"time"

	"domino/benchmark/workload"
)

var wg sync.WaitGroup

func olExecute(client Client) {
	if TxnTargetRate <= 0 {
		logger.Fatalf("Target rate is %d in open-loop settings.", TxnTargetRate)
	}
	interval := time.Duration(int64(time.Second) / int64(TxnTargetRate))
	n := int(int64(Duration)/int64(time.Second))*int(TxnTargetRate) + 100
	if n < TxnTotalNum {
		n = TxnTotalNum
	}
	txnStatList := make([]*TxnStat, n)

	switch client.(type) {
	case *EpaxosClient:
		go client.(*EpaxosClient).HandleReplies()
	}

	c := 0 // txn count
	s := time.Now()
	d := time.Since(s)
	for d < Duration || (Duration <= 0 && c < TxnTotalNum) {
		txn := Workload.GenTxn()
		txnStatList[c] = NewTxnStat(txn.TxnId, txn.ReadKeys, txn.WriteData)

		wg.Add(1)
		go olExecuteTxn(client, txn, txnStatList[c])

		time.Sleep(interval)
		d = time.Since(s)
		c++
	}

	wg.Wait()

	for i := 0; i < c; i++ {
		olPrintTxnExecStat(txnStatList[i], s)
	}
}

func olExecuteTxn(client Client, txn *workload.Txn, txnStat *TxnStat) {
	startTime := time.Now()
	isUseFp, isAccept, isFast, _ := client.SyncExecTxn(txn.ReadKeys, txn.WriteData)
	//isUseFp, isAccept, isFast, val := client.SyncExecTxn(txn.ReadKeys, txn.WriteData)
	//logger.Infof("Result txnId = %s %t %t %t %s", txn.TxnId, isUseFp, isAccept, isFast, val)
	endTime := time.Now()
	execStat := &TxnExecutionStat{
		"0",
		&startTime,
		&endTime,
		isUseFp,
		isAccept,
		isFast,
	}
	txnStat.ExecHistory = append(txnStat.ExecHistory, execStat)

	wg.Done()
}

func olPrintTxnExecStat(
	txnStat *TxnStat,
	benchmarkStartTime time.Time,
) time.Duration {
	firstExecStat := txnStat.ExecHistory[0]

	// Transaction completion time from the first time starting the transaction until
	// the time the transaction is accepted or rejected after retrying the max number of times.
	latency := firstExecStat.EndTime.Sub(*(firstExecStat.BeginTime))

	// txnId, isCommit, latency, abort num, 1st-run start time
	fmt.Println(
		//txnStat.TxnStatId, ",", // txn id (may not be the txnId in Carousel)
		bool2int[firstExecStat.IsUseFp], ",", // true : use Fast Paxos, otherwise false
		bool2int[firstExecStat.IsAccept], ",", // true: accepted, false: rejected
		bool2int[firstExecStat.IsFast], ",", // true: fast-path decision, false: slow-path decision
		int64(latency), ",", // unit: ns
		//len(txnStat.ExecHistory)-1, ",", // # of aborted
		// elapsed time (ns) that the txn first executes from the benchmark starts
		int64(firstExecStat.BeginTime.Sub(benchmarkStartTime)),
		//",",
		//firstExecStat.BeginTime.UnixNano(), // time in ns the transaction first starts
	)

	return latency
}
