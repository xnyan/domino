package dynamic

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RpcIo interface {
	InitConn(addrList []string, streamBufSize int)

	// FIFO communication and processing between replicas
	BcstReplicaMsg(addrList []string, msg *ReplicaMsg)
	SendReplicaMsg(addr string, msg *ReplicaMsg)

	// Synchoronous I/O Latency probing (not FIFO processing)
	SendReplicaProbeReq(addr string) *ReplicaProbeReply

	// Helper function
	syncSendReplicaMsg(addr string, msg *ReplicaMsg)
}

type rpcIo struct {
	streamChMap map[string]chan<- *ReplicaMsg
	syncSend    bool // synchronously sends a message
}

func newRpcIo(isSyncSend bool) *rpcIo {
	return &rpcIo{
		streamChMap: make(map[string]chan<- *ReplicaMsg),
		syncSend:    isSyncSend,
	}
}

func (io *rpcIo) setStreamCh(addr string, c chan<- *ReplicaMsg) {
	io.streamChMap[addr] = c
}

func (io *rpcIo) getStreamCh(addr string) chan<- *ReplicaMsg {
	c, ok := io.streamChMap[addr]
	if !ok {
		logger.Fatalf("Cannot find stream sender channel to addr = %s", addr)
	}
	return c
}

func (io *rpcIo) closeStream(addr string) {
	c := io.streamChMap[addr]
	close(c)
	delete(io.streamChMap, addr)
}

///////////////////////////////////
// gRPC streaming implementation
type GrpcIo struct {
	*rpcIo
	cm            *CommManager
	syncStreamMap map[string]DynamicPaxos_DeliverReplicaMsgClient
}

func NewGrpcIo(isSyncSend bool) RpcIo {
	io := &GrpcIo{
		rpcIo:         newRpcIo(isSyncSend),
		cm:            NewCommManager(),
		syncStreamMap: make(map[string]DynamicPaxos_DeliverReplicaMsgClient),
	}
	return io
}

func (io *GrpcIo) InitConn(addrList []string, streamBufSize int) {
	for _, addr := range addrList {
		io.initStream(addr, streamBufSize)
	}
}

func (io *GrpcIo) BcstReplicaMsg(addrList []string, msg *ReplicaMsg) {
	for _, addr := range addrList {
		io.SendReplicaMsg(addr, msg)
	}
}

func (io *GrpcIo) SendReplicaMsg(addr string, msg *ReplicaMsg) {
	if io.syncSend {
		io.syncSendReplicaMsg(addr, msg)
	} else {
		c := io.getStreamCh(addr)
		c <- msg
	}
}

// Synchoronous I/O
func (io *GrpcIo) SendReplicaProbeReq(addr string) *ReplicaProbeReply {
	req := &ProbeReq{}
	rpcStub := io.cm.NewRpcStub(addr)
	reply, err := rpcStub.ReplicaProbe(context.Background(), req, grpc.WaitForReady(true))
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending probe request to addr = %s", addr)
	}
	return reply
}

// Returns after sending the message
func (io *GrpcIo) syncSendReplicaMsg(addr string, msg *ReplicaMsg) {
	stream, ok := io.syncStreamMap[addr]
	if !ok {
		logger.Fatalf("Cannot find the stream for addr = %s", addr)
	}
	err := stream.Send(msg)
	if err != nil {
		logger.Fatalf("Fails to send msg = %v to addr = %s", msg, addr)
	}
}

// Inits a streaming to a replica server
func (io *GrpcIo) initStream(addr string, bufferSize int) {
	rpcStub := io.cm.NewRpcStub(addr)
	stream, err := rpcStub.DeliverReplicaMsg(context.Background(), grpc.WaitForReady(true))
	if err != nil {
		logger.Fatalf("Fails to init stream to addr = %s, err = %v", addr, err)
	}

	if io.syncSend {
		io.setSyncStream(addr, stream)
	} else {
		c := make(chan *ReplicaMsg, bufferSize)
		io.setStreamCh(addr, c)

		// Makes sure that there is only one thread sending msgs via the stream.
		go func(c <-chan *ReplicaMsg, stream DynamicPaxos_DeliverReplicaMsgClient) {
			for msg := range c {
				err := stream.Send(msg)
				if err != nil {
					logger.Fatalf("Fails to send msg = %v", msg)
				}
			}
			stream.CloseSend()
		}(c, stream)
	}
}

func (io *GrpcIo) setSyncStream(addr string, stream DynamicPaxos_DeliverReplicaMsgClient) {
	io.syncStreamMap[addr] = stream
}
