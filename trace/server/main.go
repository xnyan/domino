package main

import (
	"flag"

	"github.com/op/go-logging"

	"domino/trace/common"
	"domino/trace/node"
)

var logger = logging.MustGetLogger("Server")

var isDebug bool
var dcId string    // current datacenter id
var locFile string // datacenter and server address mapping file

var addr string // ip:port

func main() {
	// Parses command line args
	parseArgs()

	// Configuration based on args
	config()

	// Starts a server instance
	s := node.NewServer(addr)
	logger.Infof("Starting server at addr = %s", addr)
	node.StartServer(s)
}

func parseArgs() {
	flag.BoolVar(&isDebug, "d", false, "debug mode")
	flag.StringVar(&dcId, "i", "", "datacenter id")
	flag.StringVar(&locFile, "l", "", "location file")

	flag.Parse()

	if dcId == "" {
		flag.Usage()
		logger.Fatalf("Missing datacenter id.")
	}

	if locFile == "" {
		flag.Usage()
		logger.Fatalf("Missing datacenter and server address location file.")
	}
}

func config() {
	common.ConfigLogger(isDebug)

	m := common.NewMapParser()
	m.Load(locFile)

	var ok bool
	addr, ok = m.Get(dcId)
	if !ok {
		logger.Fatalf("No datacenter id = %s in file %s", dcId, locFile)
	}
}
