package clientlib

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"domino/common"
	"domino/dynamic/config"
	"domino/dynamic/dynamic"
	"domino/dynamic/latency"
)

type LatInfo struct {
	addr     string        // replica network address
	rt       time.Duration // roundtrip time including queuing delay
	qDelay   time.Duration // queuing delay (in ns) on the replica
	paxosLat time.Duration // additional paxos (in ms) latency when the replica is the leader
}

type LatTimeInfo struct {
	addr       string        // replica network address
	rt         time.Duration // roundtrip time including queuing delay
	timeOffset time.Duration // time offset between the clock time of sending (on client) and processing (on server)
	paxosLat   time.Duration // additional paxos (in ms) latency when the replica is the leader
}

type DynamicClientLib struct {
	c dynamic.Client

	Id                   string // client id
	IsReplicaColocated   bool   // true if colocated with a replica in the same datacenter
	ColocatedReplicaAddr string // the net addr of the colocated replica

	// Command counting
	count     int64
	countLock sync.Mutex

	// Clock information
	clock common.Clock

	// Replica information
	replicaNum      int
	majority        int
	fastQuorum      int
	replicaAddrList []string            // a list of all of the replicas' network addresses
	replicaAddrDir  map[string][]string // replica net addr --> other replicas' net addr
	//replicaIdAddrTable map[string]string   // replica ID --> replica net addr

	// Paxos
	isPaxosEnabled              bool
	isPaxosUseLowestLatMajority bool

	// Fast Paxos
	isAll                 bool // If the Fast Paxos uses the latency to the supermajority
	isFpLeaderSoleLearner bool // If the Fast Paxos leader is the sole learner
	fpShard               int32
	fpLeaderAddr          string

	// Network latency monitoring
	clientAddDelay       int64         // additional delay in ns
	probeInv             time.Duration // probing interval
	probeC               chan *LatInfo
	probeTimeC           chan *LatTimeInfo
	isProbeTime          bool
	fpOneWayLatPredictor *latency.LatPredictor // one-way message delay + queuing delay
	fpLatPredictor       *latency.LatPredictor // network roundtrip time + queuing delay
	paxosLatPredictor    *latency.LatPredictor // paxos latency

	latUpdateLock         sync.Mutex
	replicaLatMap         map[string]int64 // replica addr --> latency in ms
	sortedReplicaAddrList []string         // a list of replicas from low to high in latency
	//replicaLatList        []int          // roundtrip latency to each replica in ms

	predictPth float64 // >= 0 && <= 1.0 using the pth percenitle latency for prediction

	// Only for throughput experiments
	fpLoadP         int //default -1, disabled. Otherwise [0, 100] in % of using Fast Paxos
	rnd             *rand.Rand
	paxosTargetAddr string // a given target replica address when using Paxos
}

