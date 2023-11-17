package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/op/go-logging"

	"domino/common"
	"domino/dynamic/dynamic"
	"domino/kv"
)

var logger = logging.MustGetLogger("server")

// Configurations
var isDebug bool = false
var replicaId string = ""
var configFile string = ""
var replicaFile string = "" // replica location file

// Network
var NetAddr string

// Lists of replica Ids and addresses, both are ordered by Ids
var ReplicaIdList []string = make([]string, 0)
var NodeAddrList []string = make([]string, 0)
var ReplicaIdAddrMap map[string]string = make(map[string]string)

// Consensus
var PaxosShard int32 = -1 // Paxos shard id
var FpShard int32 = -1
var PaxosShardLeaderMap map[int32]string
var FpShardLeaderMap map[int32]string
var FollowerAddrList []string = make([]string, 0)
var HbInterval time.Duration = time.Duration(10 * 1000 * 1000) // Heart beat interval, unit: ns
var IsExec bool = false
var IsExecReply bool = false
var IsExecLog bool = false
var CmdChSize int = 1024 * 1024       // 1MB
var ExecChSize int = 10 * 1024 * 1024 // 10MB
var IsFpLeaderUsePaxos bool = true
var IsFpLeaderSoleLearner bool = true // If the FP leader is the only learner
var IsGrpc bool = true
var IsSyncSend bool = false
var AppLogDir string = "./"

// Latency prediction
var IsLatPrediction bool = true
var ProbeInv time.Duration
var WindowLen time.Duration
var WindowSize int
var LatPredictPth float64

// Paxos using future time
var IsPaxosFutureTime bool

// Application
var DataKeyFile string
var DataValSize int = 8 // unit: bytes

func main() {
	parseArgs()
	// Configures logging
	common.ConfigLogger(isDebug)
	// Parses the configuration file
	loadConfig(configFile, replicaFile)

	logger.Infof("Starting Dynamic Paxos Consensus")
	// Inits consensus member
	s := dynamic.NewServer(
		replicaId,
		ReplicaIdAddrMap,
		ReplicaIdList, NodeAddrList,
		//NetAddr, // TODO replaces the network address to server ID
		PaxosShard, FpShard,
		PaxosShardLeaderMap, FpShardLeaderMap,
		FollowerAddrList,
		HbInterval,
		CmdChSize,
		ExecChSize,
		IsExecReply,
		IsFpLeaderUsePaxos,
		IsGrpc, IsSyncSend,
		IsFpLeaderSoleLearner,
		IsLatPrediction,
		ProbeInv,
		WindowLen,
		WindowSize,
		LatPredictPth,
		IsPaxosFutureTime,
	)

	logger.Infof("Starting app")
	// Starts application
	startApp(s)

	logger.Infof("Starting RPC")
	// Starts RPC
	ipPort := strings.Split(NetAddr, ":")
	if len(ipPort) != 2 {
		logger.Fatalf("Invalid server network address = %s", NetAddr)
	}
	port := ipPort[1]

	s.Start(port)
}

func startApp(s *dynamic.Server) {
	kv := kv.NewKvStore()
	keyList := common.LoadKey(DataKeyFile)
	val := common.GenVal(DataValSize)
	kv.InitData(keyList, val)

	// Starts to wait on the execution channel from the consensus protocol
	go exec(s, kv)
}

func exec(s *dynamic.Server, kv *kv.KvStore) {
	w := common.NewFileWriter(AppLogDir + "/" + "kv-" + replicaId + ".log")
	logger.Info("KV waits for commands")

	c := s.GetExecCh()
	em := s.GetExecManager()

	for e := range c {
		cmd := e.GetCmd()
		shard := e.GetT().Shard

		ret := ""
		if IsExec {
			duration := time.Now().UnixNano() - e.GetStartDuration()
			logger.Infof("Executes idx = %s, duration = %v ns", cmd.Id , duration)
			ret = execCmd(kv, cmd)
		}

		if IsFpLeaderSoleLearner {
			if IsExecReply {
				if PaxosShard == shard || FpShard == shard {
					// Only shard leaders return execution results to clients
					execRet := em.GetExecRet(cmd.Id)
					execRet.C <- ret
				}
			}
		} else {
			if IsExecReply {
				// TODO Only replicas that are actually waiting for the exec resulst do so
				execRet := em.GetExecRet(cmd.Id)
				execRet.C <- ret
			}
		}

		if IsExecLog {
			w.Write(fmt.Sprintf("%s %s %s %s\n", cmd.Id, cmd.Type, cmd.Key, cmd.Val))
		}
	}
}

func execCmd(kv *kv.KvStore, cmd *dynamic.Command) string {
	if cmd.Type == "w" {
		kv.Write(cmd.Key, cmd.Val)
		return "true"
	} else if cmd.Type == "r" {
		ret, _ := kv.Read(cmd.Key)
		return ret
	}

	logger.Fatalf("Undefined cmd type = %s", cmd.Type)
	return ""
}
