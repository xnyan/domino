package clientlib

import (
	"golang.org/x/net/context"
	"domino/fastpaxos/rpc"
)

func (c *Client) DebugServerStat() {
	req := &rpc.TestRequest{}
	rpcStub := c.GetRpcStub(c.leaderAddr)
	rpcStub.Test(context.Background(), req)

	for _, addr := range c.followerAddrList {
		rpcStub = c.GetRpcStub(addr)
		rpcStub.Test(context.Background(), req)
	}
}
