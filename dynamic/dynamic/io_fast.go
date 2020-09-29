package dynamic

import (
	"io"
	"net"
	"strconv"
	"time"
)

type FastIo struct {
	*rpcIo

	PeerAddrList []string
	PeerNetIo    []NetIo
	Id           int // replica Id starting from 0
	N            int // number of replicas
	Listener     net.Listener

	paxos  Paxos   // Replica streaming
	server *Server // Client RPCs

	ready chan bool

	syncPeerNetIoMap map[string]NetIo
}

func NewFastIo(isSyncSend bool, id string, num int, nodeAddrList []string, p Paxos) *FastIo {
	i, err := strconv.Atoi(id)
	if err != nil {
		logger.Fatalf("Invalid id = %s, err = %v", id, err)
	}
	fIo := &FastIo{
		rpcIo: newRpcIo(isSyncSend),

		PeerAddrList: nodeAddrList,
		PeerNetIo:    make([]NetIo, num),
		Id:           i - 1,
		N:            num,

		ready: make(chan bool, 1),

		syncPeerNetIoMap: make(map[string]NetIo),
	}
	fIo.paxos = p
	return fIo
}

func (fIo *FastIo) SetServer(s *Server) {
	fIo.server = s
}

// Mimic the implementation of EPaxos RPC
func (fIo *FastIo) InitConn(addrList []string, streamBufSize int) {
	done := make(chan bool)
	go fIo.waitForPeerConnections(done)

	//connect to peers
	for i := 0; i < int(fIo.Id); i++ {
		addr := fIo.PeerAddrList[i]
		for {
			if conn, err := net.Dial("tcp", addr); err == nil {
				fIo.PeerNetIo[i] = NewNetIo(conn)
				break
			} else {
				time.Sleep(1e9)
			}
		}
		fIo.PeerNetIo[i].SendByte(uint8(fIo.Id))
	}

	<-done // barrier

	logger.Infof("Replica id: %d. Done connecting to peers\n", fIo.Id)

	// Listens on replica servers
	for i, netIo := range fIo.PeerNetIo {
		if i == fIo.Id {
			continue
		}
		addr := fIo.PeerAddrList[i]
		go fIo.replicaListener(i, addr, netIo)
	}

	// Sets up streaming to replica servers
	for i, netIo := range fIo.PeerNetIo {
		if i == fIo.Id {
			continue
		}
		addr := fIo.PeerAddrList[i]
		if fIo.syncSend {
			fIo.setSyncPeerNetIo(addr, netIo)
		} else {
			c := make(chan *ReplicaMsg, streamBufSize)
			fIo.setStreamCh(addr, c)
			go fIo.replicaStream(addr, c, netIo)
		}
	}
	fIo.ready <- true
}

func (fIo *FastIo) setSyncPeerNetIo(addr string, netIo NetIo) {
	fIo.syncPeerNetIoMap[addr] = netIo
}

func (fIo *FastIo) waitForPeerConnections(done chan bool) {
	var err error
	if fIo.Listener, err = net.Listen("tcp", fIo.PeerAddrList[fIo.Id]); err != nil {
		logger.Fatalf("TCP Listen error %v", err)
	}

	for i := fIo.Id + 1; i < fIo.N; i++ {
		conn, err := fIo.Listener.Accept()
		if err != nil {
			logger.Errorf("Accept error: %v", err)
			continue
		}

		netIo := NewNetIo(conn)
		id, err := netIo.RecvByte()
		if err != nil {
			logger.Errorf("Read id error: %v", err)
			continue
		}
		fIo.PeerNetIo[int(id)] = netIo
	}

	done <- true
}

func (fIo *FastIo) replicaListener(i int, addr string, netIo NetIo) {
	logger.Infof("Waiting messages from replica %d %s", i, addr)
	for {
		msgType, msg, err := netIo.RecvMsg()
		if err != nil {
			if err == io.EOF {
				logger.Infof("Connection with %d %s shuts down", i, addr)
			} else {
				logger.Fatalf("Receiving message from %d %s error: %v", i, addr, err)
			}
			break
		}

		if msgType != Msg_Type_ReplicaMsg {
			logger.Fatalf("Cannot process message type %d, message %v", msgType, msg)
		}
		fIo.server.paxos.ProcessReplicaMsg(msg.(*ReplicaMsg))
	}
}

func (fIo *FastIo) replicaStream(addr string, c <-chan *ReplicaMsg, netIo NetIo) {
	for msg := range c {
		err := netIo.SendMsg(Msg_Type_ReplicaMsg, msg)
		if err != nil {
			logger.Fatalf("Fails sending to %s message %v", addr, msg)
		}
	}
}

func (io *FastIo) BcstReplicaMsg(addrList []string, msg *ReplicaMsg) {
	for _, addr := range addrList {
		io.SendReplicaMsg(addr, msg)
	}
}

