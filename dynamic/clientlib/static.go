package clientlib

/*
import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"domino/common"
	"domino/dynamic/config"
	"domino/dynamic/dynamic"
)

//// Static Network Prediction ////
// TODO: supports using the max latency to the closest supermajority to set the timestamp for Fast Paxos
type StaticClientLib struct {
	c dynamic.Client

	Id   string // client id
	DcId string

	// Command counting
	count     int64
	countLock sync.Mutex

	// Replica information
	replicaNum   int
	majority     int
	fastQuorum   int
	replicaDcMap map[string]string // replica address (ip:port) --> dcId
	dcReplicaMap map[string]string // dcId --> replica address (ip:port)
	dcList       []string          // a list of all replica DC ids
	replicaList  []string          // a list of all replica addresses (ip:port)

	isFastPaxos bool // ture to use fast paxos

	// Multi Paxos
	paxosShardLeader       string
	paxosShardDelay        int64
	pShardLeaderDcList     []string // a list of Paxos shard leader DCs
	isUseLowestLatMajority bool     // false: uses the closest leader; true: uses the one can give lowest latency

	// Fast Paxos
	closestFpShard             int32
	closestFpShardLeader       string
	closestFpShardFollowerList []string
	fpTimeDelay                int64
	fpShardLeaderDcList        []string
	fpLeaderShardMap           map[string]int32
	// Fast Paxos Optimization for execution latency
	isFpLeaderLearner        bool
	fpExecReplicaAddr        string // may be or may not be the closes Fast Paxos shard leader
	fpNonExecReplicaAddrList []string

	// Clock information
	clock common.Clock

	// Network information
	netManager common.NetworkManager
	addDelay   int64 //additional delay for each request, unit: ns

	// Random
	rnd     *rand.Rand
	fpLoadP int
}

func newStaticClientLib(id, dcId, configFile string) *StaticClientLib {
	// Parses configurations
	p := common.NewProperties()
	p.Load(configFile)

	replicaDcMap := make(map[string]string)
	rList := p.GetStrList(config.FLAG_DYNAMIC_REPLICA_LIST, ",")
	dcList := p.GetStrList(config.FLAG_DYNAMIC_DC_LIST, ",")
	for i, r := range rList {
		addr := p.GetStr(r)
		replicaDcMap[addr] = dcList[i]
	}

	_, pShardLeaderMap, fpShardLeaderMap := config.GetShardAddrInfo(p)

	latFile := p.GetStr(config.FLAG_NETWORK_LAT_FILE)
	latTag := p.GetStr(config.FLAG_NETWORK_LAT_TAG)
	addDelay := p.GetTimeDuration(config.FLAG_DYNAMIC_CLIENT_ADD_DELAY)
	isExecReply := p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC_REPLY, "false")
	isFpLeaderUsePaxos := p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_LEADER_USE_PAXOS, "true")
	isGrpc := p.GetBoolWithDefault(config.FLAG_DYNAMIC_GRPC, "true")

	// Construct the lib instance
	lib := &StaticClientLib{
		Id:                  id,
		DcId:                dcId,
		count:               0,
		replicaNum:          len(replicaDcMap),
		replicaDcMap:        replicaDcMap,
		dcReplicaMap:        make(map[string]string),
		dcList:              make([]string, 0, len(replicaDcMap)),
		replicaList:         make([]string, 0, len(replicaDcMap)),
		pShardLeaderDcList:  make([]string, 0, len(pShardLeaderMap)),
		fpShardLeaderDcList: make([]string, 0, len(fpShardLeaderMap)),
		fpLeaderShardMap:    make(map[string]int32),
		clock:               common.NewSysNanoClock(),
		netManager:          common.NewStaticNetworkManager(latFile, latTag),
		addDelay:            int64(addDelay),
	}

	lib.isUseLowestLatMajority = p.GetBoolWithDefault(config.FLAG_PAXOS_MAJORITY_LOWEST_LAT, "false")
	lib.isFpLeaderLearner = p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_LEADER_LEARNER, "true")

	f := (lib.replicaNum - 1) / 2
	lib.majority = f + 1
	lib.fastQuorum = int(math.Ceil((3.0*float64(f))/2.0)) + 1

	lib.c = dynamic.NewDynamicClient(
		lib.replicaNum, lib.fastQuorum, isExecReply, isFpLeaderUsePaxos, isGrpc)

	for addr, dc := range lib.replicaDcMap {
		lib.dcReplicaMap[dc] = addr
		lib.dcList = append(lib.dcList, dc)
		lib.replicaList = append(lib.replicaList, addr)
	}

	for _, pLeader := range pShardLeaderMap {
		dc := lib.getDcId(pLeader)
		lib.pShardLeaderDcList = append(lib.pShardLeaderDcList, dc)
	}

	for fpShard, fpLeader := range fpShardLeaderMap {
		dc := lib.getDcId(fpLeader)
		lib.fpShardLeaderDcList = append(lib.fpShardLeaderDcList, dc)
		lib.fpLeaderShardMap[fpLeader] = int32(fpShard)
	}

	// Statically sets Paxos and Fast Paxos shards
	//isP := lib.setClosestPaxosShard()
	isP := lib.setPaxosShard()
	isFp := lib.setClosestFpShard()

	if !isP && !isFp {
		logger.Fatalf("No available Paxos or Fast Paxos shards. Check configurations!")
	}

	// Static way of choosing Paxos or Fast Paxos
	if isP && isFp {
		lib.isFastPaxos = lib.isToUseFp()
	} else if isFp {
		lib.isFastPaxos = true
	} else { // isP == true
		lib.isFastPaxos = false
	}

	//logger.Infof("Static isFastPaxos = %t", lib.isFastPaxos)

	lib.c.InitConn(lib.replicaList)

	lib.fpLoadP = p.GetIntWithDefault(config.FLAG_DYNAMIC_CLIENT_FP_LOAD, "-1")
	if lib.fpLoadP >= 0 {
		lib.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	fmt.Println("#Dynamic client static isFastPaxos =", lib.isFastPaxos, "fpLoad =", lib.fpLoadP)

	return lib
}

func (lib *StaticClientLib) rndFastPaxos() bool {
	if lib.rnd.Intn(100) < lib.fpLoadP {
		return true
	}
	return false
}

func (lib *StaticClientLib) Propose(cmd *dynamic.Command) (bool, bool, string) {
	cmd.Id = lib.genCmdId()

	// FP vs Paxos load configuration
	if lib.fpLoadP >= 0 {
		if lib.rndFastPaxos() {
			return lib.fastPaxos(cmd)
		}
		ok, ret := lib.paxos(cmd)
		return ok, false, ret
	}

	// Static Network Configuration
	if lib.isFastPaxos {
		return lib.fastPaxos(cmd)
	}

	ok, ret := lib.paxos(cmd)
	return ok, false, ret
}

func (lib *StaticClientLib) FpPropose(cmd *dynamic.Command) (bool, bool, string) {
	return lib.fastPaxos(cmd)
}

func (lib *StaticClientLib) PaxosPropose(cmd *dynamic.Command) (bool, string) {
	return lib.paxos(cmd)
}

func (lib *StaticClientLib) paxos(
	cmd *dynamic.Command,
) (bool, string) {
	addr := lib.getPaxosShardLeader()
	return lib.c.PaxosPropose(cmd, addr)
}

func (lib *StaticClientLib) fastPaxos(
	cmd *dynamic.Command,
) (bool, bool, string) {
	if lib.isFpLeaderLearner {
		shard, leader := lib.getClosestFpShardLeader()
		followerList := lib.getClosestFpShardFollower()
		t := lib.getFpTimestamp()

		//now := time.Now().UnixNano()
		//logger.Infof("cmdId = %s, curTime = %d, timestamp = %d, future = %d",
		//	cmd.Id, now, t, t-now)

		return lib.c.FastPaxosPropose(cmd, shard, leader, followerList, t)
	} else {
		shard, leader := lib.getClosestFpShardLeader()
		fpExecReplica := lib.getFpExecReplicaAddr()
		fpNonExecReplicaList := lib.getFpNonExecReplicaAddrList()
		t := lib.getFpTimestamp()

		//logger.Infof("cmdId = %s, curTime = %d, timestamp = %d, leader = %s, execReplica = %s, other = %v",
		//	cmd.Id, time.Now().UnixNano(), t, leader, fpExecReplica, fpNonExecReplicaList)

		return lib.c.ExecFastPaxosPropose(cmd, shard, leader, fpExecReplica, fpNonExecReplicaList, t)
	}
}

func (lib *StaticClientLib) Close() {
	lib.c.Close()
}

func (lib *StaticClientLib) Test() {
	lib.c.Test(lib.replicaList)
}

///// Helper Functions /////

// Thread-safe
// Returns a unique id for this client lib
func (lib *StaticClientLib) genCmdId() string {
	lib.countLock.Lock()
	id := fmt.Sprintf("%s-%d", lib.Id, lib.count)
	lib.count++
	lib.countLock.Unlock()
	return id
}

// Returns the dcId for a given network address
func (lib *StaticClientLib) getDcId(addr string) string {
	dc, ok := lib.replicaDcMap[addr]
	if !ok {
		logger.Fatalf("No dcId for addr = %s", addr)
	}
	return dc
}

// Returns the replica address in the given dc
func (lib *StaticClientLib) getAddr(dcId string) string {
	addr, ok := lib.dcReplicaMap[dcId]
	if !ok {
		logger.Fatalf("No addr for dcId = %s", dcId)
	}
	return addr
}

// Returns the fast-paxos shard that has the leader be the given address
func (lib *StaticClientLib) getFpShardByAddr(addr string) int32 {
	shard, ok := lib.fpLeaderShardMap[addr]
	if !ok {
		logger.Fatalf("No fast-paxos shard for leader addr = %s", addr)
	}
	return shard
}

// Returns a list of followers for the given shard
func (lib *StaticClientLib) calFollowerAddrList(leader string) []string {
	addrList := make([]string, 0, len(lib.replicaList)-1)
	for _, addr := range lib.replicaList {
		if addr == leader {
			continue
		}
		addrList = append(addrList, addr)
	}
	return addrList
}

func (lib *StaticClientLib) calFpNonExecReplicaAddrList(leader, execReplica string) []string {
	addrList := make([]string, 0, len(lib.replicaList)-1)
	for _, addr := range lib.replicaList {
		if addr == leader || addr == execReplica {
			continue
		}
		addrList = append(addrList, addr)
	}
	return addrList
}

func (lib *StaticClientLib) getDcList(addrList []string) []string {
	dcList := make([]string, len(addrList), len(addrList))
	for i, addr := range addrList {
		dcList[i] = lib.getDcId(addr)
	}
	return dcList
}

// Selects to use Paxos or Fast Paxos
func (lib *StaticClientLib) isToUseFp() bool {
	if lib.calFpLat() < lib.calPaxosLat(lib.getDcId(lib.paxosShardLeader)) {
		return true
	}
	return false
}

func (lib *StaticClientLib) calFpLat() int64 {
	fpDcList := lib.netManager.GetClosestQuorum(lib.DcId, lib.dcList, lib.fastQuorum)
	lat := lib.netManager.MaxOneWayNetDelay(lib.DcId, fpDcList) * 2
	return int64(lat)
}

func (lib *StaticClientLib) setPaxosShard() bool {
	if len(lib.pShardLeaderDcList) == 0 {
		return false
	}

	lib.paxosShardLeader = lib.selectPaxosShardLeader()
	lib.paxosShardDelay = lib.calOneWayNetDelay(lib.paxosShardLeader) + lib.addDelay

	logger.Debugf("Closest Paxos shard leader = %s]", lib.paxosShardLeader)

	return true
}

func (lib *StaticClientLib) setClosestFpShard() bool {
	if len(lib.fpShardLeaderDcList) == 0 {
		return false
	}

	lib.closestFpShard, lib.closestFpShardLeader = lib.calClosestFpShardLeader()
	lib.closestFpShardFollowerList = lib.calFollowerAddrList(lib.closestFpShardLeader)
	lib.fpTimeDelay = lib.calMaxOneWayDelay(lib.dcList)
	lib.fpTimeDelay += lib.addDelay

	if !lib.isFpLeaderLearner {
		lib.fpExecReplicaAddr = lib.calFpExecReplicaAddr()
		lib.fpNonExecReplicaAddrList = lib.calFpNonExecReplicaAddrList(
			lib.closestFpShardLeader, lib.fpExecReplicaAddr)
	}

	return true
}

//////Paxos//////

// Returns the cached closest Paxos shard leader
func (lib *StaticClientLib) getPaxosShardLeader() string {
	return lib.paxosShardLeader
}

// Selecting the leader being the replica that achieves lowest latency by using mencius / multi-paxos
func (lib *StaticClientLib) selectPaxosShardLeader() string {
	dcList := lib.netManager.GetClosestQuorum(lib.DcId, lib.pShardLeaderDcList, 1)
	dcId := dcList[0] // closest dcId

	if lib.isUseLowestLatMajority {
		minLat := lib.calPaxosLat(dcId)
		for _, dc := range lib.pShardLeaderDcList {
			lat := lib.calPaxosLat(dc)
			if lat < minLat {
				minLat, dcId = lat, dc
			}
		}
	}

	leader := lib.getAddr(dcId)
	return leader
}

func (lib *StaticClientLib) calPaxosLat(leaderDcId string) int64 {
	lat := lib.netManager.GetOneWayNetDelay(lib.DcId, leaderDcId) * 2
	leader := lib.getAddr(leaderDcId)
	fAddrList := lib.calFollowerAddrList(leader)
	fDcList := lib.getDcList(fAddrList)
	mDcList := lib.netManager.GetClosestQuorum(leaderDcId, fDcList, lib.majority-1)
	lat += lib.netManager.MaxOneWayNetDelay(leaderDcId, mDcList) * 2
	return int64(lat)
}

// Returns the arrival time to the cached/selected Paxos shard leader.
func (lib *StaticClientLib) getPaxosShardTimestamp() int64 {
	t := lib.paxosShardDelay + lib.clock.GetClockTime()
	return t
}

//// Calculates and returns the closest Paxos shard leader
//func (lib *StaticClientLib) calClosestPaxosShardLeader() string {
//	dcList := lib.netManager.GetClosestQuorum(lib.DcId, lib.pShardLeaderDcList, 1)
//	dcId := dcList[0]
//	addr := lib.getAddr(dcId)
//	return addr
//}

/////Fast Paxos//////

// Returns shard Id and the shard leader's net addr (i.e., ip:port)
func (lib *StaticClientLib) getClosestFpShardLeader() (int32, string) {
	return lib.closestFpShard, lib.closestFpShardLeader
}

// Calculates and returns shard Id and the shard leader's net addr (i.e., ip:port)
func (lib *StaticClientLib) calClosestFpShardLeader() (int32, string) {
	dcList := lib.netManager.GetClosestQuorum(lib.DcId, lib.fpShardLeaderDcList, 1)
	dcId := dcList[0]
	addr := lib.getAddr(dcId)
	shard := lib.getFpShardByAddr(addr)
	return shard, addr
}

// Returns the cached follower list for the cahced Fast Paxos shard
func (lib *StaticClientLib) getClosestFpShardFollower() []string {
	return lib.closestFpShardFollowerList
}

// Returns the timestamp for Fast Paxos
func (lib *StaticClientLib) getFpTimestamp() int64 {
	t := lib.fpTimeDelay + lib.clock.GetClockTime()
	return t
}

// Generates the timestamp that the request is execpted to arrive at replicas in Fast Paxos
func (lib *StaticClientLib) calFpTimestamp() int64 {
	t := lib.predictMaxArrivalTimeByDc(lib.dcList)
	t += lib.addDelay
	return t
}

// Fast Paxos execution latency optimization
func (lib *StaticClientLib) calFpExecReplicaAddr() string {
	// Chooses the closest replica
	dcList := lib.netManager.GetClosestQuorum(lib.DcId, lib.dcList, 1)
	dcId := dcList[0] // closest dcId
	addr := lib.getAddr(dcId)
	return addr
}

func (lib *StaticClientLib) getFpExecReplicaAddr() string {
	return lib.fpExecReplicaAddr
}

func (lib *StaticClientLib) getFpNonExecReplicaAddrList() []string {
	return lib.fpNonExecReplicaAddrList
}

////Timestamp////

// Calculates the one-way delay to a replica
func (lib *StaticClientLib) calOneWayNetDelay(addr string) int64 {
	dcId := lib.getDcId(addr)
	delay := lib.netManager.GetOneWayNetDelay(lib.DcId, dcId)
	return int64(delay)
}

// Calculates the max one-way delay to a set of datacenters
func (lib *StaticClientLib) calMaxOneWayDelay(dcList []string) int64 {
	delay := lib.netManager.MaxOneWayNetDelay(lib.DcId, dcList)
	return int64(delay)
}

// Predicts the arrival time to a replica
func (lib *StaticClientLib) predictArrivalTimeByAddr(addr string) int64 {
	dcId := lib.getDcId(addr)
	return lib.predictArrivalTimeByDc(dcId)
}

func (lib *StaticClientLib) predictArrivalTimeByDc(dcId string) int64 {
	delay := lib.netManager.GetOneWayNetDelay(lib.DcId, dcId)
	t := lib.clock.GetClockTime() + int64(delay)
	return t
}

// Predicts the last arrival time to a set of replicas
func (lib *StaticClientLib) predictMaxArrivalTimeByAddr(addrList []string) int64 {
	dcList := lib.getDcList(addrList)
	return lib.predictMaxArrivalTimeByDc(dcList)
}

func (lib *StaticClientLib) predictMaxArrivalTimeByDc(dcList []string) int64 {
	delay := lib.netManager.MaxOneWayNetDelay(lib.DcId, dcList)
	t := lib.clock.GetClockTime() + int64(delay)
	return t
}
*/
