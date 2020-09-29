package main

import (
	"flag"
	"math"

	"domino/common"
	"domino/dynamic/config"
)

func parseArgs() {
	flag.BoolVar(&isDebug, "d", false, "debug mode")
	flag.StringVar(&replicaId, "i", "", "server id")
	flag.StringVar(&configFile, "c", "", "configuration file")
	flag.StringVar(&replicaFile, "r", "", "replica location file")

	flag.Parse()

	if configFile == "" {
		logger.Fatalf("Invalid configuration file.")
	}
	if replicaFile == "" {
		logger.Fatalf("Invalid replica location file.")
	}
	if replicaId == "" {
		logger.Fatalf("Invalid server id.")
	}
}

func loadConfig(propertyFile, replicaFile string) {
	// Loading Replica Information
	replicaDir := config.LoadReplicaInfo(replicaFile)

	// Replica Information
	for _, rInfo := range replicaDir {
		ReplicaNum++
		addr := rInfo.GetNetAddr()
		if rInfo.IsFpLeader {
			LeaderAddr = addr
		} else {
			FollowerAddrList = append(FollowerAddrList, addr)
		}
	}
	if ReplicaNum != len(FollowerAddrList)+1 {
		logger.Fatalf("Error: ReplicaNum = %d, but follower num = %d", ReplicaNum, len(FollowerAddrList))
	}
	f := (ReplicaNum - 1) / 2
	MajorityNum = f + 1
	FastQuorum = int(math.Ceil((3.0*float64(f))/2.0)) + 1

	// Network Addr
	rInfo, ok := replicaDir[replicaId]
	if !ok {
		logger.Fatalf("Cannot find replica information for replicaId = %s", replicaId)
	}
	ServerAddr = rInfo.GetNetAddr()
	if rInfo.IsFpLeader {
		IsLeader = true
	}

	// Replica properties
	p := common.NewProperties()
	p.Load(propertyFile)

	cmdChSize := p.GetIntWithDefault(config.FLAG_DYNAMIC_PAXOS_CMD_CH_SIZE, "1048576")
	execChSize := p.GetIntWithDefault(config.FLAG_DYNAMIC_EXEC_CH_SIZE, "10485760")
	FastPaxosVoteChannelBufferSize = cmdChSize
	FastPaxosCommadChannelBufferSize = cmdChSize
	FastPaxosExecutionChannelBufferSize = execChSize

	FastPaxosIsExec = p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC, "false")
	FastPaxosIsExecReply = p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC_REPLY, "false")
	FastPaxosLogPersistEnabled = p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC_LOG, "false")

	FastPaxosLogType = "fixed"
	FastPaxosLogSize = 1024 * 1024 * 8
	FastPaxosLogFilePath = p.GetStr(config.FLAG_DYNAMIC_KV_LOG_DIR)
	FastPaxosLogFilePath += "/kv-" + replicaId + ".log" // TODO change the variable name to AppLogDir

	FastPaxosScheduler = common.NoScheduler // No scheduler for now
	FastPaxosProcessWindow = 0
	// Applications
	DataKeyFile = p.GetStr(config.FLAG_APP_DATA_KEY_FILE)
	DataValSize = p.GetIntWithDefault(config.FLAG_APP_DATA_VAL_SIZE, "8")
}
