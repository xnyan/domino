package latency

import (
	"sync"
	"time"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("LatencyPredictor")

type ProbeRet struct {
	Addr string
	Rt   time.Duration // rountrip time
}

// Thread safe latency manager
type syncLatManager struct {
	lm   *LatManager
	lock sync.Mutex
}

func (l *syncLatManager) AddProbeRet(pr *ProbeRet) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lm.AddLat(pr.Rt)
}

/*
func (l *syncLatManager) GetWindow95th() int64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.lm.GetWindow95th()
}
*/

func (l *syncLatManager) GetWindowPth(p float64) int64 {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.lm.GetWindowP(p)
}

type LatPredictor struct {
	dstTable map[string]*syncLatManager
}

func NewLatPredictor(dstList []string, windowLen time.Duration, windowSize int) *LatPredictor {
	p := &LatPredictor{
		dstTable: make(map[string]*syncLatManager),
	}
	for _, dst := range dstList {
		p.dstTable[dst] = &syncLatManager{
			lm: NewLatManager(windowLen, windowSize),
		}
	}
	return p
}

func (pm *LatPredictor) getLatMgr(addr string) *syncLatManager {
	l, ok := pm.dstTable[addr]
	if !ok {
		logger.Fatalf("There is no latency manager for addr = %s", addr)
	}
	return l
}

func (pm *LatPredictor) AddProbeRet(pr *ProbeRet) {
	l := pm.getLatMgr(pr.Addr)
	l.AddProbeRet(pr)
}

// Returns the predicted roundtrip latency in ms
func (pm *LatPredictor) PredictLat(addr string, p float64) int64 {
	l := pm.getLatMgr(addr)
	//return l.GetWindow95th()
	return l.GetWindowPth(p)
}
