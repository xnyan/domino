package main

import (
	"bufio"
	"net"
	"sync"

	"domino/common"
	"domino/epaxos/genericsmrproto"
	"domino/epaxos/state"
	//"domino/hybrid/config"
	"domino/dynamic/config"
)

type EpaxosClient struct {
	noLeader bool
	fast     bool

	leaderAddr       string
	followerAddrList []string
	targetServerAddr string
	replyAddr        string // the server that sends back the reply

	serverConnTable map[string]net.Conn      // tcp connections
	readerTable     map[string]*bufio.Reader // receiving network messages
	readerLockTable map[string]*sync.Mutex
	writerTable     map[string]*bufio.Writer // sending network messages
	writerLockTable map[string]*sync.Mutex
	txnCount        int32
	txnCountLock    sync.Mutex

	cmdTable     map[int32]*common.Future
	cmdTableLock sync.RWMutex
}

func NewEpaxosClient(dcId, targetDcId, configFile, replicaFile string, noLeader bool, fast bool) *EpaxosClient {
	// Re-uses Dynamic config and assumes that there is only one Fast Paxos shard leader
	// Uses the Fast Paxos shard leader to be the Multi-Paxos leader
	var leaderInfo, targetInfo *config.ReplicaInfo
	followerList := make([]*config.ReplicaInfo, 0)
	repDir := config.LoadReplicaInfo(replicaFile)
	for _, rInfo := range repDir {
		if rInfo.IsFpLeader {
			leaderInfo = rInfo
		} else {
			followerList = append(followerList, rInfo)
		}
		if rInfo.DcId == targetDcId {
			// The replica's DcId that this client will send requests to
			targetInfo = rInfo
		}
	}
	if targetInfo == nil {
		logger.Fatalf("Missing a valid targetDcId = %s", targetDcId)
	}

	p := common.NewProperties()
	p.Load(configFile)

	c := &EpaxosClient{
		noLeader:         noLeader,
		fast:             fast,
		leaderAddr:       leaderInfo.GetNetAddr(),
		followerAddrList: make([]string, 0),
		serverConnTable:  make(map[string]net.Conn),
		readerTable:      make(map[string]*bufio.Reader),
		readerLockTable:  make(map[string]*sync.Mutex),
		writerTable:      make(map[string]*bufio.Writer),
		writerLockTable:  make(map[string]*sync.Mutex),
		txnCount:         0,
		cmdTable:         make(map[int32]*common.Future),
	}

	for _, f := range followerList {
		c.followerAddrList = append(c.followerAddrList, f.GetNetAddr())
	}

	/*
		// chooses the closest server as the operation coordinator
		latFile := p.GetStr(config.FLAG_NETWORK_LAT_FILE)
		latTag := p.GetStr(config.FLAG_NETWORK_LAT_TAG)
		netManager := common.NewStaticNetworkManager(latFile, latTag)
		dcList := p.GetStrList(config.FLAG_DYNAMIC_DC_LIST, ",")
		closestDcL := netManager.GetClosestQuorum(dcId, dcList, 1)
		closestDc := closestDcL[0]

		// TODO clean up this code for lowest lat majority calculation
		isUseLowestLatMajority := p.GetBoolWithDefault(config.FLAG_PAXOS_MAJORITY_LOWEST_LAT, "false")
		if isUseLowestLatMajority {
			dcReplicaMap, replicaDcMap := make(map[string]string), make(map[string]string)
			replicaList := make([]string, 0)
			for i, r := range rList {
				addr := p.GetStr(r)
				replicaDcMap[addr] = dcList[i]
			}
			for addr, dc := range replicaDcMap {
				dcReplicaMap[dc] = addr
				replicaList = append(replicaList, addr)
			}
			_, pShardLeaderMap, _ := config.GetShardAddrInfo(p)
			pShardLeaderDcList := make([]string, 0)
			for _, pLeader := range pShardLeaderMap {
				dc := getDcId(pLeader, replicaDcMap)
				pShardLeaderDcList = append(pShardLeaderDcList, dc)
			}

			majority := len(replicaList)/2 + 1
			minLat := calPaxosLat(dcId, closestDc,
				netManager, dcReplicaMap, replicaDcMap, replicaList, majority)
			for _, dc := range pShardLeaderDcList {
				lat := calPaxosLat(dcId, dc,
					netManager, dcReplicaMap, replicaDcMap, replicaList, majority)
				if lat < minLat {
					minLat, closestDc = lat, dc
				}
			}
		}

		for i, dc := range dcList {
			if dc == closestDc {
				c.targetServerAddr = p.GetStr(rList[i])
			}
		}
	*/

	c.targetServerAddr = targetInfo.GetNetAddr()
	if c.noLeader {
		c.replyAddr = c.targetServerAddr // epaxos, mencius
	} else {
		c.replyAddr = c.leaderAddr // multi-paxos, generalized paxos
	}

	c.Init()

	return c
}

