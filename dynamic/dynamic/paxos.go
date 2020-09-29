package dynamic

import (
	"fmt"
	"math"
	"time"
)

type Paxos interface {
	// Paxos assignes a request with a future time when the replication completes
	EnablePaxosFutureTime(s *Server)

	Start(hbDelay time.Duration)

	// Probing message for latency monitoring
	// Returns the queuing (i.e., processing delay)
	Probe() time.Duration

	// Returns isCommitted
	PaxosLeaderAccept(cmd *Command) bool

	// Returns isAccepted (A cmd gets rejected if the timestamp passes the current time)
	FpReplicaAccept(cmd *Command, fpT *Timestamp) bool

	// Processes messages between replicas
	ProcessReplicaMsg(msg *ReplicaMsg)

	// Fast Paxos shard leader waits for consensus result
	FpWaitConsensus(t int64, cmdId string) (bool, bool)

	// Returns the Paxos shard ID that this replica is the leader
	GetCurPaxosShard() int32

	// Returns the Fast Paxos shard ID that this replica is the leader
	GetCurFpShard() int32

	// Returns the command execution channel for applications
	GetExecCh() <-chan Entry

	// Returns RPC I/O handle TODO Decouples I/O from Paxos module
	GetIo() ReplicaIo

	// Testing
	Test()
}

type DynamicPaxos struct {
	// Input channel
	cmdCh chan interface{}
	//pCh   chan *PaxosProposal
	//fpCh  chan *FpProposal
	//msgCh chan *ReplicaMsg
	//latCh chan *LatReq

	// Heart beat timer
	hbTimer    *time.Timer
	hbInterval time.Duration

	// Paxos consensus control
	pShard            int32
	pShardLeaderIdMap map[int32]string
	pLeaderIdShardMap map[string]int32
	pConsMap          map[int64]*PaxosCons

	// Fast Paxos consensus control
	replicaId          string
	fpShard            int32
	fpShardLeaderIdMap map[int32]string
	fpLeaderIdShardMap map[string]int32
	fpShardConsMap     map[int32]*FpConsManager
	fpShardExecTMap    map[int32]*FpExecTimeManager
	fpSlowConsMap      map[int64]*PaxosCons // NOTE: this is not sharded yet
	// Fast Paxos All Learners, execution latency
	isFpLeaderLearner bool

	//Replica info
	ReplicaNum       int
	Majority         int
	FastQuorum       int
	followerAddrList []string
	replicaIdAddrMap map[string]string

	cmdM *FpCmdConsRetManager
	tm   TimeManager // Time / clock manager
	lm   LogManager  // log manager
	io   ReplicaIo   // FIFO io

	cmdChSize int

	// Network I/O swtich
	IsGrpc     bool // If using gRPC
	IsSyncSend bool // If doing server network I/O and computation in the same thread

	//// Paxos uses a future timestamp for reducing execution latency.
	// When this is enabled, instead of using the current timestamp, a Paxos
	// leader assigns a future timestamp that indicates the expected time that
	// the replication should complete.  A committed instance (Paxos or Fast
	// Paxos) instance would have to wait for at most one-way message for any
	// previous Paxos instances to commit. This would reduce execution latency
	// since in Mencius an instance may have to wait one roundtrip to execute.
	// Paxos as well as in Fast Paxos.
	// NOTE: The future timestamp is calculated
	// based on dynamic latency prediction. A later request could have an smaller
	// timestamp than an earlier request's timestamp. This will invalidate the
	// FIFO requirement of Mencius, where a later request should be put in a log
	// position after any earlier accepted requests. To solve this problem, we
	// track the largest future timestamp that has been assigned, and we need to
	// guarantee that assigned timestamps increase.
	isPaxosFutureTime bool
	lastPaxosT        int64   // the last assigned Paxos timestamp
	server            *Server // handle to the server for accessing latency prediction information

	// Testing
	testCh chan bool
}

