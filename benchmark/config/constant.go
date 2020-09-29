package config

// constants
const (
	INVALID_STR = ""

	FLAG_CLIENT_TYPE = "benchmark.client"

	// Execution duration, unit: seconds
	FLAG_EXEC_DURATION = "benchmark.duration"

	// Total number of transactions to execute if duration is not positive
	FLAG_EXEC_TXN_TOTAL_NUM = "benchmark.txn.total"

	// Target number of transactons sent per second
	FLAG_EXEC_TXN_TARGET_RATE = "benchmark.txn.target.rate"

	FLAG_EXEC_TXN_OPEN_LOOP = "benchmark.openloop"

	// Random seed. 0 for a dynamic seed. Default is 0.
	FLAG_RANDOM_SEED = "benchmark.random.seed"

	// Retry configurations
	FLAG_RETRY_MAX      = "benchmark.txn.retry.max"
	FLAG_RETRY_MODE     = "benchmark.txn.retry.mode"
	FLAG_RETRY_INTERVAL = "benchmark.txn.retry.interval"
	FLAG_RETRY_MAX_SLOT = "benchmark.txn.retry.maxslot"

	FLAG_WORKLOAD_TYPE                  = "workload.type"
	FLAG_WORKLOAD_ZIPF_ALPHA            = "workload.zipf.alpha"
	FLAG_WORKLOAD_RETWIS_ADD_USER_RATIO = "workload.retwis.adduser.ratio"
	FLAG_WORKLOAD_RETWIS_FOLLOW_RATIO   = "workload.retwis.followunfollow.ratio"
	FLAG_WORKLOAD_RETWIS_POST_RATIO     = "workload.retwis.posttweet.ratio"
	FLAG_WORKLOAD_RETWIS_LOAD_RATIO     = "workload.retwis.loadtimeline.ratio"
	FLAG_WORKLOAD_YCSBT_READ_NUM        = "workload.ycsbt.readnum"
	FLAG_WORKLOAD_YCSBT_WRITE_NUM       = "workload.ycsbt.writenum"
)

const (
	DEFAULT_EXEC_DURATION        = "10s"
	DEFAULT_EXEC_TXN_TOTAL_NUM   = "1000"
	DEFAULT_EXEC_TXN_TARGET_RATE = "0"
	DEFAULT_EXEC_TXN_OPEN_LOOP   = "false"
	DEFAULT_RANDOM_SEED          = "0"
	DEFAULT_RETRY_MAX            = "0"        // max number of retrying a transaction
	DEFAULT_RETRY_MODE           = "constant" // retry strategy
	DEFAULT_RETRY_INTERVAL       = "10ms"     // time interval between retries, unit: ms
	DEFAULT_RETRY_MAX_SLOT       = "32"       // the max number of slots for exponential backoff

	// Client type
	//CLIENT_DO      = "do"
	//CLIENT_HYBRID  = "h"  // hybrid Paxos
	CLIENT_FP      = "fp" // basic Fast Paxos
	CLIENT_DYNAMIC = "d"  // dynamic Paxos
	CLIENT_EPAXOS  = "e"  // epaxos
	CLIENT_MENCIUS = "m"  // mencius
	CLIENT_GPAXOS  = "g"  // generalized paxos
	CLIENT_PAXOS   = "p"  // multi-paxos

	// Retry mode type
	RETRY_MODE_CONSTANT = "constant" // fixed retry interval
	RETRY_MODE_EXP      = "exp"      // exponential backoff

	DEFAULT_WORKLOAD = "ycsbt" // workload type
	WORKLOAD_RETWIS  = "retwis"
	WORKLOAD_YCSBT   = "ycsbt"
	WORKLOAD_ONETXN  = "onetxn" // uses only a transaction having the same set of reads & writes
)
