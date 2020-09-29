package fastpaxos

import (
	"sync"

	"domino/common"
)

const (
	INVALID_IDX = "-1"
)

type Vote struct {
	OpId string
	Idx  string
}

type FastPathManager struct {
	ReplicaNum int // Total number of replicas
	Majority   int // Majority number of replicas
	FastQuorum int // Supermajority number of replicas

	voteCh chan *Vote

	// Uses a future to store the fast-path result (i.e., idx).
	opFastRetTable     map[string]*common.Future // Uses INVALID_IDX if the fast path fails
	opFastRetTableLock sync.RWMutex
	opSlowRetTable     map[string]*common.Future
	opSlowRetTableLock sync.RWMutex

	opTable  map[string]*OpInfo  // opId --> operation voting info
	idxTable map[string]*IdxInfo // log idx --> idx voting info

	opQueue  *common.Queue // operations that do not have fast or slow path yet but all votes
	idxQueue *common.Queue // idxes that do not have fast or slow path but all votes

	isRun   bool
	runLock sync.RWMutex
}

func NewFastPathManager(
	replicaNum, majority, fastQuorum int,
	voteChBufferSize int,
) *FastPathManager {
	return &FastPathManager{
		ReplicaNum: replicaNum,
		Majority:   majority,
		FastQuorum: fastQuorum,

		voteCh: make(chan *Vote, voteChBufferSize),

		opFastRetTable: make(map[string]*common.Future),
		opSlowRetTable: make(map[string]*common.Future),

		opTable:  make(map[string]*OpInfo),
		idxTable: make(map[string]*IdxInfo),

		opQueue:  common.NewQueue(),
		idxQueue: common.NewQueue(),

		isRun: false,
	}
}

func (f *FastPathManager) GetOpFastRet(opId string) *common.Future {
	f.opFastRetTableLock.Lock()
	defer f.opFastRetTableLock.Unlock()

	if _, exists := f.opFastRetTable[opId]; !exists {
		f.opFastRetTable[opId] = common.NewFuture()
	}

	ret := f.opFastRetTable[opId]
	return ret
}

func (f *FastPathManager) RemoveOpFastRet(opId string) {
	f.opFastRetTableLock.Lock()
	defer f.opFastRetTableLock.Unlock()

	delete(f.opFastRetTable, opId)
}

func (f *FastPathManager) GetOpSlowRet(opId string) *common.Future {
	f.opSlowRetTableLock.Lock()
	defer f.opSlowRetTableLock.Unlock()

	if _, exists := f.opSlowRetTable[opId]; !exists {
		f.opSlowRetTable[opId] = common.NewFuture()
	}

	ret := f.opSlowRetTable[opId]
	return ret
}

func (f *FastPathManager) RemoveOpSlowRet(opId string) {
	f.opSlowRetTableLock.Lock()
	defer f.opSlowRetTableLock.Unlock()

	delete(f.opSlowRetTable, opId)
}

func (f *FastPathManager) GetOpInfo(opId string) *OpInfo {
	if _, exists := f.opTable[opId]; !exists {
		f.opTable[opId] = NewOpInfo()
	}
	opInfo := f.opTable[opId]

	return opInfo
}

func (f *FastPathManager) RemoveOpInfo(opId string) {
	delete(f.opTable, opId)
}

func (f *FastPathManager) GetIdxInfo(idx string) *IdxInfo {
	if _, exists := f.idxTable[idx]; !exists {
		f.idxTable[idx] = NewIdxInfo()
	}
	idxInfo := f.idxTable[idx]
	return idxInfo
}

func (f *FastPathManager) RemoveIdxInfo(idx string) {
	delete(f.idxTable, idx)
}

// This function can only be called when all of listeners have acquired the
// hendle to the fast-path and slow-path ret
func (f *FastPathManager) CleanRetHandle(opId string) {
	f.RemoveOpFastRet(opId)
	f.RemoveOpSlowRet(opId)
}

// The fast-path coordinator / leader should call this function to acquire handles to wait
// for fast-path or slow-path results.
// Returns handles (fast-path furute, slow-path future), which stores the idx to commit the operation
func (f *FastPathManager) LeaderVote(vote *Vote) (*common.Future, *common.Future) {
	// Gets the handle first in case that it is deleted when voting id done before this
	opFastRet := f.GetOpFastRet(vote.OpId)
	opSlowRet := f.GetOpSlowRet(vote.OpId)

	// Avoids being blocked on the vote channel
	go func() { f.Vote(vote) }()

	return opFastRet, opSlowRet
}

// Followers should just use this function to vote
// Accepts a vote via the fast path
func (f *FastPathManager) Vote(vote *Vote) {
	f.voteCh <- vote
}

// Starts the fast-path manager
func (f *FastPathManager) Run() {
	f.runLock.Lock()
	defer f.runLock.Unlock()
	if f.isRun {
		return
	}
	f.isRun = true

	go f.startVoteCh()
}

