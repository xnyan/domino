package dynamic

import (
	"fmt"

	"domino/common"
)

// Fast Paxos shard consensus instance manager
type FpConsManager struct {
	consMap   map[int64]*FpCons
	consQueue *common.PriorityQueueInt64

	slowConsMap map[int64]*FpCons
}

func NewFpConsManager(initSize int) *FpConsManager {
	return &FpConsManager{
		consMap:     make(map[int64]*FpCons, initSize),
		consQueue:   common.NewPriorityQueueInt64(),
		slowConsMap: make(map[int64]*FpCons),
	}
}

func (m *FpConsManager) GetFpCons(t int64) *FpCons {
	if _, ok := m.consMap[t]; !ok {
		m.consMap[t] = NewFpCons(t)
		m.consQueue.Push(t)
	}
	return m.consMap[t]
}

func (m *FpConsManager) PeekCons() (*FpCons, bool) {
	t, ok := m.consQueue.Peek()
	if !ok {
		return nil, false
	}

	if _, ok := m.consMap[t]; !ok {
		logger.Fatalf("Cannot find fast paxos consensus instance for t = %d", t)
	}
	return m.consMap[t], true
}

func (m *FpConsManager) PopCons() *FpCons {
	_, ok := m.consQueue.Peek()
	if !ok {
		return nil
	}
	t := m.consQueue.Pop()
	if _, ok := m.consMap[t]; !ok {
		logger.Fatalf("Cannot find fast paxos consensus instance for t = %d", t)
	}
	cons := m.consMap[t]
	delete(m.consMap, t)
	return cons
}

func (m *FpConsManager) AddSlowFpCons(cons *FpCons) {
	if _, ok := m.slowConsMap[cons.GetT()]; ok {
		logger.Fatalf("Cannot add slow fp cons twice t = %d", cons.GetT())
	}
	m.slowConsMap[cons.GetT()] = cons
}

func (m *FpConsManager) GetSlowFpCons(t int64) *FpCons {
	return m.slowConsMap[t]
}

func (m *FpConsManager) DelSlowFpCons(t int64) {
	delete(m.slowConsMap, t)
}

// Fast Paxos shard execution time manager
type FpExecTimeManager struct {
	// Time before which all commands to commit (may not be committed yet) are
	// all known, and continuous committed commands can get executed.
	execT int64

	rIdList []string
	tList   []int64 // a list of fast non-accept time in increasing order
	min     int     // idx for the min non-accept time in all replicas
	fqMin   int     // idx for the min time in a fast quorum with highest non-accept time

	rIdxMap map[string]int // rId --> idx in nAtList
}

func NewFpExecTimeManager(rN int, fastQuorum int, rIdList []string) *FpExecTimeManager {
	tm := &FpExecTimeManager{
		rIdList: make([]string, rN),
		tList:   make([]int64, rN),
		min:     0,
		fqMin:   rN - fastQuorum,

		rIdxMap: make(map[string]int),
	}

	for i, rId := range rIdList {
		tm.tList[i] = 0
		tm.rIdList[i] = rId
		tm.rIdxMap[rId] = i
	}

	return tm
}

func (tm *FpExecTimeManager) UpdateExecT(t int64) {
	tm.execT = t
}

func (tm *FpExecTimeManager) GetExecT() int64 {
	return tm.execT
}

func (tm *FpExecTimeManager) GetFastNonAcceptT(rId string) int64 {
	i, ok := tm.rIdxMap[rId]
	if !ok {
		logger.Fatalf("No fast non-accept time for replica = %s", rId)
	}
	return tm.tList[i]
}

func (tm *FpExecTimeManager) GetMinNonAcceptT() (int64, string) {
	return tm.tList[tm.min], tm.rIdList[tm.min]
}

func (tm *FpExecTimeManager) GetFastQuorumNonAcceptT() (int64, string) {
	return tm.tList[tm.fqMin], tm.rIdList[tm.fqMin]
}

func (tm *FpExecTimeManager) UpdateNonAcceptTime(rId string, nat int64) {
	i, ok := tm.rIdxMap[rId]
	if !ok {
		logger.Fatalf("No fast non-accept time for replica = %s", rId)
	}

	t := tm.tList[i]
	if nat <= t {
		logger.Fatalf("Cannot update fast non-accept time from %d to %d for replica = %s",
			t, nat, rId)
	}

	tm.tList[i] = nat
	for j := i + 1; j < len(tm.tList); j++ {
		if tm.tList[j] >= tm.tList[i] {
			break
		}
		// Swap
		tm.tList[i], tm.tList[j] = tm.tList[j], tm.tList[i]
		tm.rIdList[i], tm.rIdList[j] = tm.rIdList[j], tm.rIdList[i]
		tm.rIdxMap[tm.rIdList[i]] = i
		tm.rIdxMap[tm.rIdList[j]] = j
		i = j
	}
}

// Testing
func (m *FpConsManager) Test() {
	fmt.Printf("Fast Paxos consensus manager: ")
	fmt.Printf("consMap size = %d, consQueue size = %d slowConsMap size = %d\n",
		len(m.consMap), m.consQueue.Len(), len(m.slowConsMap))
}

func (tm *FpExecTimeManager) Test() {
	fmt.Printf("Fast Paxos execT manager: execT = %d\n", tm.execT)
}
