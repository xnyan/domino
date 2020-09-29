package dynamic

import (
	"io"
	"time"

	"golang.org/x/net/context"
)

func (s *Server) PaxosPropose(ctx context.Context, req *PaxosProposeReq) (*PaxosProposeReply, error) {
	//start := time.Now()
	//logger.Infof("cmdId = %s", req.Cmd.Id)
	isCommit := s.paxos.PaxosLeaderAccept(req.Cmd)
	//logger.Infof("commit cmdId = %s in %v", req.Cmd.Id, time.Now().Sub(start))

	execRet := ""
	if isCommit && s.IsExecReply {
		// Waits for execution result
		execRet = s.em.WaitExecRet(req.Cmd.Id)
	}

	rep := &PaxosProposeReply{
		IsCommit: isCommit,
		ExecRet:  execRet,
	}
	//logger.Infof("return cmdId = %s in %v", req.Cmd.Id, time.Now().Sub(start))
	return rep, nil
}

func (s *Server) FpPropose(req *FpProposeReq, stream DynamicPaxos_FpProposeServer) error {
	//now := time.Now().UnixNano()
	//logger.Infof("cmdId = %s, curTime = %d, timestamp = %d, future = %d ", req.Cmd.Id, now, req.Time.Time, now-req.Time.Time)

	isAccept := s.paxos.FpReplicaAccept(req.Cmd, req.Time)

	//logger.Infof("cmdId = %s, timestamp = %v, isFastAccept = %t", req.Cmd.Id, req.Time, isAccept)

	fastReply := &FpProposeReply{
		IsAccept: isAccept,
		IsFast:   true,
	}
	if err := stream.Send(fastReply); err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending fast-path reply.")
	}

	if s.IsFpLeaderLearner { // Only the leader is the learner
		if req.Time.Shard == s.paxos.GetCurFpShard() {
			s.waitConsensusAndExec(req, stream)
		}
	} else {
		// blocking
		isCommit, isFast := s.paxos.FpWaitConsensus(req.Time.Time, req.Cmd.Id)
		//logger.Infof("cmdId = %s, timestamp = %d, isCommit = %t, isFast = %t",
		//	req.Cmd.Id, req.Time.Time, isCommit, isFast)

		/*
			// Only the specific exec Replica and the Leader will return execution
			// results if the fast path succeeds.
			if isCommit && isFast {
				// The fast path commits the command.
				// The client-chosen replica sends commit result and execution result to the client
				if req.IsExecReply {
					s.returnReply(req, isCommit, isFast, stream)
				}
			} else if isCommit && !isFast {
				// The slow path commits the command.
				// The leader and the chosen replica both send commit result and
				// execution result to the client. The client will choose the first one.
				if req.Time.Shard == s.paxos.GetCurFpShard() || req.IsExecReply {
					s.returnReply(req, isCommit, isFast, stream)
				}
			} else if !isCommit {
				// The client chosen replica returns the result to the client
				if req.IsExecReply {
					if s.IsFpLeaderUsePaxos && s.paxos.GetCurPaxosShard() >= 0 {
						isCommit = s.paxos.PaxosLeaderAccept(req.Cmd)
					}
					s.returnReply(req, isCommit, isFast, stream)
				}
			}
		*/

		// Every replica returns an execution result, and the client just takes the
		// first one. This approach avoids having the client to select an execution
		// replica, but it increases network traffic.
		if isCommit && isFast {
			// The fast path commits the command.
			s.returnReply(req, isCommit, isFast, stream)
		} else if isCommit && !isFast {
			// The slow path commits the command.
			s.returnReply(req, isCommit, isFast, stream)
		} else if !isCommit {
			// The client chosen replica returns the result to the client
			if req.IsExecReply {
				if s.IsFpLeaderUsePaxos && s.paxos.GetCurPaxosShard() >= 0 {
					isCommit = s.paxos.PaxosLeaderAccept(req.Cmd)
				}
				s.returnReply(req, isCommit, isFast, stream)
			}
		}
	}

	//logger.Infof("FpPropose RPC returns cmd = %s", req.Cmd.Id)

	return nil
}

