package main

import (
	"flag"

	"domino/common"
	"domino/dynamic/config"
)

func err(errMsg string) {
	flag.Usage()
	logger.Fatal(errMsg)
}

func parseArgs() {
	flag.BoolVar(&isDebug, "d", false, "debug mode")
	flag.StringVar(&replicaId, "i", "", "server id")
	flag.StringVar(&configFile, "c", "", "configuration file")
	flag.StringVar(&replicaFile, "r", "", "replica location file")

	flag.Parse()

	if configFile == "" {
		err("Invalid configuration file.")
	}
	if replicaFile == "" {
		err("Invalid replica location file.")
	}
	if replicaId == "" {
		err("Invalid server id.")
	}
}

func loadConfig(propertyFile, replicaFile string) {
	// Loading Replica Information
	replicaDir := config.LoadReplicaInfo(replicaFile)

	// Shard information
	var shardNum int
	shardNum, PaxosShardLeaderMap, FpShardLeaderMap = config.GenShardInfo(replicaDir)
	// Paxos Shard
	for pShard, rId := range PaxosShardLeaderMap {
		if rId == replicaId {
			PaxosShard = pShard
		}
	}
	// Fast Paxos Shard
	for fpShard, rId := range FpShardLeaderMap {
		if rId == replicaId {
			FpShard = fpShard
		}
	}
	// Replica Information
	for rId, rInfo := range replicaDir {
		addr := rInfo.GetNetAddr()
		if rId != replicaId {
			FollowerAddrList = append(FollowerAddrList, addr)
		}
		ReplicaIdList = append(ReplicaIdList, rId)
		NodeAddrList = append(NodeAddrList, addr)
		ReplicaIdAddrMap[rId] = addr
	}
	// Network Addr
	rInfo, ok := replicaDir[replicaId]
	if !ok {
		logger.Fatalf("Cannot find replica information for replicaId = %s", replicaId)
	}
	NetAddr = rInfo.GetNetAddr()

	// Replica properties
	p := common.NewProperties()
	p.Load(propertyFile)

	HbInterval = p.GetTimeDurationWithDefault(config.FLAG_DYNAMIC_PAXOS_HB_INTERVAL, "10ms")
	IsExec = p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC, "false")
	IsExecReply = p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC_REPLY, "false")
	IsExecLog = p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC_LOG, "false")
	CmdChSize = p.GetIntWithDefault(config.FLAG_DYNAMIC_PAXOS_CMD_CH_SIZE, "1048576")
	ExecChSize = p.GetIntWithDefault(config.FLAG_DYNAMIC_EXEC_CH_SIZE, "10485760")

	IsFpLeaderUsePaxos = p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_LEADER_USE_PAXOS, "true")
	IsFpLeaderSoleLearner = p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_LEADER_SOLELEARNER, "true")
	IsGrpc = p.GetBoolWithDefault(config.FLAG_DYNAMIC_GRPC, "true")
	IsSyncSend = p.GetBoolWithDefault(config.FLAG_DYNAMIC_SYNC_SEND, "false")

	AppLogDir = p.GetStr(config.FLAG_DYNAMIC_KV_LOG_DIR)

	// Latency prediction
	IsLatPrediction = true
	ProbeInv = p.GetTimeDurationWithDefault(config.FLAG_DYNAMIC_LAT_REPLICA_PROBE_INVTERVAL, "10ms")
	WindowLen = p.GetTimeDurationWithDefault(config.FLAG_DYNAMIC_LAT_REPLICA_PROBE_WINDOW_LEN, "1s")
	WindowSize = p.GetIntWithDefault(config.FLAG_DYNAMIC_LAT_REPLICA_PROBE_WINDOW_MIN_SIZE, "10")
	LatPredictPth = p.GetFloat64WithDefault(
		config.FLAG_DYNAMIC_PREDICT_PERCENTILE,
		config.DEFAULT_DYNAMIC_PREDICT_PERCENTILE)
	if LatPredictPth < 0.0 || LatPredictPth > 1.0 {
		logger.Fatalf("Invalid percentile = %f for prediction. Expected [0.0, 1.0]", LatPredictPth)
	}

	// Paxos uses a future time
	IsPaxosFutureTime = p.GetBoolWithDefault(config.FLAG_DYNAMIC_PAXOS_FUTURE_TIME, "false")

	// Applications
	DataKeyFile = p.GetStr(config.FLAG_APP_DATA_KEY_FILE)
	DataValSize = p.GetIntWithDefault(config.FLAG_APP_DATA_VAL_SIZE, "8")

	logger.Infof("paxos shard = %d, fp shard = %d", PaxosShard, FpShard)
	logger.Infof("shard num = %d", shardNum)
	logger.Infof("followerList = %v", FollowerAddrList)
}
