package dynamic

import (
	"fmt"
	"time"

	"domino/common"
)

type LogManager interface {
	//Updates the min execution time among all Paxos and Fast Paxos shards
	UpdatePaxosMinExecT(t *Timestamp)
	UpdateFpMinExecT(t *Timestamp) bool

	// Paxos
	PaxosAcceptCmd(t *Timestamp, cmd *Command)
	PaxosCommitCmd(t *Timestamp, cmdId string)

	// Fast Paxos
	// Fast-path accept a client's command
	FpFastAcceptCmd(t *Timestamp, cmd *Command) bool
	// Slow-path accept a command
	FpAcceptCmd(t *Timestamp, cmd *Command)
	// Fast-path / slow-path commit a command
	FpCommitCmd(t *Timestamp, cmd *Command, isFast bool)

	// Preapres a shard for execution
	PrepareShardExec(shard int32)

	// Execution
	Exec()
	GetExecCh() <-chan Entry

	// Testing
	Test()
}

type SimpleLogManager struct {
	n             int32                // total number of shards
	logMap        map[int32]Log        // shard Id --> log
	shardExecTMap map[int32]*Timestamp // shard Id --> max execution time on a shard

	execQueueMap map[int32]*common.Queue // shard Id --> continuous committed entry queue
	nextExecTMap map[int32]int64         // shard Id --> next to execute time on a shard
	execCh       chan Entry
}

func NewSimpleLogManager(
	pShardList []int32,
	fpShardList []int32,
	execChSize int,
) LogManager {
	lm := &SimpleLogManager{
		logMap:        make(map[int32]Log),
		shardExecTMap: make(map[int32]*Timestamp),
		execQueueMap:  make(map[int32]*common.Queue),
		nextExecTMap:  make(map[int32]int64),
		execCh:        make(chan Entry, execChSize),
	}

	for _, shard := range pShardList {
		lm.logMap[shard] = NewPaxosLog()
		lm.shardExecTMap[shard] = &Timestamp{Time: 0, Shard: shard}
		lm.execQueueMap[shard] = common.NewQueue()
		lm.nextExecTMap[shard] = 0
	}

	for _, shard := range fpShardList {
		lm.logMap[shard] = NewFpLog()
		lm.shardExecTMap[shard] = &Timestamp{Time: 0, Shard: shard}
		lm.execQueueMap[shard] = common.NewQueue()
		lm.nextExecTMap[shard] = 0
	}

	lm.n = int32(len(pShardList) + len(fpShardList))
	if lm.n == 0 {
		logger.Fatalf("No shards")
	}

	return lm
}

func (lm *SimpleLogManager) getShardLog(shard int32) Log {
	log, ok := lm.logMap[shard]
	if !ok {
		logger.Fatalf("Cannot find Paxos log for shard = %d", shard)
	}
	return log
}

// NOTE: for each Paxos shard, the execution time can only be updated to a future one.
func (lm *SimpleLogManager) UpdatePaxosMinExecT(t *Timestamp) {
	execT := lm.shardExecTMap[t.Shard]
	if CompareTime(t, execT) <= 0 {
		logger.Fatalf("Existing execution time = %v > time = %v at shard = %d", execT, t, t.Shard)
	}
	lm.updateMinExecT(t)
}

func (lm *SimpleLogManager) UpdateFpMinExecT(t *Timestamp) bool {
	execT := lm.shardExecTMap[t.Shard]
	//if dp.isFpLeaderLearner {
	//	if CompareTime(t, execT) < 0 {
	//		logger.Fatalf("Existing execution time = %v > time = %v at shard = %d", execT, t, t.Shard)
	//	}
	//	if CompareTime(t, execT) == 0 {
	//		return false
	//	}
	//}
	if CompareTime(t, execT) <= 0 {
		return false
	}

	lm.updateMinExecT(t)
	return true
}

func (lm *SimpleLogManager) updateMinExecT(t *Timestamp) {
	lm.shardExecTMap[t.Shard] = t
}

