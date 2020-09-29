package main

import (
	"golang.org/x/net/context"

	"domino/fastpaxos/rpc"
)

// Testing
func (server *Server) Test(
	ctx context.Context,
	request *rpc.TestRequest,
) (*rpc.TestReply, error) {

	server.test()

	//server.kvStore.Print()
	server.fp.Test()

	return &rpc.TestReply{}, nil
}

func (server *Server) test() {
	server.execRetTableLock.RLock()
	defer server.execRetTableLock.RUnlock()

	logger.Infof("execRetTable size = %d", len(server.execRetTable))
	for t, _ := range server.execRetTable {
		logger.Info("execRet t = %d", t)
	}
}
