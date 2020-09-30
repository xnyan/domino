package main

import (
	"flag"
	"time"

	"github.com/op/go-logging"

	"domino/trace/common"
	"domino/trace/node"
)

var logger = logging.MustGetLogger("Client")

const (
	Flag_probe_inv      = "probe.inv"
	Flag_probe_duration = "probe.duration"
	Flag_log_dir        = "log.dir"
)

var isDebug bool
var dcId string
var configFile string
var locFile string

var dcTable map[string]string // net addr --> dc id
var logDir string
var inv time.Duration
var duration time.Duration

func main() {
	// parses command line args
	parseArgs()

	// Configuration based on args
	config()

	// Starts a client instance
	c := node.NewClient(dcId, dcTable, logDir)
	node.StartClient(c, inv, duration)
}

func parseArgs() {
	flag.BoolVar(&isDebug, "d", false, "debug mode")
	flag.StringVar(&dcId, "i", "", "datacenter id")
	flag.StringVar(&configFile, "c", "", "configuration file")
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

	if configFile == "" {
		flag.Usage()
		logger.Fatalf("Missing client configuration file")
	}
}

func config() {
	common.ConfigLogger(isDebug)

	m := common.NewMapParser()
	m.Load(locFile)
	dcTable = m.GetMap()
	if dcTable == nil || len(dcTable) == 0 {
		logger.Fatalf("Empty datacenter and server addresses in location file %s", locFile)
	}

	p := common.NewProperties()
	p.Load(configFile)
	inv = parseDuration(p, Flag_probe_inv, configFile)
	duration = parseDuration(p, Flag_probe_duration, configFile)
	logDir, _ = p.Get(Flag_log_dir)
}

func parseDuration(p *common.Properties, flag, configFile string) time.Duration {
	if s, ok := p.Get(flag); ok {
		t, err := time.ParseDuration(s)
		if err != nil {
			invalidFlagVal(flag, s, configFile)
		}
		return t
	} else {
		missFlag(flag, configFile)
	}
	return -1
}

func missFlag(flag, configFile string) {
	logger.Fatalf("Invalid %s in file %s", flag, configFile)
}

func invalidFlagVal(flag, val, configFile string) {
	logger.Fatalf("Invalid %s = %s in file %s", flag, val, configFile)
}