func (lm *SimpleLogManager) PaxosAcceptCmd(t *Timestamp, cmd *Command) {
	log := lm.getShardLog(t.Shard)
	e := &CmdEntry{
		Cmd:    cmd,
		T:      t,
		Status: ENTRY_STAT_LEADER_ACCEPTED,
		Duration: 0,
	}

	if en, ok := log.Get(t.Time); ok {
		logger.Fatalf("Cannot accept cmd = %s at t = %v [cmdId = %s]", cmd.Id, t, en.GetCmd().Id)
	}

	log.Put(t.Time, e)
}

func (lm *SimpleLogManager) PaxosCommitCmd(t *Timestamp, cmdId string) {
	log := lm.getShardLog(t.Shard)
	e, ok := log.Get(t.Time)

	if !ok {
		logger.Fatalf("Cannot commit cmd = %s at t = %v as it is not accepted yet", cmdId, t)
	}

	if e.GetCmd().Id != cmdId {
		logger.Fatalf("Cannot commit cmd = %s at t = %v [cmd = %s, t = %v, status = %d]",
			cmdId, t, e.GetCmd().Id, e.GetT(), e.GetStatus())
	}

	e.SetStartDuration(time.Now().UnixNano())

	e.SetStatus(ENTRY_STAT_COMMITTED)
}

func (lm *SimpleLogManager) FpFastAcceptCmd(t *Timestamp, cmd *Command) bool {
	log := lm.getShardLog(t.Shard)
	if _, ok := log.Get(t.Time); ok {
		return false
	}

	e := &CmdEntry{
		Cmd:    cmd,
		T:      t,
		Status: ENTRY_STAT_ACCEPTOR_ACCEPTED,
		Duration: 0,
	}
	log.Put(t.Time, e)

	return true
}

func (lm *SimpleLogManager) FpAcceptCmd(t *Timestamp, cmd *Command) {
	log := lm.getShardLog(t.Shard)

	if cmd == nil { // no-op
		log.Del(t.Time)
		return
	}

	if e, ok := log.Get(t.Time); ok {
		if e.GetCmd().Id == cmd.Id {
			e.SetStatus(ENTRY_STAT_LEADER_ACCEPTED)
			return
		}
	}

	e := &CmdEntry{
		Cmd:    cmd,
		T:      t,
		Status: ENTRY_STAT_LEADER_ACCEPTED,
		Duration: 0,
	}
	log.Put(t.Time, e)
}

func (lm *SimpleLogManager) FpCommitCmd(t *Timestamp, cmd *Command, isFast bool) {
	log := lm.getShardLog(t.Shard)
	e, ok := log.Get(t.Time)
	if isFast {
		// Fast-path commit may replace the cmd that is fast-accepted on a follower.
		if cmd == nil {
			if ok {
				log.Del(t.Time)
			}
			// NOTE: it is possible that a client reques arrives later on a replica
			// that have committed no-op by deleting the entry in the log. This is ok
			// since the replica will update its execution time to this time, and it
			// will not be able to fast accept any new command.
			return
		}

		if !ok || e.GetCmd().Id != cmd.Id {
			e = &CmdEntry{
				Cmd:    cmd,
				T:      t,
				Status: ENTRY_STAT_COMMITTED,
				Duration: 0,
			}
			log.Put(t.Time, e)
		} else {
			//if e.GetCmd().Id == cmd.Id
			e.SetStartDuration(time.Now().UnixNano())
			e.SetStatus(ENTRY_STAT_COMMITTED)
		}
	} else {
		// Slow-path commit comes after the slow-path accept due to FIFO.
		// In this case the id is the only parameter for normal commands.
		if cmd == nil {
			// The accept request for no-op should have arrived and replaced any existing command
			if ok {
				logger.Fatalf("The accept request for no-op at t = %v has no effect!", t)
			}
			return
		}

		if !ok {
			logger.Fatalf("Cannot commit cmd = %s at t = %v as it is not accepted yet", cmd.Id, t)
		}
		if e.GetCmd().Id != cmd.Id {
			logger.Fatalf("Cannot commit cmd = %s at t = %v [cmd = %s, t = %v, status = %d]",
				cmd.Id, t, e.GetCmd().Id, e.GetT(), e.GetStatus())
		}

		e.SetStartDuration(time.Now().UnixNano())
		e.SetStatus(ENTRY_STAT_COMMITTED)
	}
}