func NewDynamicPaxos(
	replicaId string,
	replicaIdAddrMap map[string]string,
	rIdList, nodeAddrList []string, // a list of replicIds and Addrs
	pShard, fpShard int32,
	pShardLeaderIdMap, fpShardLeaderIdMap map[int32]string,
	followerAddrList []string,
	hbInterval time.Duration,
	cmdChSize int,
	execChSize int,
	isGrpc, isSyncSend bool,
	isFpLeaderLearner bool,
) Paxos {
	p := &DynamicPaxos{
		cmdCh: make(chan interface{}, cmdChSize*4),
		//pCh:   make(chan *PaxosProposal, cmdChSize),
		//fpCh:  make(chan *FpProposal, cmdChSize),
		//msgCh: make(chan *ReplicaMsg, cmdChSize*4),
		//latCh: make(chan *LatReq, cmdChSize),

		hbInterval: hbInterval,

		pShard:            pShard,
		pShardLeaderIdMap: pShardLeaderIdMap,
		pLeaderIdShardMap: make(map[string]int32),
		pConsMap:          make(map[int64]*PaxosCons, 5000),

		replicaId:          replicaId,
		fpShard:            fpShard,
		fpShardLeaderIdMap: fpShardLeaderIdMap,
		fpLeaderIdShardMap: make(map[string]int32),
		fpShardConsMap:     make(map[int32]*FpConsManager),
		fpShardExecTMap:    make(map[int32]*FpExecTimeManager),
		fpSlowConsMap:      make(map[int64]*PaxosCons, 5000),
		isFpLeaderLearner:  isFpLeaderLearner,

		followerAddrList: followerAddrList,
		replicaIdAddrMap: replicaIdAddrMap,

		cmdChSize: cmdChSize,

		IsGrpc:     isGrpc,
		IsSyncSend: isSyncSend,

		isPaxosFutureTime: false,

		testCh: make(chan bool, 1),
	}

	p.ReplicaNum = len(p.followerAddrList) + 1
	f := (p.ReplicaNum - 1) / 2
	p.Majority = f + 1
	p.FastQuorum = int(math.Ceil((3.0*float64(f))/2.0)) + 1

	pShardList, fpShardList := make([]int32, 0), make([]int32, 0)
	for shard, leaderId := range p.pShardLeaderIdMap {
		pShardList = append(pShardList, shard)
		p.pLeaderIdShardMap[leaderId] = shard
		logger.Infof("Paxos shard leader %s -> shard %d", leaderId, shard)
	}

	//rIdList := make([]string, 0, len(followerAddrList)+1)
	//rIdList = append(rIdList, p.replicaId)
	//rIdList = append(rIdList, followerAddrList...)
	for fpShard, leaderId := range p.fpShardLeaderIdMap {
		fpShardList = append(fpShardList, fpShard)
		p.fpLeaderIdShardMap[leaderId] = fpShard
		p.fpShardConsMap[fpShard] = NewFpConsManager(cmdChSize)
		p.fpShardExecTMap[fpShard] = NewFpExecTimeManager(p.ReplicaNum, p.FastQuorum, rIdList)
	}

	p.cmdM = NewFpCmdConsRetManager(cmdChSize)
	p.tm = NewRealClockTm()
	p.lm = NewSimpleLogManager(pShardList, fpShardList, execChSize)
	p.io = NewStreamIo(p.IsGrpc, p.IsSyncSend, p.replicaId, p.ReplicaNum, nodeAddrList, p)

	return p
}

// Enables Paxos to assign a request with a future time indicating when the
// replicaiton would complete. This is used to reduce execution latency.
func (dp *DynamicPaxos) EnablePaxosFutureTime(s *Server) {
	dp.isPaxosFutureTime = true
	dp.server = s
	dp.lastPaxosT = dp.tm.GetCurrentTime()
}

