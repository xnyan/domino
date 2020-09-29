package main

import (
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"domino/common"
	fp "domino/fastpaxos/fastpaxos"
	"domino/fastpaxos/rpc"
)

////////////////////////////////////////
// The leader asks every replica to commit op at idx
// Non-blocking
func (server *Server) LeaderCommitOp(op *rpc.Operation, idx string) *common.Future {
	req := &rpc.CommitRequest{
		Idx: idx,
		Op:  op,
	}

	// Async RPC to followers
	commitNum := 0
	var commitLock sync.Mutex
	commitCv := sync.NewCond(&commitLock)
	for _, addr := range server.FollowerAddrList {
		go func(addr string) {
			server.SendCommitRequest(addr, req)

			commitLock.Lock()
			commitNum++
			commitLock.Unlock()
			commitCv.Signal()
		}(addr)
	}

	// Async local commit
	go func() {
		server.CommitOp(op, idx).GetValue()
		commitLock.Lock()
		commitNum++
		commitLock.Unlock()
		commitCv.Signal()
	}()

	done := common.NewFuture()

	go func() {
		commitLock.Lock()
		for commitNum < server.MajorityNum {
			commitCv.Wait()
		}
		commitLock.Unlock()
		done.SetValue(true)
	}()

	return done
}

// Blocking
func (server *Server) SendCommitRequest(addr string, req *rpc.CommitRequest) *rpc.CommitReply {
	rpcStub := server.GetRpcStub(addr)

	reply, err := rpcStub.Commit(context.Background(), req, grpc.WaitForReady(true))
	//reply, err := rpcStub.Commit(context.Background(), req)

	if err != nil {
		logger.Fatalf("Leader fails to commit operaiton id = %s at idx = %d to server addr = %s err = %v",
			req.Op.Id, req.Idx, addr, err)
	}

	return reply
}

// May block
func (server *Server) CommitOp(op *rpc.Operation, idx string) *common.Future {
	done := common.NewFuture()
	lc := &fp.LeaderCommit{Idx: fp.ParseLogIdx(idx), Op: op, Done: done}

	server.fp.RunCmd(lc)

	return done
}

/////////////////////////////////
// The leader asks every replica to accept the operation at idx via the slow path
// Non-blocking
func (server *Server) LeaderAcceptOp(op *rpc.Operation, idx string) *common.Future {
	req := &rpc.AcceptRequest{
		Idx: idx,
		Op:  op,
	}

	// Async RPC to followers
	acceptNum := 0
	var acceptLock sync.Mutex
	acceptCv := sync.NewCond(&acceptLock)
	for _, addr := range server.FollowerAddrList {
		go func(addr string) {
			server.SendAcceptRequest(addr, req)

			acceptLock.Lock()
			acceptNum++
			acceptLock.Unlock()
			acceptCv.Signal()
		}(addr)
	}

	// Async local accept
	go func() {
		server.AcceptOp(op, idx).GetValue()
		acceptLock.Lock()
		acceptNum++
		acceptLock.Unlock()
		acceptCv.Signal()
	}()

	done := common.NewFuture()

	go func() {
		acceptLock.Lock()
		for acceptNum < server.MajorityNum {
			acceptCv.Wait()
		}
		acceptLock.Unlock()
		done.SetValue(true)
	}()

	return done
}

// Blocking
func (server *Server) SendAcceptRequest(addr string, req *rpc.AcceptRequest) *rpc.AcceptReply {
	rpcStub := server.GetRpcStub(addr)

	reply, err := rpcStub.Accept(context.Background(), req, grpc.WaitForReady(true))
	//reply, err := rpcStub.Accept(context.Background(), req)

	if err != nil {
		logger.Fatalf("Leader fails to send its slow-path accepted operation id = %s at idx = %d to "+
			"server addr = %s error = %v", req.Op.Id, req.Idx, addr, err)
	}

	return reply
}

// May block
func (server *Server) AcceptOp(op *rpc.Operation, idx string) *common.Future {
	done := common.NewFuture()
	la := &fp.LeaderProposal{Idx: fp.ParseLogIdx(idx), Op: op, Done: done}

	server.fp.RunCmd(la)

	return done
}

////////////////////////////////////////
// A follower sends its fast-path accepted operation at idx to the leader
// Blocking
func (server *Server) SendPromiseRequest(req *rpc.PromiseRequest) *rpc.PromiseReply {
	rpcStub := server.GetRpcStub(server.LeaderAddr)

	reply, err := rpcStub.Promise(context.Background(), req, grpc.WaitForReady(true))
	//reply, err := rpcStub.Promise(context.Background(), req)

	if err != nil {
		logger.Fatalf("Follower fails to send its fast-path accepted operaiton id = %s at idx = %d to "+
			"leader addr = %s error = %v", req.Op.Id, req.Idx, server.LeaderAddr, err)
	}

	return reply
}

////////////////////////////////////////
// Helper functions
func (server *Server) GetRpcStub(addr string) rpc.FastPaxosRpcClient {
	conn, exists := server.cm.GetConnection(addr)
	if !exists {
		conn = server.cm.BuildConnection(addr)
	}

	rpcStub := rpc.NewFastPaxosRpcClient(conn)

	return rpcStub
}