func newDynamicClientLib(id, dcId, configFile, replicaFile, targetDcId string) *DynamicClientLib {
	isReplicaColocated, coRepAddr := false, ""
	repDir := config.LoadReplicaInfo(replicaFile)
	rAddrList := make([]string, 0)
	for _, rInfo := range repDir {
		rAddrList = append(rAddrList, rInfo.GetNetAddr())
		if rInfo.DcId == dcId {
			isReplicaColocated = true
			coRepAddr = rInfo.GetNetAddr()
		}
	}

	//rList := p.GetStrList(config.FLAG_DYNAMIC_REPLICA_LIST, ",")
	//for i, r := range rList {
	//		rList[i] = p.GetStr(r)
	//}
	f := (len(rAddrList) - 1) / 2

	p := common.NewProperties()
	p.Load(configFile)

	lib := &DynamicClientLib{
		Id:                   id,
		IsReplicaColocated:   isReplicaColocated,
		ColocatedReplicaAddr: coRepAddr,

		count: 0,
		clock: common.NewSysNanoClock(),

		replicaNum:      len(rAddrList),
		majority:        f + 1,
		fastQuorum:      int(math.Ceil((3.0*float64(f))/2.0)) + 1,
		replicaAddrList: rAddrList,
		replicaAddrDir:  make(map[string][]string),

		isPaxosEnabled:              true,
		isPaxosUseLowestLatMajority: true,
		isAll:                       p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_PREDICT_ALL, "true"),
		isFpLeaderSoleLearner:       p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_LEADER_SOLELEARNER, "true"),

		clientAddDelay: int64(p.GetTimeDuration(config.FLAG_DYNAMIC_CLIENT_ADD_DELAY)),
		probeInv:       p.GetTimeDurationWithDefault(config.FLAG_DYNAMIC_LAT_PROBE_INTERVAL, "10ms"),
		probeC:         make(chan *LatInfo, 1024*16),
		probeTimeC:     make(chan *LatTimeInfo, 1024*16),
		isProbeTime:    p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_PREDICT_TIMEOFFSET, "true"),

		replicaLatMap:         make(map[string]int64),
		sortedReplicaAddrList: make([]string, len(rAddrList), len(rAddrList)),
	}

	pShardLeaderAddrMap, fpShardLeaderAddrMap := config.GenShardAddrInfo(repDir)
	if len(pShardLeaderAddrMap) == 0 {
		lib.isPaxosEnabled = false
	} else if len(pShardLeaderAddrMap) != len(rAddrList) {
		// TODO Makes random number of replicas to be the leader of Paxos instances
		logger.Fatalf("All replicas must be the leader of Paxos instances!")
	}
	if len(fpShardLeaderAddrMap) > 1 {
		// TODO Implements to choose the closest Fast Paxos leader as in the static approach
		logger.Fatalf("There must be at most one Fast Paxos leader for now!")
	} else if len(fpShardLeaderAddrMap) == 0 {
		lib.fpShard = -1
		if !lib.isPaxosEnabled {
			logger.Fatalf("No Fast Paxos or Paxos shards!")
		}
	}
	for fpShard, fpLeaderAddr := range fpShardLeaderAddrMap {
		lib.fpShard, lib.fpLeaderAddr = fpShard, fpLeaderAddr
	}

	// Replicas
	for i, rAddr := range rAddrList {
		lib.replicaLatMap[rAddr] = 1 // 1ms
		lib.sortedReplicaAddrList[i] = rAddr
		lib.replicaAddrDir[rAddr] = make([]string, 0, len(rAddrList)-1)
		lib.replicaAddrDir[rAddr] = append(lib.replicaAddrDir[rAddr], rAddrList[:i]...)
		lib.replicaAddrDir[rAddr] = append(lib.replicaAddrDir[rAddr], rAddrList[i+1:]...)
	}

	// Probing
	windowLen := p.GetTimeDurationWithDefault(config.FLAG_DYNAMIC_LAT_PROBE_WINDOW_LEN, "1s")
	windowSize := p.GetIntWithDefault(config.FLAG_DYNAMIC_LAT_PROBE_WINDOW_MIN_SIZE, "10")
	lib.fpOneWayLatPredictor = latency.NewLatPredictor(rAddrList, windowLen, windowSize)
	lib.fpLatPredictor = latency.NewLatPredictor(rAddrList, windowLen, windowSize)
	lib.paxosLatPredictor = latency.NewLatPredictor(rAddrList, windowLen, windowSize)

	// Fast Paxos client I/O lib
	isExecReply := p.GetBoolWithDefault(config.FLAG_DYNAMIC_EXEC_REPLY, "false")
	isFpLeaderUsePaxos := p.GetBoolWithDefault(config.FLAG_DYNAMIC_FP_LEADER_USE_PAXOS, "true")
	isGrpc := p.GetBoolWithDefault(config.FLAG_DYNAMIC_GRPC, "true")
	lib.c = dynamic.NewDynamicClient(
		lib.replicaNum, lib.fastQuorum, isExecReply, isFpLeaderUsePaxos, isGrpc)

	lib.predictPth = p.GetFloat64WithDefault(
		config.FLAG_DYNAMIC_PREDICT_PERCENTILE,
		config.DEFAULT_DYNAMIC_PREDICT_PERCENTILE)
	if lib.predictPth < 0.0 || lib.predictPth > 1.0 {
		logger.Fatalf("Invalid percentile = %f for prediction. Expected [0.0, 1.0]", lib.predictPth)
	}

	// Only for throughput experiments
	lib.fpLoadP = p.GetIntWithDefault(config.FLAG_DYNAMIC_CLIENT_FP_LOAD, "-1")
	if lib.fpLoadP >= 0 {
		lib.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
		for _, rInfo := range repDir {
			if rInfo.DcId == targetDcId {
				lib.paxosTargetAddr = rInfo.GetNetAddr()
				break
			}
		}
	}

	return lib
}