func (lm *SimpleLogManager) PrepareShardExec(shard int32) {
	lm.shardExec(shard)
}

func (lm *SimpleLogManager) Exec() {
	lm.logExec()
}

func (lm *SimpleLogManager) GetExecCh() <-chan Entry {
	return lm.execCh
}

// Puts the continuous committed cmds to a ready exec queue for the shard
func (lm *SimpleLogManager) shardExec(shard int32) {
	log, execT, q := lm.logMap[shard], lm.shardExecTMap[shard], lm.execQueueMap[shard]

	isFpLog := log.IsFpLog()
	e, ok := log.PeekMin()
	for ; ok; e, ok = log.PeekMin() {
		if CompareTime(e.GetT(), execT) > 0 {
			break
		}

		if isFpLog {
			if e.GetStatus() == ENTRY_STAT_ACCEPTOR_ACCEPTED {
				// NOTE: commit and slow-path accept req must be processed before the execT update.
				// For Fast Paxos, skips fast-path pending commands that are pending
				// and before execT.  For example, a replica vote to accept a command,
				// but the leader has moved the execT.  The accept vote will not be
				// processed, and the leader may not send any accept/commit for this
				// instance. But the leader moves the execT, which implifies that this
				// instance should be put a no-op.
				log.PopMin()
				continue
			}
		}

		if e.GetStatus() != ENTRY_STAT_COMMITTED {
			break
		}

		q.Push(e)
		log.PopMin()
	}

	if ok {
		lm.nextExecTMap[shard] = e.GetT().Time
		if isFpLog { // TODO better coding
			// Fast Paxos log may have a future accepted/commit entry
			if CompareTime(e.GetT(), execT) > 0 {
				lm.nextExecTMap[shard] = execT.Time
			}
		}
	} else {
		lm.nextExecTMap[shard] = execT.Time
	}
}

func (lm *SimpleLogManager) logExec() {
	for ok, e := lm.nextExecEntry(); ok; ok, e = lm.nextExecEntry() {
		lm.execCh <- e
	}
}

// Returns the next to execute entry if available
// Returns isOk, entry, shard#
func (lm *SimpleLogManager) nextExecEntry() (bool, Entry) {
	me, mt := lm.nextExecEntryByShard(0)
	var k int32 = 0
	var i int32 = 1
	for ; i < lm.n; i++ {
		e, t := lm.nextExecEntryByShard(i)

		if t < mt {
			me, mt, k = e, t, i
		}
	}

	if me != nil {
		q := lm.execQueueMap[k]
		e, _ := q.Pop()
		return true, e.(Entry)
	}
	return false, nil
}

// Returns next-to-execute entry in the shard and its timestamp.
// Otherwise, returns nil and the next-to-execute timestamp.
func (lm *SimpleLogManager) nextExecEntryByShard(shard int32) (Entry, int64) {
	q := lm.execQueueMap[shard]
	if e, ok := q.Peek(); ok {
		return e.(Entry), e.(Entry).GetT().Time
	}
	return nil, lm.nextExecTMap[shard]
}

func (lm *SimpleLogManager) Test() {
	for shard, l := range lm.logMap {
		fmt.Printf("Log shard = %d, cmdQSize = %d, shardExecT = %v\n",
			shard, lm.execQueueMap[shard].Size(), lm.shardExecTMap[shard])
		l.Test()
		fmt.Println()
	}
}
