package dynamic

import (
	"sync"
)

type ExecRet struct {
	C chan string
}

func NewExecRet() *ExecRet {
	return &ExecRet{
		C: make(chan string, 1),
	}
}

type ExecManager struct {
	execRetMap     map[string]*ExecRet
	execRetMapLock sync.Mutex
}

func NewExecmanager() *ExecManager {
	return &ExecManager{
		execRetMap: make(map[string]*ExecRet, 200000),
	}
}

func (em *ExecManager) GetExecRet(cmdId string) *ExecRet {
	em.execRetMapLock.Lock()
	defer em.execRetMapLock.Unlock()

	if _, ok := em.execRetMap[cmdId]; !ok {
		em.execRetMap[cmdId] = NewExecRet()
	}

	return em.execRetMap[cmdId]
}

func (em *ExecManager) DelExecRet(cmdId string) {
	em.execRetMapLock.Lock()
	defer em.execRetMapLock.Unlock()

	delete(em.execRetMap, cmdId)
}

// Waits for execution result to be ready
// Blocks until execution result is available
func (em *ExecManager) WaitExecRet(cmdId string) string {
	execRet := em.GetExecRet(cmdId)
	defer em.DelExecRet(cmdId)

	return <-execRet.C
}

func (em *ExecManager) Size() int {
	em.execRetMapLock.Lock()
	defer em.execRetMapLock.Unlock()

	return len(em.execRetMap)
}
