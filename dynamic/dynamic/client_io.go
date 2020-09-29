package dynamic

import (
	"io"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ClientIo interface {
	InitConn(addrList []string)
	// Paxos, blocking RPC function
	SendPaxosProposeReq(addr string, cmd *Command) *PaxosProposeReply

	// Fast Paxos, non-blocking RPC functions
	BcstFpProposeToFollowers(addrList []string, req *FpProposeReq, fastC chan<- *FpProposeReply, wg *sync.WaitGroup)
	SendFpProposeToLeader(addr string, req *FpProposeReq, fastC, slowC chan<- *FpProposeReply, wg *sync.WaitGroup)
	// Fast Paxos with all replicas as learners, non-blocking RPC functions
	SendFpProposeToExecReplica(addr string, req *FpProposeReq, fastC, slowC chan<- *FpProposeReply, wg *sync.WaitGroup)

	// Latency monitoring
	SendProbeReq(addr string) *ProbeReply
	SendProbeTimeReq(addr string) *ProbeTimeReply

	// Testing
	SendTestReq(addr string)
}

type ClientGrpc struct {
	cm *CommManager
}

func NewClientGrpc() ClientIo {
	return &ClientGrpc{
		cm: NewCommManager(),
	}
}

func (gIo *ClientGrpc) InitConn(addrList []string) {
	var wg sync.WaitGroup
	for _, addr := range addrList {
		wg.Add(1)
		go func(addr string) {
			gIo.cm.BuildConnection(addr)
			wg.Done()
		}(addr)
	}
	wg.Wait()
}

////////// RPCs for clients ///////////////
func (gIo *ClientGrpc) SendProbeReq(addr string) *ProbeReply {
	req := &ProbeReq{}
	rpcStub := gIo.cm.NewRpcStub(addr)
	reply, err := rpcStub.Probe(context.Background(), req, grpc.WaitForReady(true))
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending probe request to addr = %s", addr)
	}
	return reply
}

func (gIo *ClientGrpc) SendProbeTimeReq(addr string) *ProbeTimeReply {
	req := &ProbeReq{}
	rpcStub := gIo.cm.NewRpcStub(addr)
	reply, err := rpcStub.ProbeTime(context.Background(), req, grpc.WaitForReady(true))
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending probe request to addr = %s", addr)
	}
	return reply
}

func (gIo *ClientGrpc) SendPaxosProposeReq(
	addr string, cmd *Command,
) *PaxosProposeReply {
	req := &PaxosProposeReq{Cmd: cmd}
	rpcStub := gIo.cm.NewRpcStub(addr)
	//reply, err := rpcStub.PaxosPropose(context.Background(), req)
	reply, err := rpcStub.PaxosPropose(context.Background(), req, grpc.WaitForReady(true))

	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending paxos proposal to addr = %s, cmdId = %s", addr, cmd.Id)
	}
	return reply
}

// Non-blocking
func (gIo *ClientGrpc) BcstFpProposeToFollowers(
	addrList []string, req *FpProposeReq, fastC chan<- *FpProposeReply, wg *sync.WaitGroup,
) {
	for _, addr := range addrList {
		wg.Add(1)
		go func(addr string) {
			reply, _ := gIo.sendFpProposeReq(addr, req)
			fastC <- reply
			wg.Done()
		}(addr)
	}
}

// Non-blocking
func (gIo *ClientGrpc) SendFpProposeToLeader(
	leaderAddr string, req *FpProposeReq, fastC, slowC chan<- *FpProposeReply, wg *sync.WaitGroup,
) {
	wg.Add(1)
	go func() {
		fastReply, stream := gIo.sendFpProposeReq(leaderAddr, req)
		fastC <- fastReply

		slowReply, err := stream.Recv()
		if err != nil {
			// TODO Does nothing if the leader does not return slow-path result, where err == io.EOF
			logger.Fatalf("Fails to recv slow-path reply from addr = %s req = %v error %v",
				leaderAddr, req, err)
		}

		slowC <- slowReply

		wg.Done()
	}()
}

// Non-blocking
func (gIo *ClientGrpc) SendFpProposeToExecReplica(
	addr string, req *FpProposeReq, fastC, slowC chan<- *FpProposeReply, wg *sync.WaitGroup,
) {
	wg.Add(1)
	go func() {
		fastReply, stream := gIo.sendFpProposeReq(addr, req)
		fastC <- fastReply

		//logger.Infof("Waiting slow reply from addr = %s cmd = %s", addr, req.Cmd.Id)
		if slowReply, err := stream.Recv(); err == nil {
			slowC <- slowReply
		} else if err != io.EOF {
			logger.Fatalf("Fails to recv slow-path reply from addr = %s req = %v error %v", addr, req, err)
		}
		//logger.Infof("Done waiting slow reply from addr = %s cmd = %s", addr, req.Cmd.Id)

		wg.Done()
	}()
}

// Blocking
func (gIo *ClientGrpc) sendFpProposeReq(
	addr string, req *FpProposeReq,
) (*FpProposeReply, DynamicPaxos_FpProposeClient) {
	rpcStub := gIo.cm.NewRpcStub(addr)
	//logger.Infof("Propose to addr = %s cmdId = %s t = %d real = %d", addr, req.Cmd.Id, req.Time, time.Now().UnixNano())
	stream, err := rpcStub.FpPropose(context.Background(), req, grpc.WaitForReady(true))

	if err != nil {
		logger.Fatalf("Fails establish fast-paxos stream to addr = %s req = %v error %v", addr, req, err)
	}

	reply, err := stream.Recv()
	if err != nil {
		logger.Fatalf("Fails to recv fast-path reply from addr = %s req = %v error %v", addr, req, err)
	}
	//logger.Infof("Reply from addr = %s cmdId = %s t = %d actual t = %d", addr, req.Cmd.Id, req.Time, reply.Time)

	return reply, stream
}

// Testing
func (gIo *ClientGrpc) SendTestReq(addr string) {
	req := &TestReq{}
	rpcStub := gIo.cm.NewRpcStub(addr)
	_, err := rpcStub.Test(context.Background(), req)

	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending test req to addr = %s", addr)
	}
}
