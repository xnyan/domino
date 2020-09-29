package main

import (
	"errors"
	"io"
	"time"

	"golang.org/x/net/context"

	"domino/common"
	fp "domino/fastpaxos/fastpaxos"
	"domino/fastpaxos/rpc"
)

// RPCs on replicas for clients to propose an operation
func (server *Server) Propose(req *rpc.ProposeRequest, stream rpc.FastPaxosRpc_ProposeServer) error {

	idxFuture := common.NewFuture()
	cp := &fp.ClientProposal{
		Op:        req.Op,
		Idx:       idxFuture,
		Delay:     time.Duration(req.Delay),
		Timestamp: req.Timestamp,
		ClientId:  req.ClientId,
	}

	// May block, which depends on the scheduling strategy
	server.fp.Schedule(cp)
	//server.fp.RunCmd(cp) // May block if fast paxos' input buffer is full

	idx := idxFuture.GetValue().(*fp.LogIdx).String() // Blocking

	reply := &rpc.ProposeReply{OpId: req.Op.Id, Idx: idx, IsSlow: false}

	// Sends response to the client
	if err := stream.Send(reply); err != nil {
		// It is possible that the client has a successful fast path from other
		// replicas, and it closes the communication.
		// In this case, cannot terminate the RPC, as the replica needs to vote for
		// the system to continue.
		logger.Errorf("Fails to send the fast-path reply. operation id = %s error: %v",
			req.Op.Id, err)
		// TODO Distinguish the above case with other types of errors. Use a more
		// elegant way instead of just printing an error.
	}

	if server.IsLeader {
		// The leader waits for the completion of the fast path
		vote := &fp.Vote{Idx: idx, OpId: req.Op.Id}

		logger.Debugf("Leader votes opId = (%s) at idx = (%s)", vote.OpId, vote.Idx)

		fastRet, slowRet := server.fp.LeaderVote(vote)
		defer server.fp.CleanRetHandle(vote.OpId)

		logger.Debugf("Leader watis for the fast-path result for opId = (%s)", vote.OpId)

		fastIdx := fastRet.GetValue().(string) // blocking

		logger.Debugf("Leader gets the fast-path result for opId = (%s): idx = (%s)",
			vote.OpId, fastIdx)

		if fastIdx != fp.INVALID_IDX {
			// The fast path succeeds
			// Asks replicas to commit the operation at the idx
			server.LeaderCommitOp(req.Op, fastIdx)
		} else {
			// The fast path fails, and the leader starts a slow path for that operation.

			// The slow-path will choose an idx for the operation. The leader
			// must make sure that the operation will be put on an idx that will
			// not have any other operation succeed via the fast path. This the
			// safety requirment. Otherwise, a client may learn an operation at
			// the idx via the fast path, but the leader puts a different
			// operation at the idx via the slow path.  (More thoughts are needed
			// for fault tolerance. If use a timeout to reject the operation, how
			// to put a no-op on the potential empty idx)

			logger.Debugf("Leader watis for the slow-path result for opId = (%s)", vote.OpId)

			slowIdx := slowRet.GetValue().(string) // blocking

			logger.Debugf("Leader gets the slow-path result for opId = (%s): idx = (%s)",
				vote.OpId, slowIdx)

			if slowIdx == fp.INVALID_IDX {
				logger.Fatalf("Invalid slow-path idx for operation id = %s", req.Op.Id)
			}

			// Asks replicas to accept the operation at the chosen idx via the slow-path
			isSlowAcceptDone := server.LeaderAcceptOp(req.Op, slowIdx)

			// Blocking
			if isSlowAcceptDone.GetValue().(bool) {
				// The slow path completes replicating the operation at the chosen
				// idx on at least a majority of the replicas

				slowPathReply := &rpc.ProposeReply{
					OpId:   req.Op.Id,
					Idx:    slowIdx,
					IsSlow: true,
				}

				// Replies the slow-path resultto the client
				if err := stream.Send(slowPathReply); err != nil {
					logger.Fatalf("Fails to send the slow-path reply. operation id = %s error: %v", req.Op.Id, err)
					return err
				}

				// Async asks replicas to commit the operation at the chosen idx
				go server.LeaderCommitOp(req.Op, slowIdx)
			} else {
				logger.Fatalf("Slow path replicaiton should not fail operation id = %s idx = %d",
					req.Op.Id, slowIdx)
			}
		}
	} else {
		// A follower sends the leader the idx that it accepts the operation at via the fast path
		promiseReq := &rpc.PromiseRequest{Idx: idx, Op: req.Op}
		// Async sends the accepted result to the leader
		go server.SendPromiseRequest(promiseReq)
	}

	return nil
}

