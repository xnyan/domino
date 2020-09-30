package node

import (
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Server communication table
type CommManager struct {
	connMap     map[string]*grpc.ClientConn // Network addr (ip:port) --> RPC connection
	connMapLock sync.RWMutex
}

func NewCommManager() *CommManager {
	return &CommManager{
		connMap: make(map[string]*grpc.ClientConn),
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
	for i := 0; err != nil && i < 10; i++ {
		t := time.Duration(100 * 1000 * 1000) // 100 ms
		logger.Errorf("Cannot build connection to server %s, error: %v", addr, err)
		logger.Errorf("Wait for %v to retry...", t)
		time.Sleep(t)
		logger.Errorf("Retry count %d", i)
		cm.connMap[addr], err = grpc.Dial(addr,
			grpc.WithInsecure(), grpc.WithReadBufferSize(0), grpc.WithWriteBufferSize(0))
	}
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

func (cm *CommManager) NewRpcStub(addr string) LatencyClient {
	conn, exists := cm.GetConnection(addr)
	if !exists {
		conn = cm.BuildConnection(addr)
	}

	rpcStub := NewLatencyClient(conn)

	return rpcStub
}
