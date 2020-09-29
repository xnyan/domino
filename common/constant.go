package common

// Configuration flags
const (
	Flag_leader        = "leader"
	Flag_follower_list = "follower.list"
	Flag_dc_list       = "dc.list"

	Flag_replica_num  = "replica.num"
	Flag_majority_num = "majority.num"
	Flag_fast_quorum  = "fast.quorum"

	Flag_stats_enabled = "stats.enabled"

	Flag_rpc_stream_server_client_buffer_size = "rpc.stream.server.client.buffer.size"

	Flag_fastpaxos_vote_ch_buffer_size    = "fastpaxos.vote.channel.buffer.size"
	Flag_fastpaxos_command_ch_buffer_size = "fastpaxos.command.channel.buffer.size"
	Flag_fastpaxos_exec_ch_buffer_size    = "fastpaxos.execution.channel.buffer.size"

	Flag_fastpaxos_log_type            = "fastpaxos.log.type"
	Flag_fastpaxos_log_manager         = "fastpaxos.log.manager"
	Flag_fastpaxos_log_size            = "fastpaxos.log.size"
	Flag_fastpaxos_log_file_path       = "fastpaxos.log.file.path"
	Flag_fastpaxos_log_persist_enabled = "fastpaxos.log.persist.enabled"

	Flag_fastpaxos_scheduler         = "fastpaxos.scheduler"
	Flag_fastpaxos_scheduler_window  = "fastpaxos.scheduler.window" // timestamp process window
	Flag_fastpaxos_delay_config_file = "fastpaxos.delay.config.file"
	Flag_fastpaxos_delay_config_tag  = "fastpaxos.delay.config.tag" // oneway-delay
	Flag_fastpaxos_delay_additional  = "fastpaxos.delay.additional" // e.g., "5ms"

	// "all" or "closest"(closest quorum)
	Flag_fastpaxos_delay_quorum_type = "fastpaxos.delay.quorum.type"

	Flag_fastpaxos_do_noop_proposer_num   = "fastpaxos.noop.proposer.num"
	Flag_fastpaxos_do_noop_ch_buffer_size = "fastpaxos.noop.channel.buffer.size"

	//Execution
	////If enable clients to wait for execution result. Default is false.
	Flag_fastpaxos_exec_mode    = "fastpaxos.exec.mode"
	Default_fastpaxos_exec_mode = "false"

	//data
	Flag_data_key_file    = "data.key.file"
	Flag_data_val_size    = "data.val.size"
	Default_data_val_size = "1024"
)

// Log
const (
	FixedLog          = "fixed"
	SegLog            = "seg"
	DefaultLogManager = "default"
	//DeterministicLogManager = "deterministic"
)

// Constants
const (
	FastpaxosLogFileSuffix = "-fastpaxos.log"
)

// Scheduler type
const (
	NoScheduler            = "no"
	DelayScheduler         = "delay"
	TimestampScheduler     = "timestamp"
	DeterministicScheduler = "deterministic"
)

// Delay calculation
const (
	DelayAll           = "all"
	DelayClosestQuorum = "closest"
)
