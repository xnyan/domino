package com

import (
	"sync"

	"github.com/op/go-logging"
	"google.golang.org/grpc"

	"domino/fastpaxos/rpc"
)

var logger = logging.MustGetLogger("comm")

// Server communication table
type CommManager struct {
	connTable     map[string]*grpc.ClientConn // Network addr (ip:port) --> RPC connection
	connTableLock sync.RWMutex

	rpcStubTable     map[string]rpc.FastPaxosRpcClient // Network addr (ip:port)--> RPC stub
	rpcStubTableLock sync.RWMutex
}

func NewCommManager() *CommManager {
	return &CommManager{
		connTable:    make(map[string]*grpc.ClientConn),
		rpcStubTable: make(map[string]rpc.FastPaxosRpcClient),
	}
}

func (cm *CommManager) BuildConnection(addr string) *grpc.ClientConn {
	cm.connTableLock.Lock()
	defer cm.connTableLock.Unlock()

	if conn, exists := cm.connTable[addr]; exists {
		return conn
	}

	logger.Debugf("Connecting to server %s", addr)
	var err error
	cm.connTable[addr], err = grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		logger.Fatalf("Cannot build connection to server %s, error: %v", addr, err)
	}

	return cm.connTable[addr]
}

func (cm *CommManager) GetConnection(addr string) (*grpc.ClientConn, bool) {
	cm.connTableLock.RLock()
	defer cm.connTableLock.RUnlock()

	conn, exists := cm.connTable[addr]
	return conn, exists
}

// Builds an RPC stub if not exists. Otherwise, returns the existing one
func (cm *CommManager) BuildRpcStub(addr string) rpc.FastPaxosRpcClient {
	cm.rpcStubTableLock.Lock()
	defer cm.rpcStubTableLock.Unlock()

	if rpcStub, exists := cm.rpcStubTable[addr]; exists {
		return rpcStub
	}

	conn := cm.BuildConnection(addr)

	cm.rpcStubTable[addr] = rpc.NewFastPaxosRpcClient(conn)

	return cm.rpcStubTable[addr]
}

// Returns the RPC stub on file
func (cm *CommManager) GetRpcStub(addr string) (rpc.FastPaxosRpcClient, bool) {
	cm.rpcStubTableLock.RLock()
	defer cm.rpcStubTableLock.RUnlock()

	rpcStub, exists := cm.rpcStubTable[addr]
	return rpcStub, exists
}

///////////////////
// Streaming pattern
/*
type StreamManager struct {
	*CommManager
	commitChTable [string]chan *rpc.CommitRequest // server addr -->

	// For leaders to call RPCs on followers
	AcceptTable [string]rpc.Consensus_AcceptClient
	CommitTable [string]rpc.Consensus_CommitClient
}

func (sm *StreamManager) SendCommitRequest(addr string, req *rpc.CommitRequest) {
	sm.commitChTableLock.RLock()
	commitCh, exists := sm.commitChTable[addr]
	sm.commitChTableLock.RUnlock()

	if !exists {
		commitCh = sm.StartCommitStream(addr)
	}

	commitCh <- req
}

func (sm *StreaManager) StartCommitStream(addr string) chan *rpc.CommitRequest {
	sm.commitChTableLock.Lock()
	defer sm.commitChTableLock.Unlock()

	if commitCh, exists := sm.commitChTable[addr]; exists {
		return commitCh
	}

	sm.commitChTable[addr] = make(chan *rpc.CommmitRequest, ServerRpcServerBufferSize)
	commitCh = sm.commitChTable[addr]

	// Starts the stream RPC
	go func() {
		rpcStub := sm.BuildRpcStub(addr)
		stream, err := rpcStub.Commit(context.Background())
		if err != nil {
			logger.Fatalf("Fails to start Commit stream rpc to addr = %s", addr)
		}
		for req := range commitCh {
			if e := stream.Send(req); e != nil {
				logger.Fatalf("Fails to send commit request operation id = %s, idx = %d", req.Op.Id, req.Idx)
			}
		}
		stream.CloseSend()
	}()

	return commitCh
}
*/
