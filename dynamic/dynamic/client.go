package dynamic

import (
	"sync"
)

type Client interface {
	// Returns isCommitted, the accepted time slot, execRet
	PaxosPropose(c *Command, leader string) (bool, string)

	// Returns isCommitted, isFast, the accepted time slot, execRet
	FastPaxosPropose(
		c *Command, shard int32, leader string, followerList []string, t int64,
	) (bool, bool, string)

	// Fast Paxos all leaners, for reducing execution latency
	ExecFastPaxosPropose(
		c *Command, shard int32, leader, execReplica string, replicaList []string, t int64,
	) (bool, bool, string)

	// Latency monitoring
	// Returns the queuing delay in ns and Paxos latency in ms
	Probe(addr string) (int64, int32)
	ProbeTime(addr string) (int64, int32)

	InitConn(addrList []string)

	// Shuts down
	Close()

	// Testing
	Test(addrList []string)
}

type DynamicClient struct {
	replicaNum int // The total number of replicas
	fastQuorum int

	// For Fast Paxos
	waitExecRet bool
	wg          sync.WaitGroup

	// If the Fast Paxos leader uses Paxos to accept commonds that arerejected in
	// Fast Paxos instances
	IsFpLeaderUsePaxos bool

	// Network I/O
	io ClientIo
}

func NewDynamicClient(
	replicaNum, fastQuorum int,
	isWaitExec, isFpLeaderUsePaxos, isGrpc bool,
) Client {
	c := &DynamicClient{
		replicaNum:         replicaNum,
		fastQuorum:         fastQuorum,
		waitExecRet:        isWaitExec,
		IsFpLeaderUsePaxos: isFpLeaderUsePaxos,
	}

	if isGrpc {
		c.io = NewClientGrpc()
	} else {
		// TODO supports network latency monitoring
		c.io = NewClientFastIo(&c.wg)
	}

	return c
}

func (c *DynamicClient) InitConn(addrList []string) {
	c.io.InitConn(addrList)
}

//Returns server-side queuing delay in ns
func (c *DynamicClient) Probe(addr string) (int64, int32) {
	reply := c.io.SendProbeReq(addr)
	return reply.QueuingDelay, reply.PaxosLat
}

func (c *DynamicClient) ProbeTime(addr string) (int64, int32) {
	reply := c.io.SendProbeTimeReq(addr)
	return reply.ProcessTime, reply.PaxosLat
}

// Proposes a command to the given Paxos instance
func (c *DynamicClient) PaxosPropose(cmd *Command, addr string) (bool, string) {
	reply := c.io.SendPaxosProposeReq(addr, cmd)

	logger.Debugf("%s reply from addr = %s %t %s", cmd.Id, addr, reply.IsCommit, reply.ExecRet)

	return reply.IsCommit, reply.ExecRet
}

// Proposes a command to the given Fast Paxos instance
func (c *DynamicClient) FastPaxosPropose(
	cmd *Command, shard int32, leader string, followerList []string, t int64,
) (bool, bool, string) {
	fastC := make(chan *FpProposeReply, c.replicaNum)
	slowC := make(chan *FpProposeReply, 1)
	req := &FpProposeReq{
		Cmd:  cmd,
		Time: &Timestamp{Time: t, Shard: shard},
	}

	// Sends request to the fast-paxos shard leader
	c.io.SendFpProposeToLeader(leader, req, fastC, slowC, &c.wg)

	// Sends request to shard followers
	c.io.BcstFpProposeToFollowers(followerList, req, fastC, &c.wg)

	// return results
	isWait, isCommit, isFast, ret := true, false, false, ""
	a, r := 0, 0 // fast-path accept / reject votes
	for isWait {
		select {
		case reply := <-fastC:
			if reply.IsAccept {
				a++
			} else {
				r++
			}
			if a == c.fastQuorum {
				isCommit = true
				isFast = true
				if !c.waitExecRet {
					isWait = false
					continue
				}
			}
			if r == c.replicaNum {
				if !c.IsFpLeaderUsePaxos {
					isWait = false // all rejects
				}
			}
		case reply := <-slowC:
			isCommit = reply.IsAccept
			isFast = reply.IsFast // Note: the slow-path result may come back earlier than the fast-path result
			ret = reply.ExecRet
			isWait = false
		}
	}

	//logger.Infof("cmdId = %s, isCommit = %t, isFast = %t, ret = %s", cmd.Id, isCommit, isFast, ret)

	return isCommit, isFast, ret
}

func (c *DynamicClient) ExecFastPaxosPropose(
	cmd *Command, shard int32, leader, execReplica string, replicaList []string, t int64,
) (bool, bool, string) {
	fastC := make(chan *FpProposeReply, c.replicaNum)
	slowC := make(chan *FpProposeReply, 2)

	execReq := &FpProposeReq{
		Cmd:         cmd,
		Time:        &Timestamp{Time: t, Shard: shard},
		IsExecReply: true,
	}
	nonExecReq := &FpProposeReq{
		Cmd:         cmd,
		Time:        &Timestamp{Time: t, Shard: shard},
		IsExecReply: false,
	}

	/*
		// Only the specific exec Replica and the Leader will return execution
		// results if the fast path succeeds.
		// Sends request to the fast-paxos shard leader and the chosen execution replica
		c.io.SendFpProposeToExecReplica(execReplica, execReq, fastC, slowC, &c.wg)
		if leader != execReplica {
			c.io.SendFpProposeToExecReplica(leader, nonExecReq, fastC, slowC, &c.wg)
			rList := make([]string, 0, len(replicaList)-1) // TODO avoids creaing such a list everytime
			for _, r := range replicaList {
				if r != leader {
					rList = append(rList, r)
				}
			}
			// Sends request to other replicas
			c.io.BcstFpProposeToFollowers(rList, nonExecReq, fastC, &c.wg)
		} else {
			// Sends request to other replicas
			c.io.BcstFpProposeToFollowers(replicaList, nonExecReq, fastC, &c.wg)
		}
	*/
	// All replicas can return execution results if the fast path succeeds. The
	// client will take the latest execution result.
	c.io.SendFpProposeToExecReplica(execReplica, execReq, fastC, slowC, &c.wg)
	for _, r := range replicaList {
		c.io.SendFpProposeToExecReplica(r, nonExecReq, fastC, slowC, &c.wg)
	}

	// return results
	isWait, isCommit, isFast, ret := true, false, false, ""
	a, r := 0, 0 // fast-path accept / reject votes
	for isWait {
		select {
		case reply := <-fastC:
			if reply.IsAccept {
				a++
			} else {
				r++
			}
			if a == c.fastQuorum {
				isCommit = true
				isFast = true
				if !c.waitExecRet {
					isWait = false
					continue
				}
			}
			if r == c.replicaNum {
				if !c.IsFpLeaderUsePaxos {
					isWait = false // all rejects
				}
			}
		case reply := <-slowC:
			isCommit = reply.IsAccept
			isFast = reply.IsFast // Note: the slow-path result may come back earlier than the fast-path result
			ret = reply.ExecRet
			isWait = false
		}
	}

	//logger.Infof("Result cmdId = %s, isCommit = %t, isFast = %t, ret = %s", cmd.Id, isCommit, isFast, ret)

	return isCommit, isFast, ret
}

func (c *DynamicClient) Close() {
	// Waits for all threads to complete
	c.wg.Wait()
}

// Testing
func (c *DynamicClient) Test(addrList []string) {
	for _, addr := range addrList {
		c.io.SendTestReq(addr)
	}
}
