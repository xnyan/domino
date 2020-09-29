package dynamic

import (
	"runtime/debug"
)

////Fast Paxos
func (dp *DynamicPaxos) handleFpProposal(p *FpProposal) {
	fpExecTM := dp.getFpShardExecTM(p.FpT.Shard)
	nat := dp.tm.GetCurrentTime()

	//logger.Infof("cmdId = %s, timestamp = %d, execT = %d nat = %d, over execT = %d, over nat = %d",
	//	p.Cmd.Id, p.FpT.Time, fpExecTM.GetExecT(), nat, p.FpT.Time-fpExecTM.GetExecT(), p.FpT.Time-nat)
	isAccept := false
	if p.FpT.Time > fpExecTM.GetExecT() {
		if nat <= p.FpT.Time {
			isAccept = dp.lm.FpFastAcceptCmd(p.FpT, p.Cmd)
		}
	}

	// Sends fast-path result to client
	p.RetC <- isAccept

	if !isAccept {
		// Does not vote rejected commands
		return
	}

	if dp.isFpLeaderLearner {
		// The FP shard leader is the sole learner in the consensus group.
		// Followers will wait for the leader's decisions on commit and execution (time).
		// This approach reduces message complexity. However, the execution time on
		// followers will have a delay of at least one-way message (from the leader
		// to followers), which would affect the execution time for the Paxos
		// shards coordinated by the followers.
		if dp.GetCurFpShard() == p.FpT.Shard {
			// FP Shard leader
			// Locally processes fast-path vote
			dp.processFpFastVote(p.FpT, p.Cmd, dp.getReplicaId())
			// Locally updates the fast-path non-accept time
			dp.processFpFastNonAcceptTime(nat, dp.getReplicaId(), dp.getReplicaId())
		} else {
			// FP Shard follower
			// Sends fast-path vote to FP shard leader, and piggybacks the not-accept timestamp
			//loger.Infof("sending nat = %d", nat)
			dp.fpVote(p.FpT, p.Cmd, isAccept, nat, dp.getReplicaId())
		}
	} else {
		// Broadcasts the vote to every other replica
		dp.bcstVote(p.FpT, p.Cmd, isAccept, nat, dp.getReplicaId())

		// Locally processes fast-path vote
		dp.processFpFastVote(p.FpT, p.Cmd, dp.getReplicaId())
		// Locally updates the fast-path non-accept time
		leaderId := dp.getFpShardLeaderId(p.FpT.Shard)
		dp.processFpFastNonAcceptTime(nat, leaderId, dp.getReplicaId())
	}

	// NOTE Alternative approach:
	// Broadcasts the vote every replica. This will allow each replica to
	// learn both the fast-path results (commit) and execution time at the same
	// time.  This approach could reduce execution latency for every replica
	// and the Paxos shards that these followers coordinate.
	// However, this approach increases message complexity and requires every
	// replica to do computation for the fast-path consensus results.
	// In this approach, if the fast path fails, each replica still depends on
	// the FP shard leader to make a decision.  Even if a command has a
	// majority of votes, a follower cannot make any decision. It has to wait
	// for the leader to make decisions. This is because if there are f
	// failures, a new leader could make a different decision. (NOTE: this may
	// be solved by allowing the new leader to collect all of the received
	// votes on followers, but this is different from standard Fast Paxos.)
}

