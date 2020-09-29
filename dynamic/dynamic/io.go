package dynamic

import ()

type ReplicaIo interface {
	InitConn(addrList []string, streamBufSize int)

	// A replica broadcasts heart beat to others
	BcstHb(addrList []string, nat int64, rId string, fpExecT *Timestamp)

	// Paxos shard leader
	BcstPaxosAcceptReq(addrList []string, t *Timestamp, cmd *Command, rId string, fpExecT *Timestamp)
	BcstPaxosFutureAcceptReq(addrList []string, t *Timestamp, cmd *Command, rId string, nat int64, fpExecT *Timestamp)
	// Piggybacks current commit time
	BcstPaxosCommitReq(addrList []string, t *Timestamp, cmdId string, nat int64, rId string, fpExecT *Timestamp)

	// Paxos shard follower
	// Piggybacks current commit time
	SendPaxosAcceptReply(addr string, t *Timestamp, cmdId string, nat int64, rId string, fpExecT *Timestamp)

	// Fast Paxos shard leader
	BcstFpAccept(addrList []string, t *Timestamp, cmd *Command, execT *Timestamp, nat int64, rId string)
	BcstFpCommit(addrList []string, t *Timestamp, cmd *Command, isFast bool, execT *Timestamp, nat int64, rId string)

	// Fast Paxos all learners
	BcstFpFastVote(addrList []string, t *Timestamp, cmd *Command, isFastAccept bool, nat int64, rId string)

	// Fast Paxos shard follower
	SendFpFastVote(addr string, t *Timestamp, cmd *Command, isFastAccept bool, nat int64, rId string)
	SendFpAcceptReply(addr string, t *Timestamp, cmdId string, nat int64, rId string)

	// Blocking
	// Synchronous I/O Latency probing
	// TODO to be piggybacked on the heart beat message
	SyncSendReplicaProbeReq(addr string) *ReplicaProbeReply
}

// Not thread-safe
type StreamIo struct {
	rpcIo RpcIo
}

func NewStreamIo(
	isGrpc, isSyncSend bool, id string, num int, nodeAddrList []string, p Paxos,
) ReplicaIo {
	io := &StreamIo{}
	if isGrpc {
		io.rpcIo = NewGrpcIo(isSyncSend)
	} else {
		io.rpcIo = NewFastIo(isSyncSend, id, num, nodeAddrList, p)
	}
	return io
}

func (io *StreamIo) InitConn(addrList []string, streamBufSize int) {
	io.rpcIo.InitConn(addrList, streamBufSize)
}

// Leader fo followers
// Shard leaders broadcast heartbeat to followers
func (io *StreamIo) BcstHb(
	addrList []string, nat int64, rId string, fpExecT *Timestamp,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_HEART_BEAT,
		NonAcceptT: nat,
		ReplicaId:  rId,
		FpExecT:    fpExecT,
	}

	io.rpcIo.BcstReplicaMsg(addrList, msg)
}

// Shard leaders broadcast Paxos accept req to followers
func (io *StreamIo) BcstPaxosAcceptReq(
	addrList []string, t *Timestamp, cmd *Command, rId string, fpExecT *Timestamp,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_PAXOS_ACCEPT_REQ,
		Cmd:        cmd,
		Time:       t,
		NonAcceptT: t.Time, // Paxos shard does not accept new command before this timestamp
		ReplicaId:  rId,
		FpExecT:    fpExecT,
	}

	io.rpcIo.BcstReplicaMsg(addrList, msg)
}

func (io *StreamIo) BcstPaxosFutureAcceptReq(
	addrList []string, t *Timestamp, cmd *Command, rId string, nat int64, fpExecT *Timestamp,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_PAXOS_ACCEPT_REQ,
		Cmd:        cmd,
		Time:       t,
		NonAcceptT: nat, // Paxos shard does not accept new command before this timestamp
		ReplicaId:  rId,
		FpExecT:    fpExecT,
	}

	io.rpcIo.BcstReplicaMsg(addrList, msg)
}