// Single thread for processing votes
func (f *FastPathManager) startVoteCh() {
	for vote := range f.voteCh {

		logger.Debugf("Fast-path processes vote opId = (%s) at idx = (%s)", vote.OpId, vote.Idx)

		// Idx
		idxInfo := f.GetIdxInfo(vote.Idx)
		opVoteN := idxInfo.Vote(vote.OpId) // fast-path vote (an op on the idx)

		// Op
		opInfo := f.GetOpInfo(vote.OpId)
		opInfo.Vote(vote.Idx)

		// Only handles the fast-path success once
		if opVoteN == f.FastQuorum {

			logger.Debugf("Fast-path succeeds for opId = (%s) at idx = (%s)", vote.OpId, vote.Idx)

			// Records fast-path results
			idxInfo.SetFastOpId(vote.OpId)
			opInfo.SetFastIdx(vote.Idx)

			// Does not commit fast-path result for an idx before receiving all of
			// the votes for the idx.  Otherwise, the idx may be committed earlier on
			// a replica before the replica has voted any operation for the idx.  In
			// this case, the replica will miss voting an idx, where it will use the
			// next available idx to vote for the same operation. The replica will
			// end up voting more idxes than other replicas.  This idx-offest for
			// voting might also lead to a cascading failure of the fast path.

			// For now, commits the fast-path result after receiving all of votes on an idx
			//TODO Uses a more elegant solution
		}

		if idxInfo.GetTotalVoteNum() == f.ReplicaNum {
			// All replicas have voted for the idx
			if idxInfo.isFast() {
				// The fast path succeeds at idx but may be an operation different from the current vote
				opId := idxInfo.GetFastOpId()
				opFastRet := f.GetOpFastRet(opId)

				logger.Debugf("Fast-path allows committing opId = (%s) at idx (%s)", opId, vote.Idx)

				// Notifies any listeners that are waiting for the fast-path result
				opFastRet.SetValue(vote.Idx) // Starts to commit the operation at idx
			} else {
				// The fast path fails at idx, uses a slow path
				if opId, exists := idxInfo.getOpId(f.Majority); exists {
					// An operation has been accepted at idx by at least a majority of replicas
					// The slow path will commit this operation at idx.
					opSlowRet := f.GetOpSlowRet(opId)

					logger.Debugf("Slow-path allows commmitting opId = (%s) at idx = (%s)", opId, vote.Idx)

					// Notifies any listeners that are waiting for the slow-path result
					// NOTE: a listener may not have the slow-path result handle yet.
					// Therefore, cannot simply remove the slow-path result.
					opSlowRet.SetValue(vote.Idx)
				} else {
					// By randomly picking an operation, the operation may be duplicated.
					// Also, another operation may starve. Has to track all of the
					// potential starving operations, and put them to future log entries
					// that will not have a fast path.
					//
					// Instead of randomly picking an operation, wait to pick one
					// operation that is determined not to have a successful fast path.

					logger.Debugf("Slow-path puts idx = (%s) to the waiting queue", vote.Idx)

					// Puts this idx into the pending idx queue
					f.idxQueue.Push(vote.Idx)
				}
			}

			// Removes the idx info because there is no more votes expected at this idx
			f.RemoveIdxInfo(vote.Idx)
		}

		if opInfo.GetTotalVoteNum() == f.ReplicaNum {
			// The operation has been accepted by all of the replicas
			if opInfo.IsFast() {
				// The fast path has succeeded for this operation at an idx (not the
				// one in the current vote)
				// Does nothing since the operation is committed only after receiving all votes on idx

				//logger.Debugf("Fast-path nothing for fast-path opId = (%s)", vote.OpId)

			} else {
				// The fast path fails for the operation
				opFastRet := f.GetOpFastRet(vote.OpId)

				logger.Debugf("Slow-path sets the failure of the fast-path for opId = (%s)", vote.OpId)

				// Notifies any listeners that the fast path fails
				opFastRet.SetValue(INVALID_IDX)

				// Does not start a slow path based on the number of votes on an
				// operation. Otherwise, the slow-path would commit an operation
				// before a replica is able to vote for that idx.  In this case, a
				// replica would miss a vote on the idx and end up using more idxes
				// for voting.
				if _, exists := opInfo.getIdx(f.Majority); exists {
					// The operation has been accepted at an idx by at least a majority of replicas.
					// The slow-path for the idx will commit the operation.

					//logger.Debugf("Slow-path nothing for majority fast-path-accepted opId = (%s)", vote.OpId)

				} else {
					logger.Debugf("Slow-path puts opId = (%s) to the waiting queue", vote.OpId)

					// Puts this operation into the pending operation queue
					f.opQueue.Push(vote.OpId)
				}
			}

			// No more expected votes for this operation
			f.RemoveOpInfo(vote.OpId)

		}

		// Goes over the pending idx queue and the pending operation queue to start slow path
		f.RunSlowPath()
	}
}

// Starts a slow path for any pending idx and pending operation
// Put one pending operation at one pending idx, where the operation and the
// idx must not have fast path or slow path yet for sure.
func (f *FastPathManager) RunSlowPath() {
	for f.opQueue.Size() > 0 && f.idxQueue.Size() > 0 {
		opId, _ := f.opQueue.Pop()
		idx, _ := f.idxQueue.Pop()

		logger.Debugf("Slow-path matches opId = (%s) at idx = (%s)", opId.(string), idx.(string))

		opSlowRet := f.GetOpSlowRet(opId.(string))

		logger.Debugf("Slow-path allows committing opId = (%s) at idx = (%s)", opId.(string), idx.(string))

		opSlowRet.SetValue(idx)

		//f.RemoveOpSlowRet(opId.(string)) // Any listener must have a handle to opSlowRet at this point
		//f.RemoveOpFastRet(vote.OpId) // Any listener must have a handle to opFastRet at this point
	}
}