// A Fast Paxos shard leader processes a fast-path accept vote
func (dp *DynamicPaxos) processFpFastVote(
	fpT *Timestamp, cmd *Command, rId string,
) {
	fpShard := fpT.Shard
	t := fpT.Time
	fpExecTM := dp.getFpShardExecTM(fpShard)
	nat := fpExecTM.GetFastNonAcceptT(rId)
	if t <= nat {
		logger.Fatalf("Replica %s fast accepts cmd = %s at t = %d passing non-accept time = %d",
			rId, cmd.Id, t, nat)
	}

	execT := fpExecTM.GetExecT()
	if t <= execT {
		// Ignores the accept vote since a command has been chosen.
		return
	}

	fpConsM := dp.getFpShardConsM(fpShard)
	cons := fpConsM.GetFpCons(t)
	n := cons.Accept(cmd)
	if n == dp.getFastQuorum() {
		cons.ChooseCmd(cmd)
		cons.SetFpCmdConsRet(cmd.Id)
		dp.SetFpCmdConsRet(t, cmd.Id, true, true)
		dp.fpRejectConflictCmd(cons)
		dp.fpCommitCmd(fpT, cmd, true)
	}

	if cons.IsCmdChosen() {
		if !cons.IsSetFpCmdConsRet(cmd.Id) {
			cons.SetFpCmdConsRet(cmd.Id)
			dp.SetFpCmdConsRet(t, cmd.Id, false, false)
		}
	}

	if cons.GetAcceptN() == dp.getReplicaNum() {
		// Only when all replicas vote a cmd on the slot
		if !cons.IsCmdChosen() {
			if dp.isFpLeaderLearner || dp.GetCurFpShard() == fpShard {
				slowCmd := cons.SelectCmd(dp.getMajority())
				cons.ChooseCmd(slowCmd) // This is not consensus result, just a flag
				fpConsM.AddSlowFpCons(cons)
				dp.MarkFpCmdSlowConsRet(cons.GetT(), cons.GetCmdIdMap())
				dp.fpAcceptCmd(fpT, slowCmd)
			}
		}
	}
}

// A Fast Paxos shard leader knows that a replica will not (fast-path) accept
// any commands <= the given timestamp, fast non-accept time. The leader can
// use this information to know a time (exec time) before which the consensus
// group will not choose any new commands.
//
// Before the leader sends the exec time to any replica, the leader must assure
// that the replica knows all timestamps that will choose a command.
func (dp *DynamicPaxos) processFpFastNonAcceptTime(nat int64, leaderId, rId string) {
	fpShard := dp.getFpShard(leaderId)
	fpExecTM := dp.getFpShardExecTM(fpShard)

	fpExecTM.UpdateNonAcceptTime(rId, nat)
	fqNat, _ := fpExecTM.GetFastQuorumNonAcceptT()
	minNat, _ := fpExecTM.GetMinNonAcceptT() // minT <= fqMinT

	fpConsM := dp.getFpShardConsM(fpShard)
	for cons, ok := fpConsM.PeekCons(); ok; cons, ok = fpConsM.PeekCons() {
		t := cons.GetT()

		if t > fqNat {
			// New command may be chosen before this command but after fqNat
			break
		}

		if cons.IsCmdChosen() {
			// The consensus instance has chosen a command (via the fast path or
			// the slow path).
			// A supermajority of replicas guarantee that they do not accept any
			// command before this instance, let's say it is time T.
			// This indicates that if all of the other replicas accept the same
			// command at a timestamp before T, the consenus group does not choose
			// that command because the supermajority has chosen a logical no-op for
			// the timestamp. That is, the no-op has the fast path succeed.
			fpConsM.PopCons()
			fpExecTM.UpdateExecT(t)
			dp.SetFpCmdConsRetExecT(t)
			continue
		}

		if t <= minNat {
			// The consensus instance has not chosen any command but all of replicas
			// guarantee that they do not accept any command before this instance.

			// The number of replicas that have fast-path accepted a cmd
			cmdN := cons.GetAcceptN()
			fpT := &Timestamp{Time: t, Shard: fpShard}
			if dp.getFastQuorum()+cmdN <= dp.getReplicaNum() {
				// The consensus is waiting for replies from at least a supermajority
				// of replicas, and these replicas guarantee they do not accept any
				// command for this instance.
				// Fast-path commits no-op
				dp.fpCommitCmd(fpT, nil, true)
				cons.ChooseCmd(nil)
				dp.fpRejectConflictCmd(cons)
			} else if dp.getMajority()+cmdN <= dp.getReplicaNum() {
				if dp.isFpLeaderLearner || dp.GetCurFpShard() == fpShard {
					// Slow-path commits a no-op
					cons.ChooseCmd(nil) // this is not a consensus result, just a flag
					fpConsM.AddSlowFpCons(cons)
					dp.MarkFpCmdSlowConsRet(cons.GetT(), cons.GetCmdIdMap())
					dp.fpAcceptCmd(fpT, nil)
				} else {
					ct := dp.tm.PrevT(t)
					fpExecTM.UpdateExecT(ct)
					dp.SetFpCmdConsRetExecT(ct)
					return
				}
			} else {
				// Does not choose a no-op here. Since Fast Paxos's recovery protocol
				// allows a random command, there is no need to choose a no-op.
				if dp.isFpLeaderLearner || dp.GetCurFpShard() == fpShard {
					slowCmd := cons.SelectCmd(dp.getMajority())
					cons.ChooseCmd(slowCmd) // this is not a consensus result, just a flag
					fpConsM.AddSlowFpCons(cons)
					dp.MarkFpCmdSlowConsRet(cons.GetT(), cons.GetCmdIdMap())
					dp.fpAcceptCmd(fpT, slowCmd)
				} else {
					ct := dp.tm.PrevT(t)
					fpExecTM.UpdateExecT(ct)
					dp.SetFpCmdConsRetExecT(ct)
					return
				}
			}

			fpConsM.PopCons()
			fpExecTM.UpdateExecT(t)
			dp.SetFpCmdConsRetExecT(t)
			continue
		} else {
			// The consenus instance has not chosen a command yet but falls beteen
			// minNat and fqNat.
			// At least a supermajority of replicas guarantees not to accept new
			// command before faNat >= t.Time. The group cannot choose a new command
			// before this instance. Since this instance has not chosen a command
			// yet, sets the commit time as the timestamp before it.
			ct := dp.tm.PrevT(t)
			fpExecTM.UpdateExecT(ct)
			dp.SetFpCmdConsRetExecT(ct)
			return

			// TODO Optimizations: if any instance can choose a command, and the
			// decision would not conflict any possible fast-path result or slow-path
			// result, we should choose the command here. For example, a majority or
			// supermajority of replicas have guaranteed not to accept any command
			// for this instance, we should commit the no-op for this instance (via
			// either the slow or fast path.) This optimization would make the commit
			// time move foward in order to reduce execution latency.
		}
	}

	// A supermajority of replicas guarantee not to accept any new command before fqNat
	// There is no consensus instance that has not chosen a command before fqNat.
	// Sets the commit time as the fqNat.
	fpExecTM.UpdateExecT(fqNat)
	dp.SetFpCmdConsRetExecT(fqNat)
}

