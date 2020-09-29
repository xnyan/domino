package dynamic

import (
	"io"
	"net"
	"sync"
	"time"
)

type FpReply struct {
	count int
	fastC chan<- *FpProposeReply
	slowC chan<- *FpProposeReply
}

type ClientFastIo struct {
	serverNetIo  map[string]NetIo            // addr --> net I/O
	serverStream map[string]chan interface{} // addr --> sending channel

	paxosReply     map[string]chan *PaxosProposeReply
	paxosReplyLock sync.Mutex

	fpReply     map[string]*FpReply
	fpReplyLock sync.Mutex

	wg *sync.WaitGroup
}

func NewClientFastIo(wg *sync.WaitGroup) ClientIo {
	fIo := &ClientFastIo{
		serverNetIo:  make(map[string]NetIo),
		serverStream: make(map[string]chan interface{}),
		paxosReply:   make(map[string]chan *PaxosProposeReply),
		fpReply:      make(map[string]*FpReply),
		wg:           wg,
	}
	return fIo
}

func (fIo *ClientFastIo) InitConn(addrList []string) {
	//var wg sync.WaitGroup
	for _, addr := range addrList {
		//wg.Add(1)
		//go func(addr string) {
		retry := 10
		for {
			if conn, err := net.Dial("tcp", addr); err == nil {
				fIo.serverNetIo[addr] = NewNetIo(conn)
				break
			} else {
				time.Sleep(2e8)
				if retry--; retry <= 0 {
					logger.Fatalf("Cannot connect to server %s error %v", addr, err)
				}
			}
		}
		//wg.Done()
		//}(addr)
	}
	//wg.Wait()

	bufferSize := 10240 * 8 // TODO Gets rid of hard coding
	// Starts reader threads
	// Single thread for each server to read messages, dispatching messages for each RPC calls
	for addr, netIo := range fIo.serverNetIo {
		go fIo.waitForServerReply(addr, netIo)
	}

	// Starts writer
	for addr, netIo := range fIo.serverNetIo {
		c := make(chan interface{}, bufferSize)
		fIo.serverStream[addr] = c
		// Single thread for each server to send messages
		go fIo.startServerStream(addr, netIo, c)
	}
}

func (fIo *ClientFastIo) startServerStream(addr string, netIo NetIo, c <-chan interface{}) {
	for msg := range c {
		switch msg.(type) {
		case *PaxosProposeReq:
			netIo.SendMsg(Msg_Type_PaxosProposeReq, msg)
		case *FpProposeReq:
			netIo.SendMsg(Msg_Type_FpProposeReq, msg)
		case *TestReq:
			netIo.SendMsg(Msg_Type_TestReq, msg)
		default:
			logger.Fatalf("Unknown message %v", msg)
		}
	}
}

func (fIo *ClientFastIo) waitForServerReply(addr string, netIo NetIo) {
	for {
		//logger.Infof("Waiting replies from server %s", addr)
		msgType, msg, err := netIo.RecvMsg()
		if err != nil {
			if err != io.EOF {
				logger.Fatalf("Receiving message from %s error: %v", addr, err)
			} else {
				logger.Debugf("Connection to server %s shuts down", addr)
				break
			}
		}

		switch msgType {
		case Msg_Type_PaxosProposeReply:
			reply := msg.(*PaxosProposeReply)
			//logger.Infof("Received Paxos reply from %s for cmdId = %s", addr, reply.CmdId)
			go func(reply *PaxosProposeReply) {
				fIo.paxosReplyLock.Lock()
				c := fIo.paxosReply[reply.CmdId]
				delete(fIo.paxosReply, reply.CmdId)
				fIo.paxosReplyLock.Unlock()

				c <- reply
				//logger.Infof("Delivered Paxos reply from %s for cmdId = %s", addr, reply.CmdId)
			}(reply)
		case Msg_Type_FpProposeReply:
			fIo.wg.Done()
			reply := msg.(*FpProposeReply)
			//logger.Infof("Received FP reply from %s for cmdId = %s", addr, reply.CmdId)
			go func(reply *FpProposeReply) {
				fIo.fpReplyLock.Lock()
				r := fIo.fpReply[reply.CmdId]
				if r.count--; r.count <= 0 {
					delete(fIo.fpReply, reply.CmdId)
				}
				fIo.fpReplyLock.Unlock()

				if reply.IsFast {
					r.fastC <- reply
					//logger.Infof("Delivered FP fast reply from %s for cmdId = %s isFast = %t", addr, reply.CmdId, reply.IsFast)
				} else {
					//logger.Infof("Delivered FP slow reply from %s for cmdId = %s isFast = %t", addr, reply.CmdId, reply.IsFast)
					r.slowC <- reply
				}
			}(reply)
		default:
			logger.Fatalf("Unknown message type %d", msgType)
		}
	}
}

func (fIo *ClientFastIo) SendProbeReq(addr string) *ProbeReply {
	logger.Fatalf("ClientFastIo has not implemented probing yet!!! Use gRPC instead.")
	return nil
}

func (fIo *ClientFastIo) SendProbeTimeReq(addr string) *ProbeTimeReply {
	logger.Fatalf("ClientFastIo has not implemented probing yet!!! Use gRPC instead.")
	return nil
}

// Blocking
func (fIo *ClientFastIo) SendPaxosProposeReq(addr string, cmd *Command) *PaxosProposeReply {
	fIo.paxosReplyLock.Lock()
	replyC := make(chan *PaxosProposeReply, 1)
	fIo.paxosReply[cmd.Id] = replyC
	fIo.paxosReplyLock.Unlock()

	req := &PaxosProposeReq{Cmd: cmd}
	streamC := fIo.serverStream[addr]
	streamC <- req

	// Waits for reply
	reply := <-replyC
	return reply
}

// Non-blocking
func (fIo *ClientFastIo) BcstFpProposeToFollowers(
	addrList []string, req *FpProposeReq, fastC chan<- *FpProposeReply, wg *sync.WaitGroup,
) {
	fIo.fpReplyLock.Lock()
	if fpReply, ok := fIo.fpReply[req.Cmd.Id]; !ok {
		fIo.fpReply[req.Cmd.Id] = &FpReply{
			count: len(addrList),
			fastC: fastC,
		}
	} else {
		fpReply.count += len(addrList)
	}
	fIo.fpReplyLock.Unlock()

	for _, addr := range addrList {
		fIo.wg.Add(1)
		streamC := fIo.serverStream[addr]
		streamC <- req
	}
}

// Non-blocking
func (fIo *ClientFastIo) SendFpProposeToLeader(
	leader string, req *FpProposeReq, fastC, slowC chan<- *FpProposeReply, wg *sync.WaitGroup,
) {
	fIo.fpReplyLock.Lock()
	if fpReply, ok := fIo.fpReply[req.Cmd.Id]; !ok {
		fIo.fpReply[req.Cmd.Id] = &FpReply{
			count: 2,
			fastC: fastC,
			slowC: slowC,
		}
	} else {
		fpReply.count += 2
		fpReply.slowC = slowC
	}
	fIo.fpReplyLock.Unlock()

	fIo.wg.Add(2)
	streamC := fIo.serverStream[leader]
	streamC <- req
}

// Non-blocking
func (io *ClientFastIo) SendFpProposeToExecReplica(
	addr string, req *FpProposeReq, fastC, slowC chan<- *FpProposeReply, wg *sync.WaitGroup,
) {
	logger.Fatalf("Not implemented for fast io")
}

func (fIo *ClientFastIo) SendTestReq(addr string) {
	req := &TestReq{}
	streamC := fIo.serverStream[addr]
	streamC <- req
}