func (dp *DynamicPaxos) Probe() time.Duration {
	start := time.Now()
	req := &LatReq{
		RetC: make(chan bool, 1),
	}
	//dp.latCh <- req
	dp.cmdCh <- req

	<-req.RetC // blocking

	return time.Since(start) // queuing delay
}

func (dp *DynamicPaxos) PaxosLeaderAccept(cmd *Command) bool {
	p := &PaxosProposal{
		Cmd:  cmd,
		RetC: make(chan bool, 1),
	}

	//dp.pCh <- p
	dp.cmdCh <- p
	isCommit := <-p.RetC

	return isCommit
}

func (dp *DynamicPaxos) FpReplicaAccept(cmd *Command, fpT *Timestamp) bool {
	p := &FpProposal{
		Cmd:  cmd,
		FpT:  fpT,
		RetC: make(chan bool, 1),
	}

	//dp.fpCh <- p
	dp.cmdCh <- p
	isFastAccept := <-p.RetC

	return isFastAccept
}

func (dp *DynamicPaxos) ProcessReplicaMsg(msg *ReplicaMsg) {
	//dp.msgCh <- msg // FIFO
	dp.cmdCh <- msg // FIFO
}

func (dp *DynamicPaxos) GetCurPaxosShard() int32 {
	return dp.pShard
}

func (dp *DynamicPaxos) GetCurFpShard() int32 {
	return dp.fpShard
}

func (dp *DynamicPaxos) GetExecCh() <-chan Entry {
	return dp.lm.GetExecCh()
}

func (dp *DynamicPaxos) GetIo() ReplicaIo {
	return dp.io
}

// Non-blocking
func (dp *DynamicPaxos) Start(hbDelay time.Duration) {
	// Wait for a while to send the first heart beat so that all servers have started up
	dp.hbTimer = time.NewTimer(hbDelay)
	go func() {
		// Use a seperate thread to get the timeout signal but process the signal
		// in the main loop
		for {
			t := <-dp.hbTimer.C
			dp.cmdCh <- t
		}
	}()
	go func() {
		dp.io.InitConn(dp.followerAddrList, dp.cmdChSize*4)
		for {
			select {
			// Use one command channel to predict queuing/processing latency
			case cmd := <-dp.cmdCh:
				switch cmd.(type) {
				case time.Time:
					dp.handleHb()
				case *PaxosProposal:
					if dp.isPaxosFutureTime {
						dp.handlePaxosProposalWithFutureTime(cmd.(*PaxosProposal))
					} else {
						dp.handlePaxosProposal(cmd.(*PaxosProposal))
					}
				case *FpProposal:
					dp.handleFpProposal(cmd.(*FpProposal))
				case *ReplicaMsg:
					dp.handleReplicaMsg(cmd.(*ReplicaMsg))
				case *LatReq:
					dp.handleLatReq(cmd.(*LatReq))
				default:
					logger.Fatalf("")
				}

			case <-dp.testCh:
				dp.handleTest()
				/*
					// Using different channels
					case <-dp.hbTimer.C:
						dp.handleHb()

					case p := <-dp.pCh:
						dp.handlePaxosProposal(p)

					case p := <-dp.fpCh:
						dp.handleFpProposal(p)

					case msg := <-dp.msgCh:
						dp.handleReplicaMsg(msg)

					case req := <-dp.latCh:
						dp.handleLatReq(req)
				*/

			}
		}
	}()
}