func (dp *DynamicPaxos) processFpAcceptReq(t *Timestamp, cmd *Command) {
	dp.lm.FpAcceptCmd(t, cmd)

	lId := dp.getFpShardLeaderId(t.Shard)
	l := dp.getReplicaAddr(lId)

	nat := dp.tm.GetCurrentTime()
	//loger.Infof("sending nat = %d", nat)
	if cmd == nil {
		dp.io.SendFpAcceptReply(l, t, NO_OP_ID, nat, dp.getReplicaId())
	} else {
		dp.io.SendFpAcceptReply(l, t, cmd.Id, nat, dp.getReplicaId())
	}
}

func (dp *DynamicPaxos) processFpAcceptReply(t *Timestamp, cmdId string) {
	fpSlowCons := dp.getFpSlowCons(t.Time)
	n := fpSlowCons.VoteAccept()

	if n == dp.Majority {
		fpConsM := dp.getFpShardConsM(t.Shard)
		fpCons := fpConsM.GetSlowFpCons(t.Time)
		if cmdId == NO_OP_ID {
			dp.fpCommitCmd(t, nil, false)
		} else {
			// Marks the consensus result for the command
			fpCons.SetFpCmdConsRet(cmdId)
			dp.SetFpCmdConsRet(t.Time, cmdId, true, false)
			// Commits the cmd via the slow path
			cmd := &Command{Id: cmdId}
			dp.fpCommitCmd(t, cmd, false)
		}
		fpConsM.DelSlowFpCons(t.Time)
		dp.fpRejectConflictCmd(fpCons)
		dp.DelFpCmdSlowConsRet(fpCons.GetT())
	}

	if n == dp.ReplicaNum {
		dp.delFpSlowCons(t.Time)
	}
}

