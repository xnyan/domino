package clientlib

import (
	"math"
	"domino/dynamic/config"
)

func (c *Client) loadConfig(configFile, replicaFile string) {
	// Loading Replica Information
	replicaDir := config.LoadReplicaInfo(replicaFile)

	// Replica Information
	for _, rInfo := range replicaDir {
		c.replicaNum++
		addr := rInfo.GetNetAddr()
		if rInfo.IsFpLeader {
			c.leaderAddr = addr
		} else {
			c.followerAddrList = append(c.followerAddrList, addr)
		}
	}
	if c.replicaNum != len(c.followerAddrList)+1 {
		logger.Fatalf("Error: replicaNum = %d, but follower num = %d", c.replicaNum, len(c.followerAddrList))
	}
	f := (c.replicaNum - 1) / 2
	c.fastQuorum = int(math.Ceil((3.0*float64(f))/2.0)) + 1
}