//// Paxos Shard Leader
func (dp *DynamicPaxos) handlePaxosProposal(p *PaxosProposal) {
	// Assigns a timestamp for the cmd
	t := &Timestamp{Time: dp.tm.GetCurrentTime(), Shard: dp.GetCurPaxosShard()}

	// Inits a consensus instance and waits for Paxos replicas to accept the cmd
	pCons := dp.initPaxosCons(t.Time, p.RetC)

	// Piggybacks Fast Paxos shard execution time if this is a Fast Paxos shard leader
	var fpExecT *Timestamp = nil
	if dp.GetCurFpShard() >= 0 {
		fpExecT = dp.getFpExecT(dp.GetCurFpShard())
	}
	//loger.Infof("sending nat = %d", t.Time)
	// Sends accept request to followers, which piggybacks the last commit time
	dp.io.BcstPaxosAcceptReq(dp.followerAddrList, t, p.Cmd, dp.getReplicaId(), fpExecT)

	// Locally accepts the command
	dp.lm.PaxosAcceptCmd(t, p.Cmd)
	pCons.VoteAccept()

	// Locally updates the Paxos shard non-accept time
	dp.updatePaxosExecT(t.Time, dp.getReplicaId())
	dp.exec()
}

//// Paxos with using a future timestamp that indicates when replication completes
func (dp *DynamicPaxos) handlePaxosProposalWithFutureTime(p *PaxosProposal) {
	// Assigns a timestamp for the cmd
	paxosDelay := dp.server.PredictPaxosLat() * 1000000 // in ns
	curTime := dp.tm.GetCurrentTime()
	futureTime := curTime + paxosDelay // Assings a future time when the replication expects to complete
	// Makes sure that the assigned future timestamps monotonically increase
	if futureTime <= dp.lastPaxosT {
		futureTime = dp.lastPaxosT + 1000 // TODO configurable / adaptive increments
	}
	dp.lastPaxosT = futureTime
	futureT := &Timestamp{Time: futureTime, Shard: dp.GetCurPaxosShard()}

	// Inits a consensus instance and waits for Paxos replicas to accept the cmd
	pCons := dp.initPaxosCons(futureT.Time, p.RetC)

	// Piggybacks Fast Paxos shard execution time if this is a Fast Paxos shard leader
	var fpExecT *Timestamp = nil
	if dp.GetCurFpShard() >= 0 {
		fpExecT = dp.getFpExecT(dp.GetCurFpShard())
	}
	//loger.Infof("sending nat = %d", t.Time)
	// Sends accept request to followers, which piggybacks the last commit time
	dp.io.BcstPaxosFutureAcceptReq(dp.followerAddrList, futureT, p.Cmd, dp.getReplicaId(), curTime, fpExecT)

	// Locally accepts the command
	dp.lm.PaxosAcceptCmd(futureT, p.Cmd)
	pCons.VoteAccept()

	// Locally updates the Paxos shard non-accept time
	dp.updatePaxosExecT(curTime, dp.getReplicaId())
	dp.exec()
}

// Sends heart beat to every body
func (dp *DynamicPaxos) handleHb() {
	// Non-accept time
	t := dp.tm.GetCurrentTime()

	var fpExecT *Timestamp = nil
	if dp.isFpLeaderLearner {
		if dp.GetCurFpShard() >= 0 {
			// Fast Paxos shard leader
			dp.processFpFastNonAcceptTime(t, dp.getReplicaId(), dp.getReplicaId())
			fpExecT = dp.getFpExecT(dp.GetCurFpShard())
			if dp.updateFpExecT(fpExecT) {
				dp.prepareShardExec(fpExecT.Shard)
			}
		}
	} else {
		// TODO Makes sure that there is only one FP shard
		for _, fpShardLeader := range dp.fpShardLeaderIdMap {
			dp.processFpFastNonAcceptTime(t, fpShardLeader, dp.getReplicaId())
		}
		if dp.GetCurFpShard() >= 0 { // Only leader sends the FP exec time
			fpExecT = dp.getFpExecT(dp.GetCurFpShard())
			if dp.updateFpExecT(fpExecT) {
				dp.prepareShardExec(fpExecT.Shard)
			}
		}
	}

	// Re-calculates the non-accept time  for heart beat. This is because a FP
	// shard leader may have already sent out a new non-accept time while upding
	// the fast-paxos non-accept time. The new one is larger than the non-accept
	// time generated in the beginning of the heartbeat.
	// Therefore, have to re-generate the non-accept time here since other replicas
	// will receive this new one later than the new one. // TODO use a better solution
	t = dp.tm.GetCurrentTime()
	if dp.GetCurPaxosShard() >= 0 {
		// Paxos Shard leader
		dp.updatePaxosExecT(t, dp.getReplicaId())
		dp.prepareShardExec(dp.GetCurPaxosShard())
	}

	//loger.Infof("sending hb nat = %d", t)
	// Sends non-accept time and Fast Paxos execution time to other replicas
	dp.io.BcstHb(dp.followerAddrList, t, dp.getReplicaId(), fpExecT)

	dp.exec()

	dp.resetHb(dp.hbInterval)
}