func (lib *DynamicClientLib) start(blocking bool, inv time.Duration) {
	if isCoPaxosLeader, _ := lib.IsPaxosLeaderColocated(); isCoPaxosLeader {
		// Simply uses the colocated Paxos leader instead of predicting network delays to replicas
		return
	}

	if lib.isProbeTime {
		lib.startTimeProcessing()
		lib.startTimeProbing(blocking, inv)
	} else {
		lib.startProcessing()
		lib.startProbing(blocking, inv)
	}
}

// Starts a thread to process probing latencies
func (lib *DynamicClientLib) startProcessing() {
	go func() {
		for latInfo := range lib.probeC {
			//logger.Infof("probe addr = %s rt = %d queing = %d pDelay = %d", latInfo.addr, latInfo.rt, latInfo.qDelay, latInfo.paxosLat)
			pLat := latInfo.rt + latInfo.paxosLat
			lib.paxosLatPredictor.AddProbeRet(&latency.ProbeRet{latInfo.addr, pLat})

			lib.fpLatPredictor.AddProbeRet(&latency.ProbeRet{latInfo.addr, latInfo.rt})

			fpOneWayLat := (latInfo.rt + latInfo.qDelay) / 2
			lib.fpOneWayLatPredictor.AddProbeRet(&latency.ProbeRet{latInfo.addr, fpOneWayLat})
		}
	}()
}

// Starts porbing thread
func (lib *DynamicClientLib) startProbing(blocking bool, inv time.Duration) {
	//logger.Infof("Predicting network latency")
	probTimer := time.NewTimer(inv)
	go func() {
		for {
			select {
			case <-probTimer.C:
				lib.probe(blocking)
				probTimer.Reset(inv)
			}
		}
	}()
}

func (lib *DynamicClientLib) probe(blocking bool) {
	var wg sync.WaitGroup
	for _, addr := range lib.replicaAddrList {
		wg.Add(1)
		go func(addr string) {
			start := time.Now()
			qDelay, paxosLat := lib.c.Probe(addr) // server-side queuing/processing delay
			rt := time.Since(start)               // network roundtrip time
			lib.probeC <- &LatInfo{
				addr:     addr,
				rt:       rt,
				qDelay:   time.Duration(qDelay),
				paxosLat: time.Duration(paxosLat * 1000000), // ms to ns
			}
			wg.Done()
		}(addr)
	}
	if blocking {
		wg.Wait()
	}
}

// Starts a thread to process probing latencies
func (lib *DynamicClientLib) startTimeProcessing() {
	//logger.Infof("Predicting time offset")
	go func() {
		for latTimeInfo := range lib.probeTimeC {
			//logger.Infof("probe addr = %s rt = %d timeOffset= %d pDelay = %d", latTimeInfo.addr, latTimeInfo.rt, latTimeInfo.timeOffset, latTimeInfo.paxosLat)
			pLat := latTimeInfo.rt + latTimeInfo.paxosLat
			lib.paxosLatPredictor.AddProbeRet(&latency.ProbeRet{latTimeInfo.addr, pLat})

			lib.fpLatPredictor.AddProbeRet(&latency.ProbeRet{latTimeInfo.addr, latTimeInfo.rt})

			lib.fpOneWayLatPredictor.AddProbeRet(&latency.ProbeRet{latTimeInfo.addr, latTimeInfo.timeOffset})
		}
	}()
}

func (lib *DynamicClientLib) startTimeProbing(blocking bool, inv time.Duration) {
	probTimer := time.NewTimer(inv)
	go func() {
		for {
			select {
			case <-probTimer.C:
				lib.probeTime(blocking)
				probTimer.Reset(inv)
			}
		}
	}()
}

func (lib *DynamicClientLib) probeTime(blocking bool) {
	var wg sync.WaitGroup
	for _, addr := range lib.replicaAddrList {
		wg.Add(1)
		go func(addr string) {
			start := time.Now()
			pTime, paxosLat := lib.c.ProbeTime(addr) // server-side queuing/processing delay
			end := time.Now()

			//logger.Infof("Probing addr = %s processT = %d start = %d offset = %d", addr, pTime, start.UnixNano(), pTime-start.UnixNano())

			lib.probeTimeC <- &LatTimeInfo{
				addr:       addr,
				rt:         end.Sub(start), // network roundtrip time
				timeOffset: time.Duration(pTime - start.UnixNano()),
				paxosLat:   time.Duration(paxosLat * 1000000), // ms to ns
			}
			wg.Done()
		}(addr)
	}
	if blocking {
		wg.Wait()
	}
}

