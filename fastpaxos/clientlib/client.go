package clientlib

import (
	"io"
	"strconv"
	"sync"

	"github.com/op/go-logging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"domino/fastpaxos/com"
	fp "domino/fastpaxos/fastpaxos"
	"domino/fastpaxos/rpc"
)

var logger = logging.MustGetLogger("client")

type Client struct {
	Id   string
	DcId string

	leaderAddr       string   //ip:port
	followerAddrList []string // a list of followers' network addresses
	replicaNum       int
	fastQuorum       int

	// Communication manger
	cm *com.CommManager

	// Waiting execution result
	waitExecRet bool

	// Operation count
	opCount int
	lock    sync.Mutex

	wg sync.WaitGroup
}

func NewClient(id, dcId, configFile, replicaFile string) *Client {
	c := &Client{
		Id:               id,
		DcId:             dcId,
		followerAddrList: make([]string, 0),
		cm:               com.NewCommManager(),
		waitExecRet:      false, //TODO implement waiting for execution result
		opCount:          0,
	}
	c.loadConfig(configFile, replicaFile)

	// Using stream I/O
	replicaAddrList := make([]string, len(c.followerAddrList))
	copy(replicaAddrList, c.followerAddrList)
	replicaAddrList = append(replicaAddrList, c.leaderAddr)
	c.initStream(replicaAddrList)

	return c
}

func (c *Client) BuildOperation(rKeyList []string, wTable map[string]string) *rpc.Operation {
	for _, k := range rKeyList {
		return &rpc.Operation{Type: fp.OP_READ, Key: k}
	}

	for k, v := range wTable {
		return &rpc.Operation{Type: fp.OP_WRITE, Key: k, Val: v}
	}

	return nil
}

// Thread safe
// Returns isCommit/isAccept, isFastPathSuccessful, execVal (if applicable)
func (c *Client) Propose(op *rpc.Operation) (bool, bool, string) {
	op.Id = c.genOpId()
	return c.fpPropose(op) // Standard Fast Paxos
}

func (c *Client) Close() {
	c.wg.Wait()
}

// Standard Fast Paxos
// Thread safe
// Returns isCommit/isAccept, isFastPathSuccessful, execVal (if applicable)
func (c *Client) fpPropose(op *rpc.Operation) (bool, bool, string) {
	req := &rpc.ProposeRequest{
		Op:       op,
		ClientId: c.Id,
	}

	fastC := make(chan *rpc.ProposeReply, c.replicaNum)
	slowC := make(chan *rpc.ProposeReply, 1)

	/*
		//// Single RPC
		// Sends request to the fast-paxos shard leader
		c.sendProposeReqToLeader(c.leaderAddr, req, fastC, slowC, &c.wg)
		// Sends request to shard followers
		c.bcstProposeReqToFollowers(c.followerAddrList, req, fastC, &c.wg)
	*/
	//// Stream RPC
	c.streamSendProposeReqToLeader(c.leaderAddr, req, fastC, slowC, &c.wg)
	c.streamBcstProposeReqToFollowers(c.followerAddrList, req, fastC, &c.wg)

	// return results
	isWait, isCommit, isFast, ret := true, false, false, ""
	fpRetTable := make(map[string]int) // idx --> count
	for isWait {
		select {
		case reply := <-fastC:
			if _, ok := fpRetTable[reply.Idx]; !ok {
				fpRetTable[reply.Idx] = 0
			}
			fpRetTable[reply.Idx]++
			if fpRetTable[reply.Idx] == c.fastQuorum {
				isCommit = true
				isFast = true
				if !c.waitExecRet {
					isWait = false
					continue
				}
			}
		case <-slowC:
			isCommit = true
			// NOTE: for now Fast Paxos only has slow path reply when the fast path fails. Does not change isFast value.
			// TODO If this is to wait execRet while the the fast path commits, has
			// to set isFast.  However, the slow path exec resulst may come back
			// earlier than the fast path commit.  Therefore, we have to put isFast
			// value in the slow-path reply too.
			//ret = reply.ExecRet        // TODO implements Fast Paxos to return execution result
			isWait = false
		}
	}
	//logger.Infof("cmdId = %s, isCommit = %t, isFast = %t, ret = %s", cmd.Id, isCommit, isFast, ret)
	return isCommit, isFast, ret
}

/////// Single RPC I/O /////
// Non-blocking
func (c *Client) sendProposeReqToLeader(
	leaderAddr string, req *rpc.ProposeRequest, fastC, slowC chan<- *rpc.ProposeReply, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fastReply, stream := c.sendProposeReq(leaderAddr, req)
		fastC <- fastReply

		slowReply, err := stream.Recv()

		if err != nil {
			// Does nothing if the leader does not return slow-path result, where err == io.EOF
			if err != io.EOF {
				logger.Fatalf("Fails to recv slow-path reply from addr = %s req = %v error %v",
					leaderAddr, req, err)
			}
		} else {
			slowC <- slowReply
		}

		wg.Done()
	}()
}

func (c *Client) bcstProposeReqToFollowers(
	followerAddrList []string, req *rpc.ProposeRequest, fastC chan<- *rpc.ProposeReply, wg *sync.WaitGroup) {
	for _, addr := range followerAddrList {
		wg.Add(1)
		go func(addr string) {
			reply, stream := c.sendProposeReq(addr, req)
			fastC <- reply
			_, err := stream.Recv()
			if err != io.EOF {
				logger.Fatalf("Follower addr = %s replies twice. operation id = %s error = %v",
					addr, req.Op.Id, err)
			}
			wg.Done()
		}(addr)
	}
}

func (c *Client) sendProposeReq(
	addr string,
	req *rpc.ProposeRequest,
) (*rpc.ProposeReply, rpc.FastPaxosRpc_ProposeClient) {
	rpcStub := c.GetRpcStub(addr)
	stream, err := rpcStub.Propose(context.Background(), req, grpc.WaitForReady(true))
	if err != nil {
		logger.Fatalf("Fails to start an RPC server addr = %s operation id = %s error = %v",
			addr, req.Op.Id, err)
	}

	reply, err := stream.Recv() // fast-path reply
	if err != nil {
		logger.Fatalf("Fails to receive the fast-path reply from server addr = %s operation id = %s error = %v",
			addr, req.Op.Id, err)
	}

	return reply, stream
}

////////////////////////////////////////
// Helper functions
// Thread-safe
func (c *Client) genOpId() string {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.opCount++

	return c.Id + "-" + strconv.Itoa(c.opCount)
}

func (c *Client) GetRpcStub(addr string) rpc.FastPaxosRpcClient {
	conn, exists := c.cm.GetConnection(addr)
	if !exists {
		conn = c.cm.BuildConnection(addr)
	}

	rpcStub := rpc.NewFastPaxosRpcClient(conn)

	return rpcStub
}