func (dp *DynamicPaxos) handleReplicaMsg(msg *ReplicaMsg) {
	switch msg.Type {
	//// Paxos
	case REPLICA_MSG_PAXOS_ACCEPT_REQ:
		dp.handlePaxosAcceptReqMsg(msg)

	case REPLICA_MSG_PAXOS_ACCEPT_REPLY:
		dp.handlePaxosAcceptReplyMsg(msg)

	case REPLICA_MSG_PAXOS_COMMIT_REQ:
		dp.handlePaxosCommitReqMsg(msg)

	//// Fast Paxos
	case REPLICA_MSG_FP_VOTE:
		dp.handleFpVote(msg)

	case REPLICA_MSG_FP_ACCEPT_REQ:
		dp.handleFpAcceptReqMsg(msg)

	case REPLICA_MSG_FP_ACCEPT_REPLY:
		dp.handleFpAcceptReplyMsg(msg)

	case REPLICA_MSG_FP_COMMIT_REQ:
		dp.handleFpCommitReqMsg(msg)

	//// Heart beat for both Paxos and Fast Paxos
	case REPLICA_MSG_HEART_BEAT:
		dp.handleReplicaHbMsg(msg)

	default:
		logger.Fatalf("Undefined replica message type: %v", msg.Type)
	}
}

// Handles a heart beat message from a replica
func (dp *DynamicPaxos) handleReplicaHbMsg(msg *ReplicaMsg) {
	//loger.Infof("handling hb nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)

	// Updates Paxos shard non-accept time
	if pShard, ok := dp.getPaxosShard(msg.ReplicaId); ok {
		dp.updatePaxosExecT(msg.NonAcceptT, msg.ReplicaId)
		dp.prepareShardExec(pShard)
	}

	if dp.isFpLeaderLearner {
		if dp.GetCurFpShard() >= 0 {
			// Fast Paxos leader
			dp.processFpFastNonAcceptTime(msg.NonAcceptT, dp.getReplicaId(), msg.ReplicaId)
			fpExecT := dp.getFpExecT(dp.GetCurFpShard())
			if dp.updateFpExecT(fpExecT) {
				dp.prepareShardExec(fpExecT.Shard)
			}
		}
	} else {
		for fpShard, fpShardLeader := range dp.fpShardLeaderIdMap {
			dp.processFpFastNonAcceptTime(msg.NonAcceptT, fpShardLeader, msg.ReplicaId)
			fpExecT := dp.getFpExecT(fpShard)
			if dp.updateFpExecT(fpExecT) {
				dp.prepareShardExec(fpExecT.Shard)
			}
		}
	}

	if msg.FpExecT != nil {
		if dp.updateFpExecT(msg.FpExecT) {
			dp.prepareShardExec(msg.FpExecT.Shard)
		}
	}

	dp.exec()
}

// A Paxos shard leader handles the accept ack from a follower
func (dp *DynamicPaxos) handlePaxosAcceptReplyMsg(msg *ReplicaMsg) {
	//loger.Infof("handling nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)
	dp.processPaxosAcceptReply(msg.Time, msg.Cmd.Id)
	dp.paxosFollowerHandleLeaderNat(false, msg.NonAcceptT, msg.ReplicaId, msg.FpExecT)
}