// NOTE This streaming RPC currently only works for the BASIC FAST PAXOS without any scheduling
func (server *Server) StreamPropose(stream rpc.FastPaxosRpc_StreamProposeServer) error {
	replyClientCh := make(chan *rpc.ProposeReply, 1024*1024*4) // TODO configurable size
	defer close(replyClientCh)

	go func() {
		for reply := range replyClientCh { // Reply to the client
			stream.Send(reply)
		}
	}()

	for {
		// Processing client proposals
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			logger.Errorf("Stream RPC fails to receive requests. Error: %v", err)
			return err
		}

		idxFuture := common.NewFuture()
		cp := &fp.ClientProposal{
			Op:        req.Op,
			Idx:       idxFuture,
			Delay:     time.Duration(req.Delay),
			Timestamp: req.Timestamp,
			ClientId:  req.ClientId,
		}
		server.fp.Schedule(cp) // May block
		go func(req *rpc.ProposeRequest, idxFuture *common.Future, cp *fp.ClientProposal) {
			idx := idxFuture.GetValue().(*fp.LogIdx).String() // Blocking
			reply := &rpc.ProposeReply{OpId: req.Op.Id, Idx: idx, IsSlow: false}
			replyClientCh <- reply // Sends fast-path response to the client

			if server.IsLeader {
				// The leader waits for the completion of the fast path
				vote := &fp.Vote{Idx: idx, OpId: req.Op.Id}
				fastRet, slowRet := server.fp.LeaderVote(vote)
				defer server.fp.CleanRetHandle(vote.OpId)
				fastIdx := fastRet.GetValue().(string) // blocking
				if fastIdx != fp.INVALID_IDX {
					server.LeaderCommitOp(req.Op, fastIdx)
				} else {
					// The fast path fails, and the leader starts a slow path for that operation.
					slowIdx := slowRet.GetValue().(string) // blocking
					if slowIdx == fp.INVALID_IDX {
						logger.Fatalf("Invalid slow-path idx for operation id = %s", req.Op.Id)
					}
					// Asks replicas to accept the operation at the chosen idx via the slow-path
					isSlowAcceptDone := server.LeaderAcceptOp(req.Op, slowIdx)
					// Blocking
					if isSlowAcceptDone.GetValue().(bool) {
						slowPathReply := &rpc.ProposeReply{
							OpId:   req.Op.Id,
							Idx:    slowIdx,
							IsSlow: true,
						}
						replyClientCh <- slowPathReply // Replies the slow-path resultto the client
						// Async asks replicas to commit the operation at the chosen idx
						go server.LeaderCommitOp(req.Op, slowIdx)
					} else {
						logger.Fatalf("Slow path replicaiton should not fail operation id = %s idx = %d",
							req.Op.Id, slowIdx)
					}
				}
			} else {
				// A follower sends the leader the idx that it accepts the operation at via the fast path
				promiseReq := &rpc.PromiseRequest{Idx: idx, Op: req.Op}
				// Async sends the accepted result to the leader
				go server.SendPromiseRequest(promiseReq)
			}
		}(req, idxFuture, cp)
	}

	return nil
}

// A replica commits an operation at an idx
func (server *Server) Commit(
	ctx context.Context,
	request *rpc.CommitRequest,
) (*rpc.CommitReply, error) {
	done := server.CommitOp(request.Op, request.Idx)
	done.GetValue()
	return &rpc.CommitReply{}, nil
}

// A replica accepts an operation at an idx via the slow path
func (server *Server) Accept(
	ctx context.Context,
	request *rpc.AcceptRequest,
) (*rpc.AcceptReply, error) {
	done := server.AcceptOp(request.Op, request.Idx)
	done.GetValue()
	return &rpc.AcceptReply{}, nil
}

