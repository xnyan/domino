package config

// Configuration flags
const (
	FLAG_DYNAMIC_FP_PREDICT_TIMEOFFSET             = "dynamic.fp.predict.timeoffset"     // default true, predict time offset
	FLAG_DYNAMIC_FP_PREDICT_ALL                    = "dynamic.fp.predict.all"            // default true, using the latency to all replicas
	FLAG_DYNAMIC_LAT_PROBE_INTERVAL                = "dynamic.lat.probe.interval"        // default 10ms
	FLAG_DYNAMIC_LAT_PROBE_WINDOW_LEN              = "dynamic.lat.probe.window.length"   // default 1s
	FLAG_DYNAMIC_LAT_PROBE_WINDOW_MIN_SIZE         = "dynamic.lat.probe.window.min.size" // default 10
	FLAG_DYNAMIC_CLIENT_ADD_DELAY                  = "dynamic.client.add.delay"
	FLAG_DYNAMIC_LAT_REPLICA_PROBE_INVTERVAL       = "dynamic.lat.replica.probe.interval"        // default 10ms
	FLAG_DYNAMIC_LAT_REPLICA_PROBE_WINDOW_LEN      = "dynamic.lat.replica.probe.window.length"   // default 1s
	FLAG_DYNAMIC_LAT_REPLICA_PROBE_WINDOW_MIN_SIZE = "dynamic.lat.replica.probe.window.min.size" // default 10

	FLAG_DYNAMIC_PREDICT_PERCENTILE    = "dynamic.lat.predict.percentile" // default 0.95, >= 0 <= 1.0
	DEFAULT_DYNAMIC_PREDICT_PERCENTILE = "0.95"

	FLAG_DYNAMIC_PAXOS_HB_INTERVAL = "dynamic.heartbeat.interval"

	FLAG_DYNAMIC_EXEC              = "dynamic.exec"       // If execute a committed cmd
	FLAG_DYNAMIC_EXEC_REPLY        = "dynamic.exec.reply" // If replies cmd execution result
	FLAG_DYNAMIC_EXEC_LOG          = "dynamic.exec.log"   // If log executed cmd (in order)
	FLAG_DYNAMIC_PAXOS_CMD_CH_SIZE = "dynamic.cmd.ch.size"
	FLAG_DYNAMIC_EXEC_CH_SIZE      = "dynamic.exec.ch.size"

	// If fast paxos rejects a command, if the fast paxos leader is a paxos
	// leader, it will use the paxos to accept the command. Default true.
	FLAG_DYNAMIC_FP_LEADER_USE_PAXOS   = "dynamic.fp.leader.use.paxos"
	FLAG_DYNAMIC_FP_LEADER_SOLELEARNER = "dynamic.fp.leader.solelearner" // default true, the leader is the sole learner
	FLAG_DYNAMIC_PAXOS_FUTURE_TIME     = "dynamic.paxos.future.time"     // default false, whether Paxos assigns a request with  a future time
	FLAG_DYNAMIC_GRPC                  = "dynamic.grpc"                  // default true, using gRPC
	FLAG_DYNAMIC_SYNC_SEND             = "dynamic.sync.send"             // default false, not sync send

	FLAG_DYNAMIC_KV_LOG_DIR = "dynamic.kv.log.dir"

	// [0, 100] percentage, negative to disable. default -1
	FLAG_DYNAMIC_CLIENT_FP_LOAD = "dynamic.client.fp.load"

	FLAG_APP_DATA_KEY_FILE = "data.key.file"
	FLAG_APP_DATA_VAL_SIZE = "data.val.size"
)
