package fastpaxos

import (
	"sync"
	"time"

	"github.com/op/go-logging"

	"domino/common"
	"domino/fastpaxos/rpc"
)

var logger = logging.MustGetLogger("fastpaxos")

const (
	OP_STOP  = "s"
	OP_WRITE = "w"
	OP_READ  = "r"
)

type ClientProposal struct {
	Op        *rpc.Operation
	Idx       *common.Future
	Delay     time.Duration
	Timestamp int64
	ClientId  string
}

type LeaderProposal struct {
	Idx  *LogIdx
	Op   *rpc.Operation
	Done *common.Future
}

type LeaderCommit struct {
	Idx  *LogIdx
	Op   *rpc.Operation
	Done *common.Future
}

type FastPaxos struct {
	// Fast path manager
	fpManager *FastPathManager

	// The channel for a replica to execute a command including:
	// (1) accpets a client's proposal via the fast path
	// (2) accpets a leader's proposal via the slow path
	// (3) commits a leader's proposal in the log
	cmdCh chan interface{}

	// The channel fo a replica to apply (execute) a committed log entry on the state machine
	execCh chan *rpc.Operation

	// Log manager
	lm LogManager

	isRun   bool
	runLock sync.RWMutex

	// Scheduling
	scheduler Scheduler
}

func NewFastPaxos(
	replicaNum, majority, fastQuorum int,
	voteChBufferSize int,
	cmdChBufferSize int,
	execChBufferSize int,
	logManager LogManager,
	scheduleType string,
	processWindow time.Duration, // for timestamp scheduler
) *FastPaxos {
	fp := &FastPaxos{
		fpManager: NewFastPathManager(replicaNum, majority, fastQuorum, voteChBufferSize),

		cmdCh:  make(chan interface{}, cmdChBufferSize),
		execCh: make(chan *rpc.Operation, execChBufferSize),

		lm: logManager,

		isRun: false,
	}

	// Scheduler
	absScheduler := NewAbstractScheduler(fp.GetCmdCh())
	switch scheduleType {
	case common.NoScheduler:
		fp.scheduler = &NoScheduler{
			absScheduler,
		}
	case common.DelayScheduler:
		fp.scheduler = &DelayScheduler{
			absScheduler,
		}
	case common.TimestampScheduler:
		fp.scheduler = NewTimestampScheduler(
			absScheduler,
		)
	case common.DeterministicScheduler:
		logger.Fatalf("Scheduler %s not implemented yet", scheduleType)
	default:
		logger.Fatalf("Unknow scheduler type %s", scheduleType)
	}

	return fp
}

func (fp *FastPaxos) Schedule(p *ClientProposal) {
	fp.scheduler.Schedule(p)
}

func (fp *FastPaxos) GetCmdCh() chan<- interface{} {
	return fp.cmdCh
}

func (fp *FastPaxos) RunCmd(cmd interface{}) {
	fp.cmdCh <- cmd
}

func (fp *FastPaxos) GetExecCh() <-chan *rpc.Operation {
	return fp.execCh
}

// Starts fast paxos
// The leader should start the fast-path manager
// Followers that are only acceptors do not need to start the fast-path manager
func (fp *FastPaxos) Run(isRunFastPathManager bool) {
	fp.runLock.Lock()
	defer fp.runLock.Unlock()
	if fp.isRun {
		return
	}
	fp.isRun = true

	// Starts log manager
	fp.lm.Run(fp.execCh)

	// Starts command channel
	go fp.startCmdCh()

	// Starts scheduler
	go fp.scheduler.Run()

	// Learners starts the fast-path manager
	if isRunFastPathManager {
		fp.fpManager.Run()
	}
}

// Single thread
func (fp *FastPaxos) startCmdCh() {
	for proposal := range fp.cmdCh {
		switch proposal.(type) {
		case *ClientProposal:
			fp.acceptClientProposal(proposal.(*ClientProposal))
		case *LeaderProposal:
			fp.acceptLeaderProposal(proposal.(*LeaderProposal))
		case *LeaderCommit:
			fp.commitLeaderProposal(proposal.(*LeaderCommit))
		default:
			logger.Fatalf("Unknown command type: %v", proposal)
		}
	}
}

func (fp *FastPaxos) acceptClientProposal(proposal *ClientProposal) {
	//TODO Ignores an operation that is committed via the slow path before the
	//client's proposal arrives.
	entry := &Entry{op: proposal.Op, status: ENTRY_FAST_ACCEPTED}
	idx, err := fp.lm.FastPathAccept(entry)

	if err != nil {
		logger.Errorf("Fast-path fails to accept opId = (%s), error: %v", err)
	}

	proposal.Idx.SetValue(idx)

	logger.Debugf("Fast-path accepts opId = (%s) at idx (%s)", entry.op.Id, idx)
}

func (fp *FastPaxos) acceptLeaderProposal(proposal *LeaderProposal) {
	defer proposal.Done.SetValue(true)

	logger.Debugf("Slow-path tries to accept opId = (%s) at idx = (%s)",
		proposal.Op.Id, proposal.Idx)

	entry := &Entry{op :proposal.Op, status: ENTRY_SLOW_ACCEPTED}
	err := fp.lm.SlowPathAccept(proposal.Idx, entry)

	if err != nil {
		logger.Errorf("Slow-path fails to accept opId = (%s) at idx = (%s), error: %v", err)
	}
}

func (fp *FastPaxos) commitLeaderProposal(commit *LeaderCommit) {
	defer commit.Done.SetValue(true)

	logger.Debugf("Leader commits opId = (%s) at idx = (%s)", commit.Op.Id, commit.Idx)

	entry := &Entry{op: commit.Op}
	err := fp.lm.Commit(commit.Idx, entry)

	if err != nil {
		logger.Errorf("Leader fails to commit opId = (%s) at idx = (%s)", commit.Op.Id, commit.Idx)
	}
}

/////////////////////////////////////////////////
// Wrapper functions for FastPathManager

func (fp *FastPaxos) LeaderVote(vote *Vote) (*common.Future, *common.Future) {
	return fp.fpManager.LeaderVote(vote)
}

func (fp *FastPaxos) Vote(vote *Vote) {
	fp.fpManager.Vote(vote)
}

func (fp *FastPaxos) CleanRetHandle(opId string) {
	fp.fpManager.CleanRetHandle(opId)
}
