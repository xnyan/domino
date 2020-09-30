package node

import (
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Probing result information
// Using a short naming struct to save logging storage space when using gob
type L struct {
	S int64 // client sending time
	E int64 // client receiving time
	C int64 // server clock time
}

type Client struct {
	cm *CommManager

	dc      string
	dcTable map[string]string // net addr --> dc id

	addrList []string
	logTable map[string]Log     // net addr --> logging instance
	chTable  map[string]chan *L // net addr --> logging input channle
}

func NewClient(dc string, dcTable map[string]string, logDir string) *Client {
	c := &Client{
		cm:       NewCommManager(),
		dc:       dc,
		dcTable:  dcTable,
		addrList: make([]string, 0),
		logTable: make(map[string]Log),
		chTable:  make(map[string]chan *L),
	}

	for dstDc, addr := range dcTable {
		if c.dc != dstDc {
			c.addrList = append(c.addrList, addr)
			c.chTable[addr] = make(chan *L, 1024*1024*16)
			if logDir == "" {
				logDir = "."
			}
			fileName := logDir + "/" + c.dc + "-" + dstDc + ".log"
			c.logTable[addr] = NewLog(fileName)
		}
	}

	if len(c.addrList) == 0 {
		logger.Fatalf("Configuration error: no target datacenter")
	}

	return c
}

func (c *Client) StartProbing(inv, duration time.Duration) {
	// Builds network connection to servers
	c.InitConn(c.addrList) // NOTE: This will async build the connections

	// Starts logging threads
	var wg sync.WaitGroup
	c.startLogging(c.addrList, wg)

	// Waits for the network connections to be ready
	c.waitConnReady(c.addrList)

	// Starts probing
	c.doStartProbing(c.addrList, inv, duration) // blocking

	// Stops logging threads
	c.closeLogging(c.addrList, wg) // blocking
}

// Inits tcp connections to the given servers
func (c *Client) InitConn(addrList []string) {
	var wg sync.WaitGroup
	for _, addr := range addrList {
		wg.Add(1)
		go func(addr string) {
			c.cm.BuildConnection(addr)
			wg.Done()
		}(addr)
	}
	wg.Wait()
}

func (c *Client) waitConnReady(addrList []string) {
	var wg sync.WaitGroup
	for _, addr := range addrList {
		wg.Add(1)
		go func(addr string) {
			logger.Debugf("Testing connection to addr = %s", addr)
			c.Probe(addr)
			logger.Debugf("Connection to addr = %s is OK", addr)
			wg.Done()
		}(addr)
	}
	wg.Wait()
	logger.Debugf("Connections are ready")
}

func (c *Client) doStartProbing(addrList []string, inv, duration time.Duration) {
	sent, rev := 0, 0 // debuging
	var revL sync.Mutex
	logger.Debugf("Starting probing")

	var wg sync.WaitGroup
	run := true
	timer := time.NewTimer(1)
	start := time.Now()
	for run {
		select {
		case <-timer.C:
			for _, addr := range addrList {
				wg.Add(1)
				sent++
				go func(addr string) {
					start := time.Now().UnixNano()
					sClock := c.Probe(addr) // blocking

					revL.Lock()
					rev++
					revL.Unlock()

					end := time.Now().UnixNano()
					ch := c.getCh(addr)
					ch <- &L{
						S: start,
						E: end,
						C: sClock,
					}
					wg.Done()
				}(addr)
			}
			elapse := time.Since(start)
			if elapse >= duration {
				run = false
			} else {
				timer.Reset(inv)
			}
		}
	}
	wg.Wait()

	logger.Debugf("Stopped probing total sent = %d, received = %d", sent, rev)
}

// Blocking
func (c *Client) Probe(addr string) int64 {
	reply := c.sendProbeReq(addr)
	return reply.ServerClock
}

// Blocking
func (c *Client) sendProbeReq(addr string) *ProbeReply {
	req := &ProbeReq{}
	rpcStub := c.cm.NewRpcStub(addr)
	reply, err := rpcStub.Probe(context.Background(), req, grpc.WaitForReady(true))
	if err != nil {
		logger.Errorf("Error: %v", err)
		logger.Fatalf("Fails sending probe request to addr = %s", addr)
	}
	return reply
}

////Logging
func (c *Client) startLogging(addrList []string, wg sync.WaitGroup) {
	logger.Debugf("Starting logging")

	var barrier sync.WaitGroup
	for _, addr := range addrList {
		wg.Add(1)
		barrier.Add(1)
		go func(addr string) {
			ch, log := c.getCh(addr), c.getLog(addr)
			barrier.Done()
			for l := range ch {
				b := LatInfoToByte(l)
				log.Write(b)
			}
			log.Flush()
			log.Close()
			wg.Done()
		}(addr)
	}
	barrier.Wait()

	logger.Debugf("Started logging")
}

func (c *Client) closeLogging(addrList []string, wg sync.WaitGroup) {
	logger.Debugf("Closing logging")

	for _, addr := range addrList {
		ch := c.getCh(addr)
		close(ch)
	}
	// Waits for all logging threads to flush data to disks
	wg.Wait()

	logger.Debugf("Closed logging")
}

////Helper functions
func (c *Client) getLog(addr string) Log {
	if log, ok := c.logTable[addr]; ok {
		return log
	} else {
		logger.Fatalf("Cannot find logger for addr = %s", addr)
	}
	return nil
}

func (c *Client) getCh(addr string) chan *L {
	if ch, ok := c.chTable[addr]; ok {
		return ch
	} else {
		logger.Fatalf("Cannot find logging channel for addr = %s", addr)
	}
	return nil
}