func (lib *DynamicClientLib) Propose(cmd *dynamic.Command) (bool, bool, bool, string) {
	cmd.Id = lib.genCmdId()

	if lib.fpLoadP >= 0 {
		// Only for throuthput experiment
		return lib.thrPropose(cmd)
	}

	isFp, delay, addr := lib.selectFp()
	//logger.Infof("%s isFp = %t delay = %d addr = %s", cmd.Id, isFp, delay, addr)

	if isFp {
		// Fast Paxos
		delay += lib.clientAddDelay // adds additional fixed delay
		isCommit, isFast, val := lib.doFpPropose(cmd, delay, addr)
		return true, isCommit, isFast, val
	} else {
		// Paxos
		ok, ret := lib.doPaxosPropose(cmd, addr)
		return false, ok, false, ret
	}
}

func (lib *DynamicClientLib) thrPropose(cmd *dynamic.Command) (bool, bool, bool, string) {
	if lib.rndFastPaxos() {
		// Uses Fast Paxos
		delay, addr := lib.predictFpOneWayDelay()
		delay += lib.clientAddDelay // adds additional fixed delay
		isCommit, isFast, val := lib.doFpPropose(cmd, delay, addr)
		return true, isCommit, isFast, val
	}
	// Uses Paxos to the specified replica for this client
	ok, ret := lib.doPaxosPropose(cmd, lib.paxosTargetAddr)
	return false, ok, false, ret
}

func (lib *DynamicClientLib) FpPropose(cmd *dynamic.Command) (bool, bool, string) {
	// TODO Implementation
	logger.Fatalf("Not implemented yet!")
	return false, false, ""
}

func (lib *DynamicClientLib) PaxosPropose(cmd *dynamic.Command) (bool, string) {
	// TODO Implementation
	logger.Fatalf("Not implemented yet!")
	return false, ""
}

// Fast Paxos
func (lib *DynamicClientLib) doFpPropose(cmd *dynamic.Command, delay int64, addr string) (bool, bool, string) {
	shard, leaderAddr := lib.getFpShardAndLeader()

	if lib.isFpLeaderSoleLearner {
		followerAddrList := lib.getOtherReplicaAddrList(leaderAddr)
		t := lib.clock.GetClockTime() + delay
		return lib.c.FastPaxosPropose(cmd, shard, leaderAddr, followerAddrList, t)
	} else {
		fpExecReplicaAddr := addr
		fpNonExecReplicaAddrList := lib.getOtherReplicaAddrList(fpExecReplicaAddr)
		t := lib.clock.GetClockTime() + delay
		return lib.c.ExecFastPaxosPropose(cmd, shard, leaderAddr,
			fpExecReplicaAddr, fpNonExecReplicaAddrList, t)
	}
}

// Paxos
func (lib *DynamicClientLib) doPaxosPropose(cmd *dynamic.Command, leaderAddr string) (bool, string) {
	return lib.c.PaxosPropose(cmd, leaderAddr)
}

func (lib *DynamicClientLib) Close() {
	lib.c.Close()
}

func (lib *DynamicClientLib) Test() {
	lib.c.Test(lib.replicaAddrList)
}

// Returns true and colocated paxos leasder's net addr
// TODO: The colocated replica may not be a Paxos leader
func (lib *DynamicClientLib) IsPaxosLeaderColocated() (bool, string) {
	if lib.isPaxosEnabled && lib.IsReplicaColocated {
		return true, lib.ColocatedReplicaAddr
	}
	return false, ""
}

// Returns isUsingFastPaxos, and the delay time for the next command
func (lib *DynamicClientLib) selectFp() (bool, int64, string) {
	lib.latUpdateLock.Lock()
	defer lib.latUpdateLock.Unlock()

	if isCoPaxosLeader, paxosLeaderAddr := lib.IsPaxosLeaderColocated(); isCoPaxosLeader {
		// Simply uses the colocated
		return false, 0, paxosLeaderAddr
	}

	paxosLeader, paxosLat := lib.predictPaxosLat() // Paxos latency
	if lib.fpShard == -1 {
		// No Fast Paxos shard
		return false, 0, paxosLeader
	}

	// TODO Combines predictFpLat() and predictFpOneWayDelay() to avoid sorting latencies twice
	fpLat := lib.predictFpLat() // Fast Paxos latency
	//logger.Infof("Prediction Paxos leader = %s lat = %d, FP lat = %d", paxosLeader, paxosLat, fpLat)
	if lib.isPaxosEnabled && paxosLat <= fpLat {
		// Uses Paxos
		return false, 0, paxosLeader
	}

	// Predicts Fast Paxos timestamp
	fpLat, fpExecReplicaAddr := lib.predictFpOneWayDelay()

	return true, fpLat * 1000000, fpExecReplicaAddr
}

