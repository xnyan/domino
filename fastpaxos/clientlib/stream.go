package clientlib

import (
	"io"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"domino/fastpaxos/rpc"
)

type replyHandle struct {
	fastC chan<- *rpc.ProposeReply
	slowC chan<- *rpc.ProposeReply
}

var replyHandleTable map[string]*replyHandle = make(map[string]*replyHandle)
var replyHandleLock sync.Mutex

func initReplyHandle(opId string, fastC, slowC chan<- *rpc.ProposeReply) {
	replyHandleLock.Lock()
	defer replyHandleLock.Unlock()
	_, ok := replyHandleTable[opId]
	if ok {
		logger.Fatalf("There is already a reply handle for operation id = %s", opId)
	}
	replyHandleTable[opId] = &replyHandle{
		fastC: fastC,
		slowC: slowC,
	}
}

func getReplyHandle(opId string) *replyHandle {
	replyHandleLock.Lock()
	defer replyHandleLock.Unlock()
	if handle, ok := replyHandleTable[opId]; ok {
		return handle
	}
	return nil
}

var sendChTable map[string]chan *rpc.ProposeRequest = make(map[string]chan *rpc.ProposeRequest)
var replyCh chan *rpc.ProposeReply = make(chan *rpc.ProposeReply, 1024*1024*4)

func (c *Client) initStream(addrList []string) {
	var wg sync.WaitGroup
	for _, addr := range addrList {
		sendChTable[addr] = make(chan *rpc.ProposeRequest, 1024*1024*4) // TODO configurable size
		wg.Add(1)
		go func(addr string) {
			c.getStreamHandle(addr)
			wg.Done()
		}(addr)
	}
	wg.Wait()

	for _, addr := range addrList {
		// Starts sending streams
		go func(addr string) {
			stream := c.getStreamHandle(addr)
			for req := range sendChTable[addr] {
				err := stream.Send(req)
				if err != nil {
					logger.Fatalf("Fails to send request to server addr = %s", addr)
				}
			}
		}(addr)

		// Starts receving streams
		go func(addr string) {
			stream := c.getStreamHandle(addr)
			for {
				reply, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					logger.Fatalf("Fails to receive reply from server addr = %s", addr)
				}
				replyCh <- reply
			}
		}(addr)
	}

	// Reply dispatcher
	go func() {
		for reply := range replyCh {
			handle := getReplyHandle(reply.OpId)
			if handle != nil {
				if reply.IsSlow {
					handle.slowC <- reply
				} else {
					handle.fastC <- reply
				}
			}
		}
	}()
}

func (c *Client) streamSendProposeReq(
	addr string,
	req *rpc.ProposeRequest,
) {
	ch, ok := sendChTable[addr]
	if !ok {
		logger.Fatalf("There is no stream shending channel for server addr = %s", addr)
	}
	ch <- req
}

/////// Streaming RPC I/O ////////
// Non-blocking
func (c *Client) streamSendProposeReqToLeader(
	leaderAddr string, req *rpc.ProposeRequest, fastC, slowC chan<- *rpc.ProposeReply, wg *sync.WaitGroup,
) {
	initReplyHandle(req.Op.Id, fastC, slowC)
	c.streamSendProposeReq(leaderAddr, req)
}

// Non-blocking
func (c *Client) streamBcstProposeReqToFollowers(
	followerAddrList []string, req *rpc.ProposeRequest, fastC chan<- *rpc.ProposeReply, wg *sync.WaitGroup,
) {
	for _, addr := range followerAddrList {
		c.streamSendProposeReq(addr, req)
	}
}

//// Stream Handle
var streamTable map[string]rpc.FastPaxosRpc_StreamProposeClient = make(map[string]rpc.FastPaxosRpc_StreamProposeClient)
var streamTableLock sync.Mutex

func (c *Client) getStreamHandle(addr string) rpc.FastPaxosRpc_StreamProposeClient {
	streamTableLock.Lock()
	defer streamTableLock.Unlock()
	stream, ok := streamTable[addr]
	if !ok {
		rpcStub := c.GetRpcStub(addr)
		var err error
		stream, err = rpcStub.StreamPropose(context.Background(), grpc.WaitForReady(true))
		if err != nil {
			logger.Fatalf("Fails to start RPC streaming to server addr = %s", addr)
		}
		streamTable[addr] = stream
	}
	return stream
}
