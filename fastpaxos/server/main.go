package main

import (
	"time"

	"github.com/op/go-logging"

	"domino/common"
)

var logger = logging.MustGetLogger("server")

// Command-line configurations
var isDebug bool = false
var replicaId string = ""
var configFile string = ""
var replicaFile string = "" // replica location file

// Configurations from config file
var ServerAddr string = ""
var IsLeader bool = false
var LeaderAddr string = ""
var FollowerAddrList []string = make([]string, 0)
var ReplicaNum int = 0
var MajorityNum int = 0
var FastQuorum int = 0

var FastPaxosVoteChannelBufferSize int = 1024 * 8
var FastPaxosCommadChannelBufferSize int = 1024 * 8
var FastPaxosExecutionChannelBufferSize int = 1024 * 8
var FastPaxosIsExec = false
var FastPaxosIsExecReply = false
var FastPaxosLogPersistEnabled bool = true

//Log type
//fixed : a log that has fixed size
//seg : a log that consists of multiple segments, dynamically increase
var FastPaxosLogType string = ""
var FastPaxosLogSize int64 = 1024 * 1024 * 4
var FastPaxosLogFilePath string = "./"

//// Scheduler
//// no : raw fast paxos
//// delay : delay the process of a proposal based on network delays across DCs
//// timestamp : order proposals based on their timestamps, and propose them
var FastPaxosScheduler string = ""
var FastPaxosProcessWindow time.Duration

//// KV store data
var DataKeyFile string = ""
var DataValSize int = 1024

///////////////////

func main() {
	// Parses command-line arguments
	parseArgs()
	// Configures logging
	common.ConfigLogger(isDebug)
	// Parses configuration files
	loadConfig(configFile, replicaFile)

	server := NewServer(replicaId, ServerAddr, IsLeader)

	server.Start()
}