// A Paxos shard follower handles the accept request from the leader
func (dp *DynamicPaxos) handlePaxosAcceptReqMsg(msg *ReplicaMsg) {
	//loger.Infof("handling nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)
	dp.processPaxosAcceptReq(msg.Time, msg.Cmd)
	dp.paxosFollowerHandleLeaderNat(false, msg.NonAcceptT, msg.ReplicaId, msg.FpExecT)
}

// Paxos shard follower handles the commit request from the leader
func (dp *DynamicPaxos) handlePaxosCommitReqMsg(msg *ReplicaMsg) {
	//loger.Infof("handling nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)
	dp.processPaxosCommitReq(msg.Time, msg.Cmd.Id)
	dp.paxosFollowerHandleLeaderNat(true, msg.NonAcceptT, msg.ReplicaId, msg.FpExecT)
}

func (dp *DynamicPaxos) paxosFollowerHandleLeaderNat(
	isCommit bool, nat int64, rId string, fpExecT *Timestamp,
) {
	// Updates Paxos shard non-accept time
	if pShard, ok := dp.getPaxosShard(rId); ok {
		dp.updatePaxosExecT(nat, rId)
		if isCommit {
			//dp.prepareShardExec(dp.getPaxosShard(rId))
			dp.prepareShardExec(pShard)
		}
	}

	// Updates Fast Paxos shard executon time if available
	if fpExecT != nil {
		if dp.updateFpExecT(fpExecT) {
			dp.prepareShardExec(fpExecT.Shard)
		}
	}

	dp.exec()
}

////Fast Paxos
////Fast Paxos leaders
// A Fast Paxos shard leader handles the accept vote from a follower
func (dp *DynamicPaxos) handleFpVote(msg *ReplicaMsg) {
	//loger.Infof("handling nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)
	dp.processFpFastVote(msg.Time, msg.Cmd, msg.ReplicaId)
	dp.fpLeaderHandleReplicaNat(msg.NonAcceptT, msg.ReplicaId)
}

// A Fast Paxos shard leader handles the accept ack from a follower
func (dp *DynamicPaxos) handleFpAcceptReplyMsg(msg *ReplicaMsg) {
	//loger.Infof("handling nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)
	dp.processFpAcceptReply(msg.Time, msg.Cmd.Id)
	dp.fpLeaderHandleReplicaNat(msg.NonAcceptT, msg.ReplicaId)
}

func (dp *DynamicPaxos) fpLeaderHandleReplicaNat(nat int64, rId string) {
	// Updates the piggybacked Paxos shard non-accept time
	dp.updatePaxosExecT(nat, rId)

	if dp.isFpLeaderLearner {
		if dp.GetCurFpShard() < 0 {
			logger.Fatalf("This is not a Fast Paxos shard leader, replicaId = %s", dp.getReplicaId())
		}

		// Handles a non-accept time for this Fast Paxos shard
		dp.processFpFastNonAcceptTime(nat, dp.getReplicaId(), rId)

		// Executes
		fpExecT := dp.getFpExecT(dp.GetCurFpShard())
		if dp.updateFpExecT(fpExecT) {
			dp.prepareShardExec(fpExecT.Shard)
		}
		dp.exec()
	} else {
		for fpShard, fpShardLeader := range dp.fpShardLeaderIdMap {
			dp.processFpFastNonAcceptTime(nat, fpShardLeader, rId)
			fpExecT := dp.getFpExecT(fpShard)
			if dp.updateFpExecT(fpExecT) {
				dp.prepareShardExec(fpExecT.Shard)
			}
		}
		dp.exec()
	}
}

//// Fast Paxos followers
// A Fast Paxos shard follower handles the accept request from the leader
func (dp *DynamicPaxos) handleFpAcceptReqMsg(msg *ReplicaMsg) {
	//loger.Infof("handling nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)
	dp.processFpAcceptReq(msg.Time, msg.Cmd)
	dp.fpReplicaHandleLeaderExecT(msg.NonAcceptT, msg.FpExecT, msg.ReplicaId)
}