func (io *FastIo) SendReplicaMsg(addr string, msg *ReplicaMsg) {
	if io.syncSend {
		io.syncSendReplicaMsg(addr, msg)
	} else {
		c := io.getStreamCh(addr)
		c <- msg
	}
}

func (fIo *FastIo) SendReplicaProbeReq(addr string) *ReplicaProbeReply {
	// TODO implementation
	logger.Fatalf("SendProbeReq() not implemented yet!")
	return nil
}

// Returns after sending the message
func (fIo *FastIo) syncSendReplicaMsg(addr string, msg *ReplicaMsg) {
	netIo, ok := fIo.syncPeerNetIoMap[addr]
	if !ok {
		logger.Fatalf("Cannot find NetIo for addr = %s", addr)
	}
	err := netIo.SendMsg(Msg_Type_ReplicaMsg, msg)
	if err != nil {
		logger.Fatalf("Fails sending to %s message %v", addr, msg)
	}
}

////////////////////////////////////////
// Handles client RPCs
func (fIo *FastIo) WaitForClientConn() {
	<-fIo.ready
	for {
		conn, err := fIo.Listener.Accept()
		if err != nil {
			logger.Errorf("Accept error: %v", err)
			continue
		}
		go fIo.clientListener(conn)
	}
}

func (fIo *FastIo) clientListener(conn net.Conn) {
	netIo := NewNetIo(conn)
	syncNetIo := NewSyncNetIo(netIo)
	for {
		msgType, msg, err := netIo.RecvMsg()
		if err != nil {
			if err != io.EOF {
				logger.Fatalf("Receiving client message error: %v", err)
			}
			break
		}

		switch msgType {
		case Msg_Type_PaxosProposeReq:
			go fIo.handlePaxosPropose(msg.(*PaxosProposeReq), syncNetIo)
		case Msg_Type_FpProposeReq:
			go fIo.handleFpPropose(msg.(*FpProposeReq), syncNetIo)
		case Msg_Type_TestReq:
			go fIo.handleTestReq(msg.(*TestReq), syncNetIo)
		default:
			logger.Fatalf("Cannot process message type %d, message %v", msgType, msg)
		}
	}
}

func (fIo *FastIo) handlePaxosPropose(req *PaxosProposeReq, syncNetIo *SyncNetIo) {
	isCommit := fIo.server.paxos.PaxosLeaderAccept(req.Cmd)

	execRet := ""
	if isCommit && fIo.server.IsExecReply {
		// Waits for execution result
		execRet = fIo.server.em.WaitExecRet(req.Cmd.Id)
	}

	rep := &PaxosProposeReply{
		IsCommit: isCommit,
		ExecRet:  execRet,
		CmdId:    req.Cmd.Id,
	}

	syncNetIo.SendMsg(Msg_Type_PaxosProposeReply, rep)
}

func (fIo *FastIo) handleFpPropose(req *FpProposeReq, syncNetIo *SyncNetIo) {
	isAccept := fIo.server.paxos.FpReplicaAccept(req.Cmd, req.Time)
	fastReply := &FpProposeReply{
		IsAccept: isAccept,
		IsFast:   true,
		CmdId:    req.Cmd.Id,
	}
	if err := syncNetIo.SendMsg(Msg_Type_FpProposeReply, fastReply); err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending fast-path reply.")
	}

	if req.Time.Shard == fIo.server.paxos.GetCurFpShard() {
		// Fast Paxos shard leader (coordinator) waits for consensus result
		isCommit, isFast := fIo.server.paxos.FpWaitConsensus(req.Time.Time, req.Cmd.Id) // blocking

		if fIo.server.IsFpLeaderUsePaxos {
			if !isCommit && fIo.server.paxos.GetCurPaxosShard() >= 0 {
				// Uses Paxos shard to accept a command that is rejected by Fast Paxos
				isCommit = fIo.server.paxos.PaxosLeaderAccept(req.Cmd)
			}
		}

		ret := ""
		if isCommit && fIo.server.IsExecReply {
			// Waits for execution result
			ret = fIo.server.em.WaitExecRet(req.Cmd.Id)
		}

		// Sends slow-path and/or execution reply to the client
		rep := &FpProposeReply{IsAccept: isCommit, IsFast: isFast, ExecRet: ret, CmdId: req.Cmd.Id}
		if err := syncNetIo.SendMsg(Msg_Type_FpProposeReply, rep); err != nil {
			logger.Errorf("Error: %v", err)
			logger.Fatalf("Fails sending consensus result for cmdId = %s, t = %v",
				req.Cmd.Id, req.Time)
		}
	}

}

func (fIo *FastIo) handleTestReq(msg *TestReq, syncNetIo *SyncNetIo) {
	fIo.server.doTest()
}