// Shard leaders broadcast Paxos commit req to followers
func (io *StreamIo) BcstPaxosCommitReq(
	addrList []string, t *Timestamp, cmdId string, nat int64, rId string, fpExecT *Timestamp,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_PAXOS_COMMIT_REQ,
		Cmd:        &Command{Id: cmdId}, // The id is used to assert correctness
		Time:       t,
		NonAcceptT: nat,
		ReplicaId:  rId,
		FpExecT:    fpExecT,
	}

	io.rpcIo.BcstReplicaMsg(addrList, msg)
}

// Follower to leader
// A follower sends a Paxos accept reply to the shard leader
func (io *StreamIo) SendPaxosAcceptReply(
	addr string, t *Timestamp, cmdId string, nat int64, rId string, fpExecT *Timestamp,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_PAXOS_ACCEPT_REPLY,
		Cmd:        &Command{Id: cmdId},
		Time:       t,
		IsAccept:   true,
		NonAcceptT: nat,
		ReplicaId:  rId,
		FpExecT:    fpExecT,
	}

	io.rpcIo.SendReplicaMsg(addr, msg)
}

//// Fast Paxos
// A leader sends slow-path accept req to followers
func (io *StreamIo) BcstFpAccept(
	addrList []string, t *Timestamp, cmd *Command, execT *Timestamp, nat int64, rId string,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_FP_ACCEPT_REQ,
		Cmd:        cmd,
		Time:       t,
		FpExecT:    execT,
		NonAcceptT: nat,
		ReplicaId:  rId,
	}

	io.rpcIo.BcstReplicaMsg(addrList, msg)
}

// A leader sends fast-path/slow-path commit req to followers
func (io *StreamIo) BcstFpCommit(
	addrList []string, t *Timestamp, cmd *Command, isFast bool,
	execT *Timestamp, nat int64, rId string,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_FP_COMMIT_REQ,
		Cmd:        cmd,
		Time:       t,
		IsAccept:   isFast, // Re-use the field to indicate if this is a fast-path commit
		FpExecT:    execT,
		NonAcceptT: nat,
		ReplicaId:  rId,
	}

	io.rpcIo.BcstReplicaMsg(addrList, msg)
}

func (io *StreamIo) BcstFpFastVote(
	addrList []string, t *Timestamp, cmd *Command, isFastAccept bool, nat int64, rId string,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_FP_VOTE,
		Cmd:        cmd,
		Time:       t,
		IsAccept:   isFastAccept,
		NonAcceptT: nat,
		ReplicaId:  rId,
	}

	io.rpcIo.BcstReplicaMsg(addrList, msg)
}

// A Fast Paxos shard follower sends a fast-path vote to the leader
func (io *StreamIo) SendFpFastVote(
	addr string, t *Timestamp, cmd *Command, isFastAccept bool, nat int64, rId string,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_FP_VOTE,
		Cmd:        cmd,
		Time:       t,
		IsAccept:   isFastAccept,
		NonAcceptT: nat,
		ReplicaId:  rId,
	}

	io.rpcIo.SendReplicaMsg(addr, msg)
}

func (io *StreamIo) SendFpAcceptReply(
	addr string, t *Timestamp, cmdId string, nat int64, rId string,
) {
	msg := &ReplicaMsg{
		Type:       REPLICA_MSG_FP_ACCEPT_REPLY,
		Cmd:        &Command{Id: cmdId},
		Time:       t,
		IsAccept:   true,
		NonAcceptT: nat,
		ReplicaId:  rId,
	}

	io.rpcIo.SendReplicaMsg(addr, msg)
}

// Blocking
func (io *StreamIo) SyncSendReplicaProbeReq(addr string) *ReplicaProbeReply {
	return io.rpcIo.SendReplicaProbeReq(addr)
}
