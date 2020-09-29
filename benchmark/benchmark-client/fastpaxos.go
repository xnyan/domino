package main

import (
	"domino/fastpaxos/clientlib"
)

type FastPaxosClient struct {
	lib *clientlib.Client
}

func NewFastPaxosClient(id, dcId, configFile, replicaFile string) *FastPaxosClient {
	c := &FastPaxosClient{
		lib: clientlib.NewClient(id, dcId, configFile, replicaFile),
	}
	return c
}

func (c *FastPaxosClient) ExecTxn(
	rKeyList []string, wTable map[string]string,
) (bool, bool, bool, string) {
	op := c.lib.BuildOperation(rKeyList, wTable)
	isCommitted, isFast, ret := c.lib.Propose(op)
	//logger.Infof("op = %s execRet = %s", op.Id, ret)
	return true, isCommitted, isFast, ret
}

func (c *FastPaxosClient) SyncExecTxn(
	rKeyList []string, wTable map[string]string,
) (bool, bool, bool, string) {
	return c.ExecTxn(rKeyList, wTable)
}

func (c *FastPaxosClient) Close() {
	c.lib.Close()
}
