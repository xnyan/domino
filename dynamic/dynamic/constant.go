package dynamic

const (
	NO_OP_ID = "NO_OP_ID"
)

const (
	INVALID_FP_IDX = -1
	INVALID_SHARD  = -1
)

// Message types for communications among replicas
const (
	REPLICA_MSG_HEART_BEAT         = 0
	REPLICA_MSG_PAXOS_ACCEPT_REQ   = 1
	REPLICA_MSG_PAXOS_ACCEPT_REPLY = 2
	REPLICA_MSG_PAXOS_COMMIT_REQ   = 3
	//REPLICA_MSG_PAXOS_COMMIT_REPLY = 4
	REPLICA_MSG_FP_VOTE         = 11
	REPLICA_MSG_FP_ACCEPT_REQ   = 12
	REPLICA_MSG_FP_ACCEPT_REPLY = 13
	REPLICA_MSG_FP_COMMIT_REQ   = 14
	//REPLICA_MSG_FP_COMMIT_REPLY    = 15
)

const (
	ENTRY_STAT_COMMITTED         = 1
	ENTRY_STAT_LEADER_ACCEPTED   = 2
	ENTRY_STAT_ACCEPTOR_ACCEPTED = 3 // the fast path in Fast Paxos
)

const (
	DEFAULT_PAXOS_LOG_INIT_SIZE = 1024 * 32
)

const (
	DEFAULT_CHANNEL_BUFFER_SIZE = 10240 * 8
)

const (
	// Time to wait for other servers to init before sending the first heart beat
	DEFAULT_FIRST_HEART_BEAT_DELAY = 5 * 1000 * 1000 * 1000 // ns
)