// A Fast Paxos shard follower handles the accept request from the leader
func (dp *DynamicPaxos) handleFpCommitReqMsg(msg *ReplicaMsg) {
	//loger.Infof("handling nat = %d from %s", msg.NonAcceptT, msg.ReplicaId)
	isFast := msg.IsAccept
	dp.processFpCommitReq(msg.Time, msg.Cmd, isFast)
	dp.fpReplicaHandleLeaderExecT(msg.NonAcceptT, msg.FpExecT, msg.ReplicaId)
}

func (dp *DynamicPaxos) fpReplicaHandleLeaderExecT(
	nat int64, fpExecT *Timestamp, leaderId string,
) {
	// Updates the piggybacked Paxos shard non-accept time
	dp.updatePaxosExecT(nat, leaderId)
	// No need to prepare the Paxos shard execution since Paxos entries are
	// committed in an increasing order.

	if dp.isFpLeaderLearner {
		// Updates execution time for the Fast Paxos shard
		if fpExecT == nil {
			logger.Fatalf("Fast Paxos shard exec time is nil from leader = %s, nat = %d", leaderId, nat)
		}
		if dp.updateFpExecT(fpExecT) {
			dp.prepareShardExec(fpExecT.Shard)
		}
		dp.exec()
	} else {
		if fpExecT != nil {
			if dp.updateFpExecT(fpExecT) {
				dp.prepareShardExec(fpExecT.Shard)
			}
			dp.exec()
		}
	}
}

func (dp *DynamicPaxos) updatePaxosExecT(nat int64, leaderId string) {
	//loger.Infof("update nat = %d for id = %s", nat, leaderId)
	shard, ok := dp.getPaxosShard(leaderId)
	if !ok {
		return
	}
	pNat := &Timestamp{Time: nat, Shard: shard}
	dp.lm.UpdatePaxosMinExecT(pNat)
}

func (dp *DynamicPaxos) updateFpExecT(t *Timestamp) bool {
	return dp.lm.UpdateFpMinExecT(t)
}

func (dp *DynamicPaxos) prepareShardExec(shard int32) {
	dp.lm.PrepareShardExec(shard)
}

func (dp *DynamicPaxos) exec() {
	dp.lm.Exec()
}

//// Paxos Shard Follower
func (dp *DynamicPaxos) processPaxosAcceptReq(t *Timestamp, cmd *Command) {
	dp.lm.PaxosAcceptCmd(t, cmd)

	// Piggybacks Fast Paxos shard execution time if this is a Fast Paxos shard leader
	var fpExecT *Timestamp = nil
	if dp.GetCurFpShard() >= 0 {
		fpExecT = dp.getFpExecT(dp.GetCurFpShard())
	}
	// Send the result to the Paxos shard leader
	lId := dp.getPaxosShardLeaderId(t.Shard)
	l := dp.getReplicaAddr(lId)
	nat := dp.tm.GetCurrentTime()
	//loger.Infof("sending nat = %d", nat)
	dp.io.SendPaxosAcceptReply(l, t, cmd.Id, nat, dp.getReplicaId(), fpExecT)
}

