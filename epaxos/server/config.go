package main

import (
	"flag"
	"strconv"
	//"strings"

	"domino/common"
	"domino/dynamic/config"
)

func parseArgs() {
	flag.BoolVar(&isDebug, "d", false, "debug mode")
	flag.StringVar(&serverId, "i", "", "server id")
	flag.StringVar(&configFile, "c", "", "server configuration file")
	flag.StringVar(&replicaFile, "r", "", "replica location file")
	flag.StringVar(&protocolType, "p", "", "Protocol type: e (EPaxos), m (Mencius), p (Multi-Paxos) or g (GPaxos) (optional)")
	flag.StringVar(&isMenciusOpt, "m", "", "true or false: enable Mencius commit early (optional)")
	flag.StringVar(&isThriftyOpt, "t", "", "true of false: enable thrifty for EPaxos and Multi-Paxos (optional)")

	flag.Parse()

	if serverId == "" {
		flag.Usage()
		logger.Fatal("Invalid server id.")
	}
	if configFile == "" {
		flag.Usage()
		logger.Fatal("Invalid configuration file.")
	}

	if replicaFile == "" {
		flag.Usage()
		logger.Fatal("Invalid replication location file.")
	}
}

const (
	Flag_protocol   = "epaxos.protocol" //e: epaxos, m: mencius; g: gpaxos; p: multi-paxos
	Flag_exec       = "epaxos.exec"
	Flag_dreply     = "epaxos.dreply"
	Flag_beacon     = "epaxos.beacon"
	Flag_durable    = "epaxos.durable"
	Flag_thrifty    = "epaxos.thrifty"
	Flag_cpuprofile = "epaxos.cpuprofile"

	// Added mencius optimization
	Flag_mencius_early_commit_ack = "mencius.early_commit_ack"

	// Added by @skoya76
	Flag_measure_commit_to_exec_time = "epaxos.measure.cmtexec.time"
)

func loadConfig(configFile, replicaFile string) {
	// NOTE nodeList must be sorted by the serverId that is from 1 to n
	replicaDir := config.LoadReplicaInfo(replicaFile)
	nodeList = make([]string, len(replicaDir))
	for rId, rInfo := range replicaDir {
		rAddr := rInfo.GetNetAddr()
		if r, e := strconv.Atoi(rId); e != nil {
			logger.Fatalf("Invalid EPaxos replicaId = %s")
		} else {
			nodeList[r-1] = rAddr
		}
	}

	rInfo, ok := replicaDir[serverId]
	if !ok {
		logger.Fatalf("No replica info for replicaId = %s", serverId)
	}

	// Sets current server's network address
	myAddr = rInfo.Ip
	if p, e := strconv.Atoi(rInfo.Port); e != nil {
		logger.Fatalf("Invalid port %s error %v", rInfo.Port, e)
	} else {
		portnum = p
	}

	p := common.NewProperties()
	p.Load(configFile)

	// Init data for KV store
	keyFile = p.GetStr(config.FLAG_APP_DATA_KEY_FILE)
	valLen = p.GetInt(config.FLAG_APP_DATA_VAL_SIZE)

	// protocol
	pType := p.GetStr(Flag_protocol)
	if protocolType != "" {
		pType = protocolType
	}
	switch pType {
	case "e":
		doEpaxos = true
	case "m":
		doMencius = true
	case "g":
		doGpaxos = true
	case "p":
		doPaxos = true
	default:
		logger.Fatalf("Unknown protocol type = %s", pType)
	}

	exec = p.GetBoolWithDefault(Flag_exec, "false")
	dreply = p.GetBoolWithDefault(Flag_dreply, "false")
	beacon = p.GetBoolWithDefault(Flag_beacon, "false")
	durable = p.GetBoolWithDefault(Flag_durable, "false")
	thrifty = p.GetBoolWithDefault(Flag_thrifty, "false")
	if isThriftyOpt != "" {
		if v, err := strconv.ParseBool(isThriftyOpt); err != nil {
			flag.Usage()
			logger.Fatalf("Invalid command-line isThriftyOpt = %s", isThriftyOpt)
		} else {
			thrifty = v
		}
	}

	cpuprofile = p.GetWithDefault(Flag_cpuprofile, cpuprofile)

	// Mencius optimization
	mencius_early_commit_ack = p.GetBoolWithDefault(Flag_mencius_early_commit_ack, "false")
	if isMenciusOpt != "" {
		if v, err := strconv.ParseBool(isMenciusOpt); err != nil {
			flag.Usage()
			logger.Fatalf("Invalid command-line isMenciusOpt = %s", isMenciusOpt)
		} else {
			mencius_early_commit_ack = v
		}
	}
	measure_commit_to_exec_time = p.GetBoolWithDefault(Flag_measure_commit_to_exec_time, "false")
}

func getBool(ct map[string]string, flag string, d bool) bool {
	if v, ok := ct[flag]; ok {
		if ret, err := strconv.ParseBool(v); err != nil {
			logger.Fatalf("Invalid %s = %s error: %v", flag, v, err)
		} else {
			return ret
		}
	}
	return d
}

func getStr(ct map[string]string, flag string, d string) string {
	if v, ok := ct[flag]; ok {
		return v
	}
	return d
}
