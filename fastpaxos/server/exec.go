package main

import (
	"domino/common"
	fp "domino/fastpaxos/fastpaxos"
	"domino/fastpaxos/rpc"
)

type ExecRet struct {
	OpId string
	Ret  string
}

func (s *Server) execOpWithRet(op *rpc.Operation, retHandle *common.Future) {
	if op == nil {
		logger.Fatalf("Operation applied to kv-store should not be nil")
	}
	logger.Infof("write")

	switch op.Type {
	case fp.OP_WRITE:
		s.kvStore.Write(op.Key, op.Val)
		retHandle.SetValue(&ExecRet{OpId: op.Id, Ret: "true"})
	case fp.OP_READ:
		val, _ := s.kvStore.Read(op.Key)
		retHandle.SetValue(&ExecRet{OpId: op.Id, Ret: val})
	case fp.OP_STOP:
		//logger.Infof("Execution stopped")
		retHandle.SetValue(&ExecRet{OpId: op.Id, Ret: "stopped"})
	}
}

func (s *Server) GetExecRetHandle(t int64) *common.Future {
	s.execRetTableLock.Lock()
	defer s.execRetTableLock.Unlock()

	if _, exists := s.execRetTable[t]; !exists {
		s.execRetTable[t] = common.NewFuture()
	}

	return s.execRetTable[t]
}

func (s *Server) DelExecHandle(t int64) {
	s.execRetTableLock.Lock()
	defer s.execRetTableLock.Unlock()

	delete(s.execRetTable, t)
}
