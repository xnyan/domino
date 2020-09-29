package config

import (
	"bufio"
	"os"
	"sort"
	"strings"

	"github.com/op/go-logging"
	//"domino/common"
)

var logger = logging.MustGetLogger("DynamicCommon")

type ReplicaInfo struct {
	RId           string // replica id
	DcId          string // datacenter id
	Ip            string // Ip
	Port          string // Port
	IsPaxosLeader bool   // Is leader of a Paxos Shard
	IsFpLeader    bool   // Is leader of an Fast Paxos Shard
}

func (r *ReplicaInfo) GetNetAddr() string {
	return r.Ip + ":" + r.Port
}

func LoadReplicaInfo(filePath string) map[string]*ReplicaInfo {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	repDir := make(map[string]*ReplicaInfo) // replica id --> replica info
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skips invalid lines
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		fieldList := strings.Fields(line)
		rInfo := &ReplicaInfo{
			RId:           fieldList[0],
			DcId:          fieldList[1],
			Ip:            fieldList[2],
			Port:          fieldList[3],
			IsPaxosLeader: false,
			IsFpLeader:    false,
		}
		if fieldList[4] == "L" {
			rInfo.IsPaxosLeader = true
		}
		if fieldList[5] == "L" {
			rInfo.IsFpLeader = true
		}
		repDir[rInfo.RId] = rInfo
	}
	return repDir
}

func GenShardInfo(repDir map[string]*ReplicaInfo) (int, map[int32]string, map[int32]string) {
	rIdList := make([]string, 0, len(repDir))
	for rId, _ := range repDir {
		rIdList = append(rIdList, rId)
	}
	sort.Strings(rIdList)

	shardNum := 0
	pShardLeaderMap := make(map[int32]string)
	for _, rId := range rIdList {
		rInfo := repDir[rId]
		if rInfo.IsPaxosLeader {
			pShardLeaderMap[int32(shardNum)] = rId
			shardNum++
		}
	}

	fpShardLeaderMap := make(map[int32]string)
	for _, rId := range rIdList {
		rInfo := repDir[rId]
		if rInfo.IsFpLeader {
			fpShardLeaderMap[int32(shardNum)] = rId
			shardNum++
		}
	}

	return shardNum, pShardLeaderMap, fpShardLeaderMap
}

func GenShardAddrInfo(repDir map[string]*ReplicaInfo) (map[int32]string, map[int32]string) {
	rIdList := make([]string, 0, len(repDir))
	for rId, _ := range repDir {
		rIdList = append(rIdList, rId)
	}
	sort.Strings(rIdList)

	shardNum := 0
	pShardLeaderMap := make(map[int32]string)
	for _, rId := range rIdList {
		rInfo := repDir[rId]
		if rInfo.IsPaxosLeader {
			pShardLeaderMap[int32(shardNum)] = rInfo.GetNetAddr()
			shardNum++
		}
	}

	fpShardLeaderMap := make(map[int32]string)
	for _, rId := range rIdList {
		rInfo := repDir[rId]
		if rInfo.IsFpLeader {
			fpShardLeaderMap[int32(shardNum)] = rInfo.GetNetAddr()
			shardNum++
		}
	}

	return pShardLeaderMap, fpShardLeaderMap
}

/*
func GetShardInfo(p *common.Properties) (int, map[int32]string, map[int32]string) {
	shardNum := 0

	pShardLeaderMap := make(map[int32]string)
	pShardReplicaList := p.GetStrList(FLAG_DYNAMIC_PAXOS_SHARD_REPLICA_LIST, ",")
	for _, rId := range pShardReplicaList {
		//addr := p.GetStr(rId)
		//pShardLeaderMap[int32(shardNum)] = addr
		pShardLeaderMap[int32(shardNum)] = rId
		shardNum++
	}

	fpShardLeaderMap := make(map[int32]string)
	fpShardReplicaList := p.GetStrList(FLAG_DYNAMIC_FP_SHARD_REPLICA_LIST, ",")
	for _, rId := range fpShardReplicaList {
		//addr := p.GetStr(rId)
		//fpShardLeaderMap[int32(shardNum)] = addr
		fpShardLeaderMap[int32(shardNum)] = rId
		shardNum++
	}
	return shardNum, pShardLeaderMap, fpShardLeaderMap
}

func GetShardAddrInfo(p *common.Properties) (int, map[int32]string, map[int32]string) {
	shardNum := 0

	pShardLeaderMap := make(map[int32]string)
	pShardReplicaList := p.GetStrList(FLAG_DYNAMIC_PAXOS_SHARD_REPLICA_LIST, ",")
	for _, rId := range pShardReplicaList {
		addr := p.GetStr(rId)
		pShardLeaderMap[int32(shardNum)] = addr
		shardNum++
	}

	fpShardLeaderMap := make(map[int32]string)
	fpShardReplicaList := p.GetStrList(FLAG_DYNAMIC_FP_SHARD_REPLICA_LIST, ",")
	for _, rId := range fpShardReplicaList {
		addr := p.GetStr(rId)
		fpShardLeaderMap[int32(shardNum)] = addr
		shardNum++
	}
	return shardNum, pShardLeaderMap, fpShardLeaderMap
}
*/
