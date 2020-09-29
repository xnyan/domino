package dynamic

import (
	"time"
)

type LatInfo struct {
	addr string
	rt   time.Duration // rountrip time
}

type ReplicaLatInfo struct {
	addr    string                   // the src replica
	rtTable map[string]time.Duration // the rountrip time to the other replicas from the src
}

type LatReq struct {
	RetC chan bool
}

func (dp *DynamicPaxos) handleLatReq(req *LatReq) {
	req.RetC <- true
}