func calPaxosLat(
	cDcId, leaderDcId string,
	netManager common.NetworkManager,
	dcReplicaMap, replicaDcMap map[string]string, replicaAddrList []string,
	majority int,
) int64 {
	lat := netManager.GetOneWayNetDelay(cDcId, leaderDcId) * 2
	leader := getAddr(leaderDcId, dcReplicaMap)
	fAddrList := calFollowerAddrList(leader, replicaAddrList)
	fDcList := getDcList(fAddrList, replicaDcMap)
	mDcList := netManager.GetClosestQuorum(leaderDcId, fDcList, majority-1)
	lat += netManager.MaxOneWayNetDelay(leaderDcId, mDcList) * 2
	return int64(lat)
}

// Returns the dcId for a given network address
func getDcId(addr string, replicaDcMap map[string]string) string {
	dc, ok := replicaDcMap[addr]
	if !ok {
		logger.Fatalf("No dcId for addr = %s", addr)
	}
	return dc
}

// Returns the replica address in the given dc
func getAddr(dcId string, dcReplicaMap map[string]string) string {
	addr, ok := dcReplicaMap[dcId]
	if !ok {
		logger.Fatalf("No addr for dcId = %s", dcId)
	}
	return addr
}

// Returns a list of followers for the given shard
func calFollowerAddrList(leader string, replicaList []string) []string {
	addrList := make([]string, 0)
	for _, addr := range replicaList {
		if addr == leader {
			continue
		}
		addrList = append(addrList, addr)
	}
	return addrList
}

func getDcList(addrList []string, replicaDcMap map[string]string) []string {
	dcList := make([]string, len(addrList), len(addrList))
	for i, addr := range addrList {
		dcList[i] = getDcId(addr, replicaDcMap)
	}
	return dcList
}

// thread-safe
func (c *EpaxosClient) GetTxnId() int32 {
	c.txnCountLock.Lock()
	c.txnCount++
	c.txnCountLock.Unlock()
	return c.txnCount
}

type EpaxosRet struct {
	IsAccept bool
	Val      string
	IsFast	bool
}

func (c *EpaxosClient) SyncExecTxn(
	rKeyList []string, wTable map[string]string,
) (bool, bool, bool, string) {
	args := c.buildArgs(rKeyList, wTable)
	if args == nil {
		logger.Errorf("Empty read/write sets.")
		return false, false, false, ""
	}
	done := c.WaitCmd(args.CommandId)
	//logger.Infof("Sending cmdId = %d, op = %d key = %s val = %s", args.CommandId, args.Command.Op, args.Command.K, args.Command.V)
	c.Propose(args)
	ret := done.GetValue().(*EpaxosRet)
	//logger.Infof("Received cmdId = %d isAccept = %t", args.CommandId, ret.IsAccept)
	c.DelCmd(args.CommandId)
	return false, ret.IsAccept, ret.IsFast, ret.Val
}

func (c *EpaxosClient) WaitCmd(cmdId int32) *common.Future {
	c.cmdTableLock.Lock()
	defer c.cmdTableLock.Unlock()
	if _, ok := c.cmdTable[cmdId]; !ok {
		c.cmdTable[cmdId] = common.NewFuture()
	}
	return c.cmdTable[cmdId]
}

func (c *EpaxosClient) DelCmd(cmdId int32) {
	c.cmdTableLock.Lock()
	defer c.cmdTableLock.Unlock()
	delete(c.cmdTable, cmdId)
}

func (c *EpaxosClient) HandleReplies() {
	reader := c.GetReader(c.replyAddr)
	for {
		reply := new(genericsmrproto.ProposeReplyTS)
		if err := reply.Unmarshal(reader); err != nil {
			logger.Fatal("Error when reading:", err)
		}
		c.cmdTableLock.RLock()
		done, ok := c.cmdTable[reply.CommandId]
		c.cmdTableLock.RUnlock()
		if !ok {
			logger.Fatalf("No thread is waiting for cmdId = %d", reply.CommandId)
		}
		if reply.OK != 0 {
			if reply.Slowpath != 0 {
				done.SetValue(&EpaxosRet{IsAccept: true, Val: string(reply.Value), IsFast: true})
			} else {
				done.SetValue(&EpaxosRet{IsAccept: true, Val: string(reply.Value), IsFast: false})
			}
		} else {
			if reply.Slowpath != 0 {
				done.SetValue(&EpaxosRet{IsAccept: false, Val: string(reply.Value), IsFast: true })
			} else {
				done.SetValue(&EpaxosRet{IsAccept: false, Val: string(reply.Value), IsFast: false})
			}
		}
	}
}