// The leader receives follower's accept response via the fast path
func (server *Server) Promise(
	ctx context.Context,
	request *rpc.PromiseRequest,
) (*rpc.PromiseReply, error) {

	if !server.IsLeader {
		logger.Fatalf("This is not a leader. server addr = %s", server.NetAddr)
		return nil, errors.New("Server " + server.NetAddr + " is not a leader")
	}

	server.fp.Vote(&fp.Vote{OpId: request.Op.Id, Idx: request.Idx})

	reply := &rpc.PromiseReply{}

	return reply, nil
}

/////////////////////////
//// Streaming RPCs

// RPCs on replicas for clients to propose an operation
/*
func (server *Server) StreamPropose(stream rpc.Consensus_StreamProposeServer) error {
	return nil

		replyClientCh := make(chan *rpc.ProposeReply, RpcStreamServerToClientBufferSize)
		defer close(replyClientCh)

		go func() {
			for reply := range replyClientCh {
				stream.Send(reply)
			}
		}()

		for {
			req, err := stream.Recv() // Client proposal

			if err == io.EOF {
				return nil
			}

			if err != nil {
				logger.Errorf("Stream RPC fails to receive requests. Error: %v", err)
				return err
			}

			idxFuture := common.NewFuture()
			cp := &fp.ClientProposal{req.Op, idxFuture}

			server.fp.RunCmd(cp) // May block if fast paxos' input buffer is full

			idx := idxFuture.GetValue().(int) // Blocking

			reply := &rpc.ProposeReply{OpId: req.Op.Id, Idx: int64(idx), IsSlow: false}

			// Sends response to the client
			replyClientCh <- reply

			if server.IsLeader {
				// The leader waits for the completion of the fast path
				vote := &fp.Vote{Idx: idx, OpId: req.Op.Id}
				fastRet, slowRet := server.fp.LeaderVote(vote)

				// Uses a thread to avoid blocking handling the client's next operation if
				// the client sees the success of the fast path earlier than the leader does.
				go func() {
					fastIdx := fastRet.GetValue().(int) // blocking
					if fastIdx != fp.INVALID_IDX {
						// The fast path succeeds
						// Asks replicas to commit the operation at the idx
						server.LeaderCommitOp(req.Op, fastIdx)
					} else {
						// The fast path fails, and the leader starts a slow path for that operation.

						// The slow-path will choose an idx for the operation. The leader
						// must make sure that the operation will be put on an idx that will
						// not have any other operation succeed via the fast path. This the
						// safety requirment. Otherwise, a client may learn an operation at
						// the idx via the fast path, but the leader puts a different
						// operation at the idx via the slow path.  (More thoughts are needed
						// for fault tolerance. If use a timeout to reject the operation, how
						// to put a no-op on the potential empty idx)
						slowIdx := slowRet.GetValue().(int) // blocking
						if slowIdx == fp.INVALID_IDX {
							logger.Fatalf("Invalid slow-path idx for operation id = %s", req.Op.Id)
						}

						// Asks replicas to accept the operation at the chosen idx via the slow-path
						isSlowAcceptDone := server.LeaderAcceptOp(req.Op, slowIdx)

						if isSlowAcceptDone.GetValue().(bool) {
							// The slow path completes replicating the operation at the chosen
							// idx on at least a majority of the replicas

							// Replies to the client
							slowPathReply := &rpc.ProposeReply{
								OpId:   req.Op.Id,
								Idx:    int64(slowIdx),
								IsSlow: true,
							}
							replyClientCh <- slowPathReply

							// Asks replicas to commit the operation at the chosen idx
							server.LeaderCommitOp(req.Op, slowIdx)
						} else {
							logger.Fatalf("Slow path replicaiton should not fail operation id = %s idx = %d",
								req.Op.Id, slowIdx)
						}
					}
				}()
			} else {
				// A follower sends the leader the idx that it accepts the operation at via the fast path
				promiseReq := &rpc.PromiseRequest{Idx: int64(idx), Op: req.Op}
				server.SendPromiseRequest(promiseReq)
			}
		}

		return nil

}
*/
