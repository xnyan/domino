package fastpaxos

import (
	"sync"
	"time"
)

type Scheduler interface {
	Schedule(*ClientProposal)
	Run()
}

type AbstractScheduler struct {
	Scheduler // interface

	outCh chan<- interface{}
}

func NewAbstractScheduler(ch chan<- interface{}) *AbstractScheduler {
	s := &AbstractScheduler{
		outCh: ch,
	}
	return s
}

//////////////////////////////
// No Scheduler, which is just the basic Fast Paxos
type NoScheduler struct {
	*AbstractScheduler
}

func (s *NoScheduler) Schedule(p *ClientProposal) {
	s.outCh <- p
}

func (s *NoScheduler) Run() {
	logger.Infof("NoScheduler starts")
}

//////////////////////////////
// Delay Scheduler, where delays a proposal based on the latency between DCs
type DelayScheduler struct {
	*AbstractScheduler
}

func (s *DelayScheduler) Schedule(p *ClientProposal) {
	// The delay time is given by a client since the client knows which replicas
	// it sends requests to, and it is the best candidate to monitor the network
	// delays from itself to the replicas.

	logger.Debugf("Client proposal opId = (%s) delay %v", p.Op.Id, p.Delay)

	t := time.NewTimer(p.Delay)
	<-t.C
	s.outCh <- p
}

func (s *DelayScheduler) Run() {
	logger.Infof("DelayScheduler starts")
}

//////////////////////////////
// Timestamp scheduler, where clients' proposals are sorted based on their
// timestamp within a time window
type TimestampScheduler struct {
	*AbstractScheduler

	Queue			*PriorityQueue
	queueLock		*sync.Mutex
	modifiedFlag 	SynchronizedFlag
}

const defaultTimeoutDuration = int64(10 * time.Second)

func NewTimestampScheduler(absScheduler *AbstractScheduler) *TimestampScheduler {
	return &TimestampScheduler{
		AbstractScheduler: 	absScheduler,

		Queue:			NewPriorityQueue(),
		queueLock:		&sync.Mutex{},
		modifiedFlag:	NewSynchronizedFlag(),
	}
}

func (s *TimestampScheduler) Schedule(p *ClientProposal) {
	s.addProposal(p)
}

func (s *TimestampScheduler) Run() {
	logger.Infof("TimestampScheduler starts")

	timer := time.NewTimer(0)

	for {
		s.queueLock.Lock()

		//logger.Debugf("TimestampScheduler: Clearing modifiedFlag (%v)", s.modifiedFlag.clearFlag())
		s.modifiedFlag.clearFlag()

		// Query current system time
		sysT := time.Now().UnixNano()

		// Process all elapsed proposals
		//logger.Debugf("TimestampScheduler: Processing elapsed proposals (sysT = %v)", sysT)
		var timeoutTimestamp int64
		for {
			if p := s.Queue.Peek(); p == nil {
				// No proposal
				//logger.Debugf("TimestampScheduler: No proposal in queue")
				timeoutTimestamp = sysT + defaultTimeoutDuration // Wait default timeout duration
				break
			} else if p.Timestamp <= sysT {
				// Elapsed proposal
				//logger.Debugf("TimestampScheduler: Elapsed proposal (opId = %v; timestamp = %v)", p.Op.Id, p.Timestamp)
				s.outCh <- s.Queue.Pop() // may block
			} else {
				// Non-elapsed proposal
				//logger.Debugf("TimestampScheduler: Non-elapsed proposal (opId = %v; timestamp = %v)", p.Op.Id, p.Timestamp)
				timeoutTimestamp = p.Timestamp // Wait for next proposal to elapse
				break
			}
		}

		s.queueLock.Unlock()


		// Wait until next proposal elapses
		WaitLoop:
			for {
				timeout := timeoutTimestamp - time.Now().UnixNano()
				if timeout < 0 {
					//logger.Debugf("TimestampScheduler: Calculated negative timeout; breaking WaitLoop")
					break WaitLoop
				}

				//logger.Debugf("TimestampScheduler: Waiting for next proposal; timeout = %v", timeout)
				resetTimer(timer, time.Duration(timeout))
				select {
				case <-timer.C:
					//logger.Debugf("TimestampScheduler: Timeout event; breaking WaitLoop")
					break WaitLoop
				case <-s.modifiedFlag:
					// A proposal is guaranteed to be in the queue if s.modifiedFlag is set

					s.queueLock.Lock()
					//logger.Debugf("TimestampScheduler: modifiedFlag is set")

					//logger.Debugf("TimestampScheduler: Clearing modifiedFlag (%v)", s.modifiedFlag.clearFlag())
					s.modifiedFlag.clearFlag()

					if newTimeoutTimestamp := s.Queue.Peek().Timestamp; newTimeoutTimestamp < timeoutTimestamp {
						timeoutTimestamp = newTimeoutTimestamp
					}

					s.queueLock.Unlock()
				}
			}
	}
}

// Receives a proposal
func (s *TimestampScheduler) addProposal(p *ClientProposal) {
	//logger.Debugf("TimestampScheduler: Received opId = (%s) timestamp = %d", p.Op.Id, p.Timestamp)

	s.queueLock.Lock()
	//logger.Debugf("TimestampScheduler: Scheduling opId = (%s) timestamp = %d", p.Op.Id, p.Timestamp)

	s.Queue.Push(p)
	s.modifiedFlag.setFlag()

	s.queueLock.Unlock()
}

func resetTimer(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}

	t.Reset(d)
}


type SynchronizedFlag chan bool

func NewSynchronizedFlag() SynchronizedFlag {
	return make(chan bool, 1)
}

func (ch SynchronizedFlag) setFlag() bool {
	select	{
	case ch <- true:
		return false
	default:
		return true
	}
}

func (ch SynchronizedFlag) clearFlag() bool {
	select {
	case _, ok := <-ch:
		if ok {
			return true
		} else {
			panic("Channel closed")
		}
	default:
		return false
	}
}