// TODO when executes commands, a replica may have a pending fast-path
// consensus instance because it fast accepts a command but its vote is
// rejected due to passing the exec time
func (dp *DynamicPaxos) processFpCommitReq(t *Timestamp, cmd *Command, isFast bool) {
	dp.lm.FpCommitCmd(t, cmd, isFast)

	if !dp.isFpLeaderLearner {
		// Follower
		fpConsM := dp.getFpShardConsM(t.Shard)
		fpCons := fpConsM.GetFpCons(t.Time)
		fpCons.ChooseCmd(cmd)
		if cmd != nil && cmd.Id != NO_OP_ID {
			// Marks the consensus result for the command
			fpCons.SetFpCmdConsRet(cmd.Id)
			dp.SetFpCmdConsRet(t.Time, cmd.Id, true, false)
		}
		dp.fpRejectConflictCmd(fpCons)
	}
}

func (dp *DynamicPaxos) fpVote(
	t *Timestamp, cmd *Command, isFastAccept bool, nat int64, rId string) {
	lId := dp.getFpShardLeaderId(t.Shard)
	l := dp.getReplicaAddr(lId)
	if isFastAccept {
		dp.io.SendFpFastVote(l, t, cmd, isFastAccept, nat, rId)
	} else {
		dp.io.SendFpFastVote(l, t, &Command{Id: cmd.Id}, isFastAccept, nat, rId)
	}
}

func (dp *DynamicPaxos) bcstVote(
	t *Timestamp, cmd *Command, isFastAccept bool, nat int64, rId string) {
	if isFastAccept {
		dp.io.BcstFpFastVote(dp.followerAddrList, t, cmd, isFastAccept, nat, rId)
	} else {
		dp.io.BcstFpFastVote(dp.followerAddrList, t, &Command{Id: cmd.Id}, isFastAccept, nat, rId)
	}
}

func (dp *DynamicPaxos) fpAcceptCmd(t *Timestamp, cmd *Command) {
	fpSlowCons := dp.initFpSlowCons(t.Time)

	// Piggybacks the Fast Paxos shard execution time
	execT := dp.getFpExecT(t.Shard)
	// Piggybacks non-accept time (only for Paxos shard)
	nat := dp.tm.GetCurrentTime()
	//loger.Infof("sending nat = %d", nat)
	// Asks every follower to accept the cmd at the time
	dp.io.BcstFpAccept(dp.followerAddrList, t, cmd, execT, nat, dp.getReplicaId())

	// Locally accepts
	dp.lm.FpAcceptCmd(t, cmd)
	fpSlowCons.VoteAccept()

	// Locally handles execution timestamp
	dp.fpReplicaHandleLeaderExecT(nat, execT, dp.getReplicaId())
}

func (dp *DynamicPaxos) fpCommitCmd(t *Timestamp, cmd *Command, isFast bool) {
	// Piggybacks the Fast Paxos shard execution time
	execT := dp.getFpExecT(t.Shard)
	// Piggybacks non-accept time (only for Paxos shard)
	nat := dp.tm.GetCurrentTime()
	//loger.Infof("sending nat = %d", nat)

	if dp.isFpLeaderLearner || isFast == false {
		// Asks followers to commit the command
		dp.io.BcstFpCommit(dp.followerAddrList, t, cmd, isFast, execT, nat, dp.getReplicaId())
	}

	// Locally commits
	dp.lm.FpCommitCmd(t, cmd, isFast)

	// Locally handles execution timestamp
	dp.fpReplicaHandleLeaderExecT(nat, execT, dp.getReplicaId())
}

//func (dp *DynamicPaxos) getFpExecT() *Timestamp {
//	fpExecTM := dp.getFpShardExecTM(dp.GetCurFpShard())
//	execT := &Timestamp{Time: fpExecTM.GetExecT(), Shard: dp.GetCurFpShard()}
//	return execT
//}
func (dp *DynamicPaxos) getFpExecT(shard int32) *Timestamp {
	fpExecTM := dp.getFpShardExecTM(shard)
	execT := &Timestamp{Time: fpExecTM.GetExecT(), Shard: shard}
	return execT
}