func (dp *DynamicPaxos) processPaxosAcceptReply(t *Timestamp, cmdId string) {
	pCons := dp.getPaxosCons(t.Time)
	n := pCons.VoteAccept()
	if n == dp.Majority {
		// Async sends reply to clients
		retC := pCons.GetRetC()
		retC <- true

		// Piggybacks Fast Paxos shard execution time if this is a Fast Paxos shard leader
		var fpExecT *Timestamp = nil
		if dp.GetCurFpShard() >= 0 {
			fpExecT = dp.getFpExecT(dp.GetCurFpShard())
		}
		// Sends commit request to followers, which piggybacks Paxos non-accept time
		nat := dp.tm.GetCurrentTime()
		//loger.Infof("sending nat = %d", nat)
		dp.io.BcstPaxosCommitReq(dp.followerAddrList, t, cmdId, nat, dp.getReplicaId(), fpExecT)

		// Commits locally
		dp.processPaxosCommitReq(t, cmdId)
		dp.updatePaxosExecT(nat, dp.getReplicaId())
		dp.prepareShardExec(dp.GetCurPaxosShard())
		dp.exec()
	}

	if n == dp.ReplicaNum {
		dp.delPaxosCons(t.Time)
	}
}

//// Paxos Shard Follower
func (dp *DynamicPaxos) processPaxosCommitReq(t *Timestamp, cmdId string) {
	dp.lm.PaxosCommitCmd(t, cmdId)
}

//// Testing
func (dp *DynamicPaxos) Test() {
	dp.hbTimer.Stop()
	// Waits 1 second to process any buffered messages
	time.Sleep(1000 * 1000 * 1000)
	dp.testCh <- true
}

func (dp *DynamicPaxos) handleTest() {
	fmt.Printf("Paxos consenus instances size = %d\n", len(dp.pConsMap))
	dp.lm.Test()
	dp.fastPaxosTest()
}

func (dp *DynamicPaxos) fastPaxosTest() {
	fmt.Printf("Local fast paxos shard = %d\n", dp.fpShard)
	fmt.Printf("Fast Paxos consensus instances size = %d\n", len(dp.fpSlowConsMap))
	dp.cmdM.Test()
	for fpShard, cm := range dp.fpShardConsMap {
		fmt.Printf("Fast Paxos shard = %d\n", fpShard)
		cm.Test()
		dp.fpShardExecTMap[fpShard].Test()
	}
}

/////Helper Functions/////
//// General helper functions
// Resets the hearbeat timer, assuming that the timer has fired
func (dp *DynamicPaxos) resetHb(d time.Duration) {
	dp.hbTimer.Reset(d)
}

func (dp *DynamicPaxos) getReplicaId() string {
	return dp.replicaId
}

func (dp *DynamicPaxos) getReplicaNum() int {
	return dp.ReplicaNum
}

func (dp *DynamicPaxos) getFastQuorum() int {
	return dp.FastQuorum
}

func (dp *DynamicPaxos) getMajority() int {
	return dp.Majority
}

//// Paxos helper functions
func (dp *DynamicPaxos) getReplicaAddr(rId string) string {
	addr, ok := dp.replicaIdAddrMap[rId]
	if !ok {
		logger.Fatalf("Cannot find the addr for replica id = %s", rId)
	}
	return addr
}

func (dp *DynamicPaxos) getPaxosShardLeaderId(shard int32) string {
	l, ok := dp.pShardLeaderIdMap[shard]
	if !ok {
		logger.Fatalf("Cannot find the leader addr for paxos shard = %d", shard)
	}
	return l
}

func (dp *DynamicPaxos) getPaxosShard(leaderId string) (int32, bool) {
	shard, ok := dp.pLeaderIdShardMap[leaderId]
	if !ok {
		//logger.Fatalf("Cannot find the Paxos shard for leader = %s", leaderId)
		return -1, false
	}
	return shard, true
}

func (dp *DynamicPaxos) initPaxosCons(t int64, c chan bool) *PaxosCons {
	pCons := NewPaxosCons(c)
	dp.pConsMap[t] = pCons
	return pCons
}

func (dp *DynamicPaxos) getPaxosCons(t int64) *PaxosCons {
	p, ok := dp.pConsMap[t]
	if !ok {
		logger.Fatalf("Cannot find Paxos Consenus Instance for t = %d", t)

	}
	return p
}

func (dp *DynamicPaxos) delPaxosCons(t int64) {
	delete(dp.pConsMap, t)
}
