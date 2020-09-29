package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"domino/benchmark/config"
	"domino/benchmark/workload"
	"domino/common"
)

// Config variables
var IsDebug bool
var ClientId string
var DcId string
var TargetDcId string // for EPaxos, the replica's DcId that a client will send requests to. If empty, the client will
var ConfigFile string
var ReplicaFile string
var ProtocolType string // Overwrites config file option
var Duration time.Duration
var RandomSeed int64

// If waits for execution result
var IsExec bool
var KeyFile string
var ValSize int

// Transaction rate settings
var TxnTotalNum int
var TxnTargetRate int
var IsOpenLoop bool

// Retry settings
var RetryMode string
var RetryMax int64
var RetryInterval time.Duration
var RetrySlotMax int64 // only for exponential backoff
// Workload settings
var Workload workload.Workload

var ClientType string

// Configuration flags
const (
	CONFIG_DEBUG         = "d"
	CONFIG_DEBUG_HELPER  = "debug mode."
	CONFIG_DEBUG_DEFAULT = false

	CONFIG_ID         = "i"
	CONFIG_ID_HELPER  = "client id"
	CONFIG_ID_DEFAULT = config.INVALID_STR

	CONFIG_DC         = "dc"
	CONFIG_DC_HELPER  = "datacenter id"
	CONFIG_DC_DEFAULT = config.INVALID_STR

	CONFIG_TARGET         = "t"
	CONFIG_TARGET_HELPER  = "target replica's dc id"
	CONFIG_TARGET_DEFAULT = config.INVALID_STR

	CONFIG_CLIENT         = "c"
	CONFIG_CLIENT_HELPER  = "Configuration file for client-side lib. <REQUIRED>."
	CONFIG_CLIENT_DEFAULT = config.INVALID_STR

	CONFIG_REPLICA         = "r"
	CONFIG_REPLICA_HELPER  = "Replica location file"
	CONFIG_REPLICA_DEFAULT = config.INVALID_STR

	CONFIG_PROTOCOL         = "p"
	CONFIG_PROTOCOL_HELPER  = "Protocol Type: d (Dynamic), e (EPaxos), m (Mencius), or p (MultiPaxos) (optional)"
	CONFIG_PROTOCOL_DEFAULT = config.INVALID_STR
)

func ParseArgs() {
	flag.BoolVar(
		&IsDebug,
		CONFIG_DEBUG,
		CONFIG_DEBUG_DEFAULT,
		CONFIG_DEBUG_HELPER,
	)

	flag.StringVar(
		&ClientId,
		CONFIG_ID,
		CONFIG_ID_DEFAULT,
		CONFIG_ID_HELPER,
	)

	flag.StringVar(
		&DcId,
		CONFIG_DC,
		CONFIG_DC_DEFAULT,
		CONFIG_DC_HELPER,
	)

	flag.StringVar(
		&TargetDcId,
		CONFIG_TARGET,
		CONFIG_TARGET_DEFAULT,
		CONFIG_TARGET_HELPER,
	)

	flag.StringVar(
		&ConfigFile,
		CONFIG_CLIENT,
		CONFIG_CLIENT_DEFAULT,
		CONFIG_CLIENT_HELPER,
	)

	flag.StringVar(
		&ReplicaFile,
		CONFIG_REPLICA,
		CONFIG_REPLICA_DEFAULT,
		CONFIG_REPLICA_HELPER,
	)

	flag.StringVar(
		&ProtocolType,
		CONFIG_PROTOCOL,
		CONFIG_PROTOCOL_DEFAULT,
		CONFIG_PROTOCOL_HELPER,
	)

	flag.Parse()

	if !validStr(ConfigFile) || !validStr(ReplicaFile) || !validStr(ClientId) || !validStr(DcId) {
		flag.Usage()
		os.Exit(-1)
	}
}

