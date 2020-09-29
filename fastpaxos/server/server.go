package main

import (
	"net"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"domino/common"
	"domino/fastpaxos/com"
	fp "domino/fastpaxos/fastpaxos"
	"domino/fastpaxos/rpc"
)

type Server struct {
	Id         string
	NetAddr    string // ip:port
	Ip         string
	Port       string
	grpcServer *grpc.Server

	IsLeader         bool
	LeaderAddr       string // ip:port
	FollowerAddrList []string
	ReplicaAddrList  []string // the addr of all replicas

	ReplicaNum  int
	MajorityNum int
	FastQuorum  int

	// Communincation channel manager
	cm *com.CommManager

	// Fast Paxos
	fp *fp.FastPaxos

	// Key-value store
	kvStore          *KvStore
	execMode         bool
	execRetTable     map[int64]*common.Future // timestamp --> exec result
	execRetTableLock sync.RWMutex
}

func NewServer(
	id, netAddr string, isLeader bool,
) *Server {
	server := &Server{
		Id:      id,
		NetAddr: netAddr,

		IsLeader:         isLeader,
		LeaderAddr:       LeaderAddr,
		FollowerAddrList: FollowerAddrList,

		ReplicaNum:  ReplicaNum,
		MajorityNum: MajorityNum,
		FastQuorum:  FastQuorum,

		cm: com.NewCommManager(),

		kvStore:      NewKvStore(),
		execRetTable: make(map[int64]*common.Future),
	}

	logger.Debugf("Server leader addr: %s", server.LeaderAddr)
	logger.Debugf("Server follower addr list: %s", server.FollowerAddrList)

	server.ReplicaAddrList = append(server.FollowerAddrList, server.LeaderAddr)

	ipPort := strings.Split(server.NetAddr, ":")
	if len(ipPort) != 2 {
		logger.Fatalf("Invalid server network address = %s", server.NetAddr)
	}
	server.Ip = ipPort[0]
	server.Port = ipPort[1]

	server.InitFastPaxos()

	// Creates the RPC server instance
	server.grpcServer = grpc.NewServer()
	rpc.RegisterFastPaxosRpcServer(server.grpcServer, server)
	reflection.Register(server.grpcServer)

	return server
}

func (server *Server) InitFastPaxos() {
	// Log
	var l fp.Log
	switch FastPaxosLogType {
	case common.FixedLog:
		l = fp.NewFixedLog(FastPaxosLogSize)
	case common.SegLog:
		// TODO uses configurable # of segs instead of hard coding
		l = fp.NewSegLog(1, FastPaxosLogSize)
	default:
		logger.Fatalf("Invalid log type: %s", FastPaxosLogType)
	}

	// Log persistence
	pLogger := fp.NewLogger(FastPaxosLogFilePath)

	// Log Manager
	var lm fp.LogManager = fp.NewDefaultLogManager(l, FastPaxosLogPersistEnabled, pLogger)

	// Fast Paxos
	server.fp = fp.NewFastPaxos(
		ReplicaNum,
		MajorityNum,
		FastQuorum,
		FastPaxosVoteChannelBufferSize,
		FastPaxosCommadChannelBufferSize,
		FastPaxosExecutionChannelBufferSize,
		lm,
		FastPaxosScheduler,
		FastPaxosProcessWindow,
	)
}

func (server *Server) Start() {
	logger.Infof("Starting Fast Paxos service")
	server.runFastPaxos(server.IsLeader)
	logger.Infof("Started Fast Paxos service")

	logger.Infof("Starting Key-value store")
	server.initKvStore()
	server.runKvStore()
	logger.Infof("Started Key-value store")

	// Starts RPC service
	rpcListener, err := net.Listen("tcp", ":"+server.Port)
	if err != nil {
		logger.Errorf("Fails to listen on port %s \nError: %v", server.Port, err)
	}

	logger.Infof("Starting RPC services")

	err = server.grpcServer.Serve(rpcListener)
	if err != nil {
		logger.Errorf("Cannot start RPC services. \nError: %v", err)
	}
}

func (server *Server) runFastPaxos(isFastPathManager bool) {
	server.fp.Run(isFastPathManager)
}

func (server *Server) initKvStore() {
	keyList := common.LoadKey(DataKeyFile)
	val := common.GenVal(DataValSize)
	server.kvStore.InitData(keyList, val)
}

func (server *Server) runKvStore() {
	server.runKvStoreWithFp()
}

func (server *Server) runKvStoreWithFp() {
	execCh := server.fp.GetExecCh()

	if FastPaxosIsExec {
		go func() {
			for op := range execCh {
				server.execOp(op)
			}
		}()
	}
}

func (server *Server) execOp(op *rpc.Operation) {
	if op == nil {
		logger.Fatalf("Operation applied to kv-store should not be nil")
	}

	if op.Type == fp.OP_WRITE {
		server.kvStore.Write(op.Key, op.Val)
	} else if op.Type == fp.OP_READ {
		// Performs a read
		server.kvStore.Read(op.Key)
	} else if op.Type == fp.OP_STOP {
		logger.Infof("Execution stopped")
	}
}
