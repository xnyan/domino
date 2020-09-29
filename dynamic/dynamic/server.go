package dynamic

import (
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"domino/common"
	"domino/dynamic/latency"
)

type Server struct {
	// Network
	grpcServer *grpc.Server

	// Consenus
	paxos      Paxos
	paxosShard int32

	IsExecReply bool // True if reply execution ret
	em          *ExecManager

	IsFpLeaderUsePaxos bool
	IsFpLeaderLearner  bool

	IsGrpc bool

	// Latency predictions
	IsLatPrediction  bool
	followerAddrList []string      // other replicas' addr as a list
	majorityIdx      int           // the majority idx in the followers' list
	probeInv         time.Duration // probing interval
	probeC           chan *latency.ProbeRet
	latPredictor     *latency.LatPredictor
	predictPth       float64 // using the pth percentile for prediction. >= 0 && <= 1.0
}

func NewServer(
	replicaId string,
	replicaIdAddrMap map[string]string,
	rIdList, nodeAddrList []string,
	pShard, fpShard int32,
	pShardLeaderMap, fpShardLeaderMap map[int32]string,
	followerAddrList []string,
	hbInterval time.Duration,
	cmdChSize int,
	execChSize int,
	isExecReply bool,
	isFpLeaderUsePaxos bool,
	isGrpc, isSyncSend bool,
	isFpLeaderLearner bool,
	isLatPrediction bool,
	probeInv time.Duration,
	windowLen time.Duration, windowSize int,
	predictPth float64,
	isPaxosFutureTime bool, // Execution latency optimization switch
) *Server {
	s := &Server{
		paxos: NewDynamicPaxos(
			replicaId,
			replicaIdAddrMap,
			rIdList, nodeAddrList,
			pShard, fpShard,
			pShardLeaderMap, fpShardLeaderMap,
			followerAddrList,
			hbInterval,
			cmdChSize,
			execChSize,
			isGrpc, isSyncSend,
			isFpLeaderLearner,
		),
		paxosShard:         pShard,
		IsExecReply:        isExecReply,
		em:                 NewExecmanager(),
		IsFpLeaderUsePaxos: isFpLeaderUsePaxos,
		IsFpLeaderLearner:  isFpLeaderLearner,
		IsGrpc:             isGrpc,
	}
	s.IsLatPrediction = isLatPrediction
	s.followerAddrList = make([]string, len(followerAddrList), len(followerAddrList))
	copy(s.followerAddrList, followerAddrList)
	s.majorityIdx = (len(s.followerAddrList)+1)>>1 - 1
	s.probeInv = probeInv
	s.probeC = make(chan *latency.ProbeRet, 10240*len(s.followerAddrList))
	s.latPredictor = latency.NewLatPredictor(s.followerAddrList, windowLen, windowSize)
	s.predictPth = predictPth

	if isPaxosFutureTime {
		s.paxos.EnablePaxosFutureTime(s)
	}

	return s
}

func (s *Server) Start(port string) {
	// Starts the consensus protocol
	// Time to wait for other servers to init before sending the first heart beat
	s.paxos.Start(time.Duration(DEFAULT_FIRST_HEART_BEAT_DELAY))

	// Latency probing
	if s.IsLatPrediction {
		if s.paxosShard >= 0 {
			go s.startLatPrediction(false, s.probeInv)
		}
	}

	// Starts rpc server
	if s.IsGrpc {
		s.startGRpc(port)
	} else {
		s.startFastRpc()
	}
}

func (s *Server) startLatPrediction(blocking bool, inv time.Duration) {
	time.Sleep(time.Duration(DEFAULT_FIRST_HEART_BEAT_DELAY))
	s.startLatProcessing()
	s.startLatProbing(blocking, inv)
}

func (s *Server) startLatProcessing() {
	go func() {
		for pr := range s.probeC {
			s.latPredictor.AddProbeRet(pr)
		}
	}()
}

func (s *Server) startLatProbing(blocking bool, inv time.Duration) {
	go func() { // Probing thread
		probTimer := time.NewTimer(inv)
		for {
			select {
			case <-probTimer.C:
				if blocking {
					s.probe() // blocking
				} else {
					go s.probe()
				}
				probTimer.Reset(inv)
			}
		}
	}()
}

// Blocking
func (s *Server) probe() {
	io := s.paxos.GetIo()
	var wg sync.WaitGroup

	for _, addr := range s.followerAddrList {
		wg.Add(1)
		go func(addr string) {
			start := time.Now()
			io.SyncSendReplicaProbeReq(addr)      // blocking
			rt := time.Since(start).Nanoseconds() // network roundtrip time + processing delay
			s.probeC <- &latency.ProbeRet{addr, time.Duration(rt)}
			wg.Done()
		}(addr)
	}

	wg.Wait()
}

// Predicts the Paxos latency (in ms) when this is the leader
func (s *Server) PredictPaxosLat() int64 {
	latList := make([]int64, len(s.followerAddrList), len(s.followerAddrList))
	for i, addr := range s.followerAddrList {
		latList[i] = s.latPredictor.PredictLat(addr, s.predictPth)
	}
	common.BubbleSort64n(latList)
	return latList[s.majorityIdx]
}

func (s *Server) startFastRpc() {
	s.paxos.(*DynamicPaxos).io.(*StreamIo).rpcIo.(*FastIo).SetServer(s)
	s.paxos.(*DynamicPaxos).io.(*StreamIo).rpcIo.(*FastIo).WaitForClientConn()
}

func (s *Server) startGRpc(port string) {
	// Creates the gRPC server instance
	s.grpcServer = grpc.NewServer()
	RegisterDynamicPaxosServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	// Starts RPC service
	rpcListener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Errorf("Fails to listen on port %s \nError: %v", port, err)
	}

	logger.Infof("Starting gRPC services")

	err = s.grpcServer.Serve(rpcListener)
	if err != nil {
		logger.Errorf("Cannot start gRPC services. \nError: %v", err)
	}
}

func (s *Server) GetExecCh() <-chan Entry {
	return s.paxos.GetExecCh()
}

func (s *Server) GetExecManager() *ExecManager {
	return s.em
}

func (s *Server) doTest() {
	s.paxos.Test()

	fmt.Println("execRetMap size =", s.em.Size())
}