// Returns the addr of the replica that achieves the lowest Paxos latency, and the latency in ms
func (lib *DynamicClientLib) predictPaxosLat() (string, int64) {
	leader := lib.replicaAddrList[0]
	minLat := lib.paxosLatPredictor.PredictLat(leader, lib.predictPth)
	//logger.Infof("paxos leader addr = %s, lat = %d", leader, minLat)
	for i := 1; i < len(lib.replicaAddrList); i++ {
		lat := lib.paxosLatPredictor.PredictLat(lib.replicaAddrList[i], lib.predictPth)
		//logger.Infof("paxos leader addr = %s, lat = %d", lib.replicaAddrList[i], lat)
		if lat < minLat {
			leader, minLat = lib.replicaAddrList[i], lat
		}
	}
	return leader, minLat
}

// Returns Fast Paxos latency in ms
func (lib *DynamicClientLib) predictFpLat() int64 {
	latList := make([]int64, lib.replicaNum, lib.replicaNum)
	for i, addr := range lib.replicaAddrList {
		latList[i] = lib.fpLatPredictor.PredictLat(addr, lib.predictPth)
	}
	common.BubbleSort64n(latList)
	return latList[lib.fastQuorum-1]
}

// Returns the one-way delay (in ms) for Fast Paxos timestamp, and the fpExecReplicaAddr
func (lib *DynamicClientLib) predictFpOneWayDelay() (int64, string) {
	noUpdate := true
	for _, addr := range lib.replicaAddrList {
		lat := lib.fpOneWayLatPredictor.PredictLat(addr, lib.predictPth)
		if lat != lib.replicaLatMap[addr] {
			lib.replicaLatMap[addr] = lat
			if noUpdate {
				noUpdate = false
			}
		}
	}

	if !noUpdate {
		lib.sortReplicaByLat()
	}

	// Determines the estimated latency for using Fast Paxos
	fpAddr := lib.sortedReplicaAddrList[lib.replicaNum-1]
	if !lib.isAll {
		fpAddr = lib.sortedReplicaAddrList[lib.fastQuorum-1]
	}
	fpLat := lib.replicaLatMap[fpAddr]

	// TODO the closest replica within the closest fast quorum may not achieve
	// the lowest execution latency.
	// Calculates the expected execution latency for every replica if all
	// replicas are learners (i.e., broadcasting is enabled), otherwise the
	// leader should be the expected execution replica which is configured on the
	// leader replica.
	// NOTE: currently, every replica would send its execution result back when
	// every replica is a learner.  This approach can achive lowest execution
	// latency but cost more network traffict. In this approach, this specified
	// execution replica is only to deal with rejected requests.
	return fpLat, lib.sortedReplicaAddrList[0]
}

func (lib *DynamicClientLib) sortReplicaByLat() {
	for i := 0; i < len(lib.sortedReplicaAddrList); i++ {
		noSwap := true
		for j := 0; j+1+i < len(lib.sortedReplicaAddrList); j++ {
			a, b := lib.sortedReplicaAddrList[j], lib.sortedReplicaAddrList[j+1]
			if lib.replicaLatMap[a] > lib.replicaLatMap[b] {
				lib.sortedReplicaAddrList[j], lib.sortedReplicaAddrList[j+1] = b, a
				noSwap = false
			}
		}
		if noSwap {
			return
		}
	}
}

///// Helper Functions /////
// Returns an unique command id for this client lib
// Thread-safe
func (lib *DynamicClientLib) genCmdId() string {
	lib.countLock.Lock()
	id := fmt.Sprintf("%s-%d", lib.Id, lib.count)
	lib.count++
	lib.countLock.Unlock()
	return id
}

//// Fast Paxos
// Returns the shard and network address of the chosen Fast Paxos leader
func (lib *DynamicClientLib) getFpShardAndLeader() (int32, string) {
	return lib.fpShard, lib.fpLeaderAddr
}

// Returns all of the replicas' network addresses except the given address
func (lib *DynamicClientLib) getOtherReplicaAddrList(addr string) []string {
	return lib.replicaAddrDir[addr]
}

func (lib *DynamicClientLib) rndFastPaxos() bool {
	if lib.rnd.Intn(100) < lib.fpLoadP {
		return true
	}
	return false
}