func (c *EpaxosClient) buildArgs(
	rKeyList []string, wTable map[string]string,
) *genericsmrproto.Propose {
	args := &genericsmrproto.Propose{c.GetTxnId(), state.Command{state.PUT, "", ""}, 0}

	for _, rk := range rKeyList {
		args.Command.Op = state.GET
		args.Command.K = state.Key(rk)
		args.Command.V = state.Value("")
		return args
	}

	for wk, wv := range wTable {
		args.Command.Op = state.PUT
		args.Command.K = state.Key(wk)
		args.Command.V = state.Value(wv)
		return args
	}

	return nil
}

func (c *EpaxosClient) ExecTxn(
	rKeyList []string, wTable map[string]string,
) (bool, bool, bool, string) {
	if args := c.buildArgs(rKeyList, wTable); args == nil {
		logger.Errorf("Empty read/write sets.")
		return false, false, false, ""
	} else {
		c.Propose(args)
		_, ok, val , isFast := c.WaitReply(c.replyAddr)
		return false, ok, isFast, val
	}
}

func (c *EpaxosClient) Close() {
	for _, conn := range c.serverConnTable {
		if conn != nil {
			conn.Close()
		}
	}
}

//Returns the server address that will send the reply
//thread-safe
func (c *EpaxosClient) Propose(args *genericsmrproto.Propose) string {
	if c.noLeader {
		// epaxos, mencius
		//logger.Infof("Sending proposal cmdId = %d to target addr = %s", args.CommandId, c.targetServerAddr)
		c.SendProposal(c.targetServerAddr, args)
		return c.targetServerAddr
	} else {
		// multi-paxos
		//logger.Infof("Sending proposal cmdId = %d to leader addr = %s", args.CommandId, c.targetServerAddr)
		c.SendProposal(c.leaderAddr, args)
		if c.fast {
			// generalized paxos, sends the proposal to everybody
			for _, addr := range c.followerAddrList {
				c.SendProposal(addr, args)
			}
		}
		return c.leaderAddr
	}
}

//Returns commondId, OK, value
func (c *EpaxosClient) WaitReply(addr string) (int32, bool, string, bool) {
	reader := c.GetReader(addr)

	l := c.GetReaderLock(addr)
	l.Lock()
	defer l.Unlock()

	reply := new(genericsmrproto.ProposeReplyTS)
	if err := reply.Unmarshal(reader); err != nil {
		logger.Fatal("Error when reading:", err)
	}
	if reply.OK != 0 {
		if reply.Slowpath != 0 {
			return reply.CommandId, true, string(reply.Value), true
		} else {
			return reply.CommandId, true, string(reply.Value), false
		}
	} else {
		if reply.Slowpath != 0 {
			return reply.CommandId, false, string(reply.Value), true
		} else {
			return reply.CommandId, false, string(reply.Value), false
		}
	}
}

//thread-safe
func (c *EpaxosClient) SendProposal(addr string, args *genericsmrproto.Propose) {
	writer := c.GetWriter(addr)

	l := c.GetWriterLock(addr)
	l.Lock()
	defer l.Unlock()

	writer.WriteByte(genericsmrproto.PROPOSE)
	args.Marshal(writer)
	writer.Flush()
}

func (c *EpaxosClient) GetWriterLock(addr string) *sync.Mutex {
	return c.writerLockTable[addr]
}

func (c *EpaxosClient) GetWriter(addr string) *bufio.Writer {
	return c.writerTable[addr]
}

func (c *EpaxosClient) GetReaderLock(addr string) *sync.Mutex {
	return c.readerLockTable[addr]
}
func (c *EpaxosClient) GetReader(addr string) *bufio.Reader {
	return c.readerTable[addr]
}

func (c *EpaxosClient) Init() {
	c.BuildConn(c.leaderAddr)
	c.InitConnLock(c.leaderAddr)
	for _, addr := range c.followerAddrList {
		c.BuildConn(addr)
		c.InitConnLock(addr)
	}
}

func (c *EpaxosClient) InitConnLock(addr string) {
	c.readerLockTable[addr] = &sync.Mutex{}
	c.writerLockTable[addr] = &sync.Mutex{}
}

// Builds connections to every server
func (c *EpaxosClient) BuildConn(addr string) {
	var err error
	c.serverConnTable[addr], err = net.Dial("tcp", addr)
	if err != nil {
		logger.Fatalf("Error connecting to replica %s", addr)
	}
	c.readerTable[addr] = bufio.NewReader(c.serverConnTable[addr])
	c.writerTable[addr] = bufio.NewWriter(c.serverConnTable[addr])
}