func LoadBenchmarkConfig() {
	p := common.NewProperties()
	p.Load(ConfigFile)

	ClientType = p.GetWithDefault(config.FLAG_CLIENT_TYPE, "")
	if ProtocolType != config.INVALID_STR {
		ClientType = ProtocolType
	}

	// Execution duration in seconds
	Duration = p.GetTimeDurationWithDefault(config.FLAG_EXEC_DURATION, config.DEFAULT_EXEC_DURATION)

	TxnTotalNum = p.GetIntWithDefault(
		config.FLAG_EXEC_TXN_TOTAL_NUM,
		config.DEFAULT_EXEC_TXN_TOTAL_NUM)

	TxnTargetRate = p.GetIntWithDefault(
		config.FLAG_EXEC_TXN_TARGET_RATE,
		config.DEFAULT_EXEC_TXN_TARGET_RATE)

	// Open loop switch
	IsOpenLoop = p.GetBoolWithDefault(
		config.FLAG_EXEC_TXN_OPEN_LOOP,
		config.DEFAULT_EXEC_TXN_OPEN_LOOP)

	// Random seed
	RandomSeed = p.GetInt64WithDefault(config.FLAG_RANDOM_SEED, config.DEFAULT_RANDOM_SEED)
	if RandomSeed == 0 {
		RandomSeed = int64(time.Now().Nanosecond())
	}
	rand.Seed(RandomSeed)

	// Retry configuration
	RetryMode = p.GetWithDefault(config.FLAG_RETRY_MODE, config.DEFAULT_RETRY_MODE)
	if RetryMode != config.RETRY_MODE_CONSTANT && RetryMode != config.RETRY_MODE_EXP {
		logger.Fatalf("Invalid transaction retry mode. Expects %s or %s. Default is %s",
			config.RETRY_MODE_CONSTANT, config.RETRY_MODE_EXP, config.RETRY_MODE_EXP)
	}

	RetryMax = p.GetInt64WithDefault(config.FLAG_RETRY_MAX, config.DEFAULT_RETRY_MAX)

	RetryInterval = p.GetTimeDurationWithDefault(config.FLAG_RETRY_INTERVAL, config.DEFAULT_RETRY_INTERVAL)

	RetrySlotMax = p.GetInt64WithDefault(
		config.FLAG_RETRY_MAX_SLOT,
		config.DEFAULT_RETRY_MAX_SLOT)

	// Execution mode
	IsExec = p.GetBoolWithDefault(
		common.Flag_fastpaxos_exec_mode,
		common.Default_fastpaxos_exec_mode)

	// Key file
	KeyFile = p.GetStr(common.Flag_data_key_file)
	if !validStr(KeyFile) {
		logger.Fatalf("Key file is not specified")
	}
	// Load keys
	keyList := common.LoadKey(KeyFile)

	// Value size
	ValSize = p.GetIntWithDefault(common.Flag_data_val_size, common.Default_data_val_size)

	// Workload configuration
	if Workload = createWorkload(p, keyList, ValSize); Workload == nil {
		logger.Fatal("Failed to create a workload instance.")
	}

	logger.Debug("Benchmark Configurations:", p.GetPropMap())
	logger.Debug("Debug mode =", IsDebug)
	logger.Debug("Client config file =", ConfigFile)
	logger.Debug("Running duration =", Duration)
	logger.Debug("Total number of transactions =", TxnTotalNum)
	logger.Debug("Target transaction rate =", TxnTargetRate)
	logger.Debug("Random seed =", RandomSeed)
	logger.Debug("Transaction retry mode =", RetryMode)
	logger.Debug("Transaction retry max num =", RetryMax)
	logger.Debug("Transaction retry interval =", RetryInterval)
	logger.Debug("Transaction retry slot max =", RetrySlotMax)
	logger.Debug("Workload type =", Workload)
	logger.Debug("Key file =", KeyFile)
	logger.Debug("Value size =", ValSize)
	logger.Debug("Do FastPaxos Execution mode (TODO delete )=", IsExec)
}

func validStr(str string) bool {
	if str == config.INVALID_STR {
		return false
	}
	return true
}

func createWorkload(
	p *common.Properties,
	keyList []string,
	valSize int,
) workload.Workload {
	zipfAlpha := p.GetFloat64WithDefault(
		config.FLAG_WORKLOAD_ZIPF_ALPHA,
		workload.WORKLOAD_DEFAULT_ZIPF_ALPHA)
	baseWorkload := workload.NewAbstractWorkload(keyList, zipfAlpha, valSize)
	workloadType := p.GetWithDefault(config.FLAG_WORKLOAD_TYPE, config.DEFAULT_WORKLOAD)

	logger.Debugf("Workload configured as %s", workloadType)

	switch workloadType {
	case config.WORKLOAD_RETWIS:
		addUserRatio := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_RETWIS_ADD_USER_RATIO,
			workload.RETWIS_DEFAULT_ADD_USER_RATIO)

		followRatio := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_RETWIS_FOLLOW_RATIO,
			workload.RETWIS_DEFAULT_FOLLOW_UNFOLLOW_RATIO)

		postRatio := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_RETWIS_POST_RATIO,
			workload.RETWIS_DEFAULT_POST_TWEET_RATIO)

		loadRatio := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_RETWIS_LOAD_RATIO,
			workload.RETWIS_DEFAULT_LOAD_TIMELINE_RATIO)

		return workload.NewRetwisWorkload(
			baseWorkload,
			addUserRatio,
			followRatio,
			postRatio,
			loadRatio,
		)

	case config.WORKLOAD_YCSBT:
		readNum := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_YCSBT_READ_NUM,
			workload.YCSBT_DEFAULT_READ_NUM_PER_TXN)

		writeNum := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_YCSBT_WRITE_NUM,
			workload.YCSBT_DEFAULT_WRITE_NUM_PER_TXN)

		return workload.NewYcsbtWorkload(
			baseWorkload,
			readNum,
			writeNum,
		)

	case config.WORKLOAD_ONETXN:
		// TODO defines properties for one_txn workload instead of re-using the ycsbt ones
		readNum := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_YCSBT_READ_NUM,
			workload.YCSBT_DEFAULT_READ_NUM_PER_TXN)

		writeNum := p.GetIntWithDefault(
			config.FLAG_WORKLOAD_YCSBT_WRITE_NUM,
			workload.YCSBT_DEFAULT_WRITE_NUM_PER_TXN)

		return workload.NewOneTxnWorkload(
			baseWorkload,
			readNum,
			writeNum,
		)

	default:
		logger.Fatalf("Unkown workload type: %s", workloadType)
	}

	return nil
}
