package dynamic

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"

	"domino/common"
)

type Log interface {
	PeekMin() (Entry, bool)
	PopMin() (Entry, bool)

	Put(t int64, e Entry)
	Get(t int64) (Entry, bool)

	Del(t int64)

	IsFpLog() bool

	// Testing
	Test()
}

type PaxosLog struct {
	eMap map[int64]Entry // time --> entry
	tQ   *common.Queue   // time are sorted in an increasing order
}

func NewPaxosLog() Log {
	return &PaxosLog{
		eMap: make(map[int64]Entry, DEFAULT_PAXOS_LOG_INIT_SIZE),
		tQ:   common.NewQueue(),
	}
}

func (l *PaxosLog) IsFpLog() bool {
	return false
}

// Assuming that all entries are added in the increasing order of the time
func (l *PaxosLog) PeekMin() (Entry, bool) {
	if t, ok := l.tQ.Peek(); ok {
		return l.eMap[t.(int64)], true
	}
	return nil, false
}

func (l *PaxosLog) PopMin() (Entry, bool) {
	if t, ok := l.tQ.Pop(); ok {
		e := l.eMap[t.(int64)]
		delete(l.eMap, t.(int64))
		return e, true
	}
	return nil, false
}

// Puts the cmd to the Paxos shard log, where cmds are ordered by the timestamps.
// NOTE: this function should be called in the order of the timestamp assignment.
func (l *PaxosLog) Put(t int64, e Entry) {
	l.eMap[t] = e
	l.tQ.Push(t)
}

func (l *PaxosLog) Get(t int64) (Entry, bool) {
	e, ok := l.eMap[t]
	return e, ok
}

func (l *PaxosLog) Del(t int64) {
	logger.Fatalf("Should not execute here on a Paxos shard log. t = %d", t)
}

func (l *PaxosLog) Test() {
	fmt.Printf("Paxos log entry map size = %d, queue size = %d\n", len(l.eMap), l.tQ.Size())
	for t, e := range l.eMap {
		fmt.Printf("%v %s %d", t, e.GetCmd().Id, e.GetStatus())
	}
}

type FpLog struct {
	tMap *treemap.Map // backed by red-black tree
}

// Fast Paxos log
func NewFpLog() Log {
	l := &FpLog{
		tMap: treemap.NewWith(utils.Int64Comparator),
	}
	return l
}

func (l *FpLog) IsFpLog() bool {
	return true
}

func (l *FpLog) PeekMin() (Entry, bool) {
	_, e := l.tMap.Min()
	if e == nil {
		return nil, false
	}
	return e.(Entry), true
}

func (l *FpLog) PopMin() (Entry, bool) {
	i, e := l.tMap.Min()
	if e == nil {
		return nil, false
	}
	l.tMap.Remove(i)
	return e.(Entry), true
}

func (l *FpLog) Put(t int64, e Entry) {
	l.tMap.Put(t, e)
}

func (l *FpLog) Get(t int64) (Entry, bool) {
	if e, ok := l.tMap.Get(t); ok {
		return e.(Entry), ok
	}
	return nil, false
}

func (l *FpLog) Del(t int64) {
	l.tMap.Remove(t)
}

func (l *FpLog) Test() {
	fmt.Printf("Fast Paxos log entry map size = %d\n", l.tMap.Size())
	itr := l.tMap.Iterator()
	for itr.Next() {
		e := itr.Value().(Entry)
		fmt.Println(e.GetT(), e.GetStatus())
	}
}
