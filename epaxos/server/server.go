package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"

	"github.com/op/go-logging"

	domino_common "domino/common"
	"domino/epaxos/common"
	"domino/epaxos/epaxos"
	"domino/epaxos/gpaxos"
	"domino/epaxos/mencius"
	"domino/epaxos/paxos"
)

//var procs *int = flag.Int("p", 2, "GOMAXPROCS. Defaults to 2")

var logger = logging.MustGetLogger("server")

// command-line
var isDebug bool = false
var serverId string = ""
var configFile string = ""
var replicaFile string = ""
var protocolType string = ""
var isMenciusOpt string = ""
var isThriftyOpt string = ""

// config file
var serverAddr string = ""
var keyFile string = ""
var valLen int = 1024
var doMencius bool = false
var doGpaxos bool = false
var doEpaxos bool = false
var doPaxos bool = false
var exec bool = false
var dreply bool = false
var beacon bool = false
var durable bool = false
var thrifty bool = false
var cpuprofile string = ""
var mencius_early_commit_ack bool = false
var measure_commit_to_exec_time bool = false

var myAddr string = ""
var portnum int = 0
var replicaId int = -1
var nodeList []string = make([]string, 0)

func main() {
	parseArgs()
	domino_common.ConfigLogger(isDebug)
	loadConfig(configFile, replicaFile)

	keyList := common.LoadKey(keyFile)
	initVal := common.GenVal(valLen)

	//runtime.GOMAXPROCS(*procs)
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt)
		go catchKill(interrupt)
	}

	if rId, err := strconv.Atoi(serverId); err != nil {
		logger.Fatalf("error: %v", err)
	} else {
		replicaId = rId - 1 // NOTE: 0 is the leader
	}

	logger.Infof("Server starting on port %d\n", portnum)
	logger.Infof("replicaId = %d, nodeList = %s", replicaId, nodeList)
	logger.Infof("exec = %t, dreply = %t, thrify = %t", exec, dreply, thrifty)

	if doEpaxos {
		logger.Info("Starting Egalitarian Paxos replica...")
		rep := epaxos.NewReplica(replicaId, nodeList, thrifty, exec, dreply, beacon, durable, keyList, initVal, measure_commit_to_exec_time)
		rpc.Register(rep)
	} else if doMencius {
		logger.Info("Starting Mencius replica...")
		mencius.EARLY_COMMIT_ACK = mencius_early_commit_ack
		thrifty = false // Mencius does not have thrifty optimization. Enabling this could cause execution failures
		logger.Infof("Mencius early commit = %t, thrifty = %t", mencius.EARLY_COMMIT_ACK, thrifty)
		rep := mencius.NewReplica(replicaId, nodeList, thrifty, exec, dreply, durable, keyList, initVal, measure_commit_to_exec_time)
		rpc.Register(rep)
	} else if doGpaxos {
		logger.Info("Starting Generalized Paxos replica...")
		rep := gpaxos.NewReplica(replicaId, nodeList, thrifty, exec, dreply, keyList, initVal, measure_commit_to_exec_time)
		rpc.Register(rep)
	} else if doPaxos {
		logger.Info("Starting classic Paxos replica...")
		rep := paxos.NewReplica(replicaId, nodeList, thrifty, exec, dreply, durable, keyList, initVal, measure_commit_to_exec_time)
		rpc.Register(rep)
	} else {
		logger.Fatalf("No protocol type specified")
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portnum+1000))
	if err != nil {
		log.Fatal("listen error:", err)
	}

	http.Serve(l, nil)
}

func catchKill(interrupt chan os.Signal) {
	<-interrupt
	if cpuprofile != "" {
		pprof.StopCPUProfile()
	}
	fmt.Println("Caught signal")
	os.Exit(0)
}