func (s *Server) waitConsensusAndExec(
	req *FpProposeReq, stream DynamicPaxos_FpProposeServer,
) {
	// Fast Paxos shard leader (coordinator) waits for consensus result
	isCommit, isFast := s.paxos.FpWaitConsensus(req.Time.Time, req.Cmd.Id) // blocking

	//logger.Infof("cmdId = %s consensus result isCommit = %t isFast = %t", req.Cmd.Id, isCommit, isFast)

	if s.IsFpLeaderUsePaxos {
		if !isCommit && s.paxos.GetCurPaxosShard() >= 0 {
			// Uses Paxos shard to accept a command that is rejected by Fast Paxos
			isCommit = s.paxos.PaxosLeaderAccept(req.Cmd)
		}
	}

	ret := ""
	if isCommit && s.IsExecReply {
		// Waits for execution result
		ret = s.em.WaitExecRet(req.Cmd.Id)
	}

	// Sends slow-path and/or execution reply to the client
	rep := &FpProposeReply{IsAccept: isCommit, IsFast: isFast, ExecRet: ret}
	if err := stream.Send(rep); err != nil {
		logger.Fatalf("Fails sending consensus result for cmdId = %s, t = %v, error %v",
			req.Cmd.Id, req.Time, err)
	}
}

func (s *Server) returnReply(
	req *FpProposeReq, isCommit, isFast bool, stream DynamicPaxos_FpProposeServer,
) {
	ret := ""
	//start := time.Now()
	if s.IsExecReply && isCommit {
		ret = s.em.WaitExecRet(req.Cmd.Id) // Waits for execution result
	}
	//duration := time.Now().Sub(start)
	//logger.Infof("cmdId = %s isCommit = %t execDuration = %v", req.Cmd.Id, isCommit, duration)
	rep := &FpProposeReply{IsAccept: isCommit, IsFast: isFast, ExecRet: ret}
	if err := stream.Send(rep); err != nil {
		logger.Fatalf("Fails sending consensus result for cmdId = %s, t = %v, error %v",
			req.Cmd.Id, req.Time, err)
	}
}

// Latency probe for clients
func (s *Server) Probe(ctx context.Context, req *ProbeReq) (*ProbeReply, error) {
	qDelay := s.paxos.Probe()
	paxosLat := s.PredictPaxosLat() // Predicts the Paxos latency (in ms) when this is the leader
	return &ProbeReply{QueuingDelay: qDelay.Nanoseconds(), PaxosLat: int32(paxosLat)}, nil
}

// Time offset probe for clients.
/*
To consider the clock synchronization between clients and servers, a client predicts
the time when its request can be processed on each server.
Each server returns its clock time to the client, and the client can calculate the
offset between its clock time and the server's clock time, which includes the
network delay and queuing delay.
The client will use the time offset to set its request timestamp.
*/
func (s *Server) ProbeTime(ctx context.Context, req *ProbeReq) (*ProbeTimeReply, error) {
	s.paxos.Probe()
	t := time.Now().UnixNano()
	paxosLat := s.PredictPaxosLat() // Predicts the Paxos latency (in ms) when this is the leader
	return &ProbeTimeReply{ProcessTime: t, PaxosLat: int32(paxosLat)}, nil
}

func (s *Server) ReplicaProbe(ctx context.Context, req *ProbeReq) (*ReplicaProbeReply, error) {
	qDelay := s.paxos.Probe()
	return &ReplicaProbeReply{QueuingDelay: qDelay.Nanoseconds()}, nil
}

func (s *Server) DeliverReplicaMsg(stream DynamicPaxos_DeliverReplicaMsgServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			logger.Fatalf("Error: %v", err)
		}

		s.paxos.ProcessReplicaMsg(msg)
	}
	stream.SendAndClose(&EmptyReply{})
	return nil
}

func (s *Server) Test(ctx context.Context, req *TestReq) (*TestReply, error) {
	s.doTest()
	rep := &TestReply{}
	return rep, nil
}
