package dynamic

import (
	"fmt"
	"sync"

	"domino/common"
)

type FpCmdConsRet struct {
	T        int64
	CmdId    string
	IsAccept bool
	IsFast   bool
}

type FpCmdSlowConsRet struct {
	T        int64
	CmdIdMap map[string]bool
}

type FpCmdCons struct {
	set bool
	c   chan *FpCmdConsRet
}

func NewFpCmdCons() *FpCmdCons {
	return &FpCmdCons{set: false, c: make(chan *FpCmdConsRet, 1)}
}

type FpCmdConsRetManager struct {
	inputCh chan interface{}

	tQueue       *common.PriorityQueueInt64
	fpCmdConsMap map[int64]map[string]*FpCmdCons
	lock         sync.Mutex

	slowCmdConsMap map[int64]map[string]bool
}

func NewFpCmdConsRetManager(size int) *FpCmdConsRetManager {
	m := &FpCmdConsRetManager{
		inputCh:        make(chan interface{}, size),
		tQueue:         common.NewPriorityQueueInt64(),
		fpCmdConsMap:   make(map[int64]map[string]*FpCmdCons, size),
		slowCmdConsMap: make(map[int64]map[string]bool),
	}
	go m.start()
	return m
}

func (m *FpCmdConsRetManager) start() {
	for c := range m.inputCh {
		switch v := c.(type) {
		case int64:
			m.handleExecT(v)
		case *FpCmdConsRet:
			m.handleFpCmdConsRet(v)
		case *FpCmdSlowConsRet:
			m.handleFpCmdSlowConsRet(v)
		}
	}
}

func (m *FpCmdConsRetManager) handleExecT(execT int64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for t, ok := m.tQueue.Peek(); ok; t, ok = m.tQueue.Peek() {
		if t > execT {
			break
		}

		m.tQueue.Pop()

		if cMap, ok := m.fpCmdConsMap[t]; ok {
			slowMap, slow := m.slowCmdConsMap[t]
			for cmdId, cmdCons := range cMap {
				if slow {
					if _, found := slowMap[cmdId]; found {
						continue
					}
				}

				if !cmdCons.set {
					cmdCons.set = true
					cmdCons.c <- &FpCmdConsRet{IsAccept: false, IsFast: false}
				}
			}
		}
	}
}

func (m *FpCmdConsRetManager) handleFpCmdSlowConsRet(ret *FpCmdSlowConsRet) {
	if ret.CmdIdMap == nil {
		// Delete
		if _, ok := m.slowCmdConsMap[ret.T]; !ok {
			logger.Fatalf("The slow cons ret maps DOEST NOT exist!")
		}
		delete(m.slowCmdConsMap, ret.T)
	} else {
		// Add
		if _, ok := m.slowCmdConsMap[ret.T]; ok {
			logger.Fatalf("The slow cons ret maps exists!")
		}
		m.slowCmdConsMap[ret.T] = ret.CmdIdMap
	}
}

func (m *FpCmdConsRetManager) handleFpCmdConsRet(ret *FpCmdConsRet) {
	cmdCons := m.getFpCmdCons(ret.T, ret.CmdId, false)
	if cmdCons.set {
		logger.Fatalf("Cannot set cmd ret twice, t = %d cmdId = %s", ret.T, ret.CmdId)
	}
	cmdCons.set = true
	cmdCons.c <- ret
}

func (m *FpCmdConsRetManager) getFpCmdCons(
	t int64, cmdId string, isQueue bool,
) *FpCmdCons {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.fpCmdConsMap[t]; !ok {
		m.fpCmdConsMap[t] = make(map[string]*FpCmdCons)
		if isQueue {
			m.tQueue.Push(t)
		}
	}

	cMap := m.fpCmdConsMap[t]
	if _, found := cMap[cmdId]; !found {
		m.fpCmdConsMap[t][cmdId] = NewFpCmdCons()
	}

	return m.fpCmdConsMap[t][cmdId]
}

func (m *FpCmdConsRetManager) delFpCmdConsRetCh(t int64, cmdId string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if cMap, ok := m.fpCmdConsMap[t]; ok {
		if _, found := cMap[cmdId]; found {
			delete(cMap, cmdId)
		}
		if len(cMap) == 0 {
			delete(m.fpCmdConsMap, t)
		}
	}
}

func (m *FpCmdConsRetManager) Test() {
	fmt.Printf("Fp Command Consensus Resulst Manager: ")
	fmt.Printf("tQueue size = %d, fpCmdConsMap size = %d, slowCmdConsMap size = %d\n",
		m.tQueue.Len(), len(m.fpCmdConsMap), len(m.slowCmdConsMap))
}
