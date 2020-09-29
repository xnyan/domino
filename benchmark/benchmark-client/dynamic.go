package main

import (
	"domino/dynamic/clientlib"
	"domino/dynamic/dynamic"
)

type DynamicClient struct {
	lib clientlib.ClientLib
}

func NewDynamicClient(id, dcId, configFile, replicaFile, targetDcId string) *DynamicClient {
	c := &DynamicClient{
		lib: clientlib.NewClientLib(id, dcId, configFile, replicaFile, targetDcId),
	}
	return c
}

func (c *DynamicClient) ExecTxn(
	rKeyList []string, wTable map[string]string,
) (bool, bool, bool, string) {
	cmd := buildDynamicCommand(rKeyList, wTable)
	isUseFp, isCommitted, isFast, ret := c.lib.Propose(cmd)
	//logger.Infof("cmd = %s execRet = %s", cmd.Id, ret)
	return isUseFp, isCommitted, isFast, ret
}

func (c *DynamicClient) SyncExecTxn(
	rKeyList []string, wTable map[string]string,
) (bool, bool, bool, string) {
	return c.ExecTxn(rKeyList, wTable)
}

func (c *DynamicClient) Close() {
	c.lib.Close()
}

//Helper function
func buildDynamicCommand(rKeyList []string, wTable map[string]string) *dynamic.Command {
	for _, k := range rKeyList {
		return &dynamic.Command{Type: "r", Key: k}
	}

	for k, v := range wTable {
		return &dynamic.Command{Type: "w", Key: k, Val: v}
	}

	return nil
}
