package dynamic

import (
	"sync"

	"google.golang.org/grpc"
)

// Server communication table
type CommManager struct {
	connMap     map[string]*grpc.ClientConn // Network addr (ip:port) --> RPC connection
	connMapLock sync.RWMutex

	rpcStubMap     map[string]DynamicPaxosClient // Network addr (ip:port)--> RPC stub
	rpcStubMapLock sync.RWMutex
}

func NewCommManager() *CommManager {
	return &CommManager{
		connMap:    make(map[string]*grpc.ClientConn),
		rpcStubMap: make(map[string]DynamicPaxosClient),
	}
}

func (cm *CommManager) BuildConnection(addr string) *grpc.ClientConn {
	cm.connMapLock.Lock()
	defer cm.connMapLock.Unlock()

	if conn, exists := cm.connMap[addr]; exists {
		return conn
	}

	logger.Debugf("Connecting to server %s", addr)
	var err error
	cm.connMap[addr], err = grpc.Dial(addr,
		grpc.WithInsecure(), grpc.WithReadBufferSize(0), grpc.WithWriteBufferSize(0))

	if err != nil {
		logger.Fatalf("Cannot build connection to server %s, error: %v", addr, err)
	}

	return cm.connMap[addr]
}

func (cm *CommManager) GetConnection(addr string) (*grpc.ClientConn, bool) {
	cm.connMapLock.RLock()
	defer cm.connMapLock.RUnlock()

	conn, exists := cm.connMap[addr]
	return conn, exists
}

// Builds an RPC stub if not exists. Otherwise, returns the existing one
func (cm *CommManager) BuildRpcStub(addr string) DynamicPaxosClient {
	cm.rpcStubMapLock.Lock()
	defer cm.rpcStubMapLock.Unlock()

	if rpcStub, exists := cm.rpcStubMap[addr]; exists {
		return rpcStub
	}

	conn := cm.BuildConnection(addr)

	cm.rpcStubMap[addr] = NewDynamicPaxosClient(conn)

	return cm.rpcStubMap[addr]
}

// Returns the RPC stub on file
func (cm *CommManager) GetRpcStub(addr string) (DynamicPaxosClient, bool) {
	cm.rpcStubMapLock.RLock()
	defer cm.rpcStubMapLock.RUnlock()

	rpcStub, exists := cm.rpcStubMap[addr]
	return rpcStub, exists
}

func (cm *CommManager) NewRpcStub(addr string) DynamicPaxosClient {
	conn, exists := cm.GetConnection(addr)
	if !exists {
		conn = cm.BuildConnection(addr)
	}

	rpcStub := NewDynamicPaxosClient(conn)

	return rpcStub
}