// Fast Paxos helper functions
func (dp *DynamicPaxos) getFpShardLeaderId(shard int32) string {
	l, ok := dp.fpShardLeaderIdMap[shard]
	if !ok {
		logger.Fatalf("Cannot find the leader addr for fast paxos shard = %d", shard)
	}
	return l
}

func (dp *DynamicPaxos) getFpShard(leaderId string) int32 {
	shard, ok := dp.fpLeaderIdShardMap[leaderId]
	if !ok {
		debug.PrintStack()
		logger.Fatalf("Cannot find the Fast Paxos shard for leader = %s", leaderId)
	}
	return shard
}

// Constructs and returns the Fast paxos shard' non-accept timestamp
func (dp *DynamicPaxos) getFpShardNat(t int64, leaderId string) *Timestamp {
	shard := dp.getFpShard(leaderId)
	return &Timestamp{Time: t, Shard: shard}
}

func (dp *DynamicPaxos) getFpShardConsM(fpShard int32) *FpConsManager {
	if m, ok := dp.fpShardConsMap[fpShard]; ok {
		return m
	}
	logger.Fatalf("Cannot find fast paxos consensus manager for shard = %d", fpShard)
	return nil
}

func (dp *DynamicPaxos) getFpShardExecTM(fpShard int32) *FpExecTimeManager {
	if m, ok := dp.fpShardExecTMap[fpShard]; ok {
		return m
	}
	logger.Fatalf("Cannot find ExecTimeManager for fast paxos shard = %d", fpShard)
	return nil
}

func (dp *DynamicPaxos) initFpSlowCons(t int64) *PaxosCons {
	fpCons := NewPaxosCons(nil)
	dp.fpSlowConsMap[t] = fpCons
	return fpCons
}

func (dp *DynamicPaxos) getFpSlowCons(t int64) *PaxosCons {
	fpCons, ok := dp.fpSlowConsMap[t]
	if !ok {
		logger.Fatalf("Cannot find Paxos Consenus Instance for t = %d", t)

	}
	return fpCons
}

func (dp *DynamicPaxos) delFpSlowCons(t int64) {
	delete(dp.fpSlowConsMap, t)
}

// Fast Paxos consenus result for the shard leader to send to the client
// Marks consensus results for commands that are determined to reject on the time slot
func (dp *DynamicPaxos) fpRejectConflictCmd(fpCons *FpCons) {
	if !fpCons.IsCmdChosen() {
		logger.Fatalf("Should not begin rejection until a command is chosen. t = %d", fpCons.GetT())
	}

	t := fpCons.GetT()
	cmdIdList := fpCons.GetRejectCmdIdList()
	for _, cmdId := range cmdIdList {
		fpCons.SetFpCmdConsRet(cmdId)
		dp.SetFpCmdConsRet(t, cmdId, false, false)
	}
}

// Blocking until the consensus result is known
func (dp *DynamicPaxos) FpWaitConsensus(t int64, cmdId string) (bool, bool) {
	cmdCons := dp.cmdM.getFpCmdCons(t, cmdId, true)
	ret := <-cmdCons.c
	dp.cmdM.delFpCmdConsRetCh(t, cmdId)
	return ret.IsAccept, ret.IsFast
}

func (dp *DynamicPaxos) SetFpCmdConsRet(t int64, cmdId string, isAccept, isFast bool) {
	dp.cmdM.inputCh <- &FpCmdConsRet{T: t, CmdId: cmdId, IsAccept: isAccept, IsFast: isFast}
}

func (dp *DynamicPaxos) MarkFpCmdSlowConsRet(t int64, cmdIdMap map[string]bool) {
	dp.cmdM.inputCh <- &FpCmdSlowConsRet{T: t, CmdIdMap: cmdIdMap}
}

func (dp *DynamicPaxos) DelFpCmdSlowConsRet(t int64) {
	dp.cmdM.inputCh <- &FpCmdSlowConsRet{T: t, CmdIdMap: nil}
}

func (dp *DynamicPaxos) SetFpCmdConsRetExecT(t int64) {
	dp.cmdM.inputCh <- t
}
