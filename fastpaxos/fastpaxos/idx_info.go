package fastpaxos

import (
	"strconv"
)

type IdxInfo struct {
	OpVoteTable  map[string]int
	TotalVoteNum int
	FastOpId     string // fast-path op
}

func NewIdxInfo() *IdxInfo {
	return &IdxInfo{
		OpVoteTable:  make(map[string]int),
		TotalVoteNum: 0,
		FastOpId:     "",
	}
}

func (idxInfo *IdxInfo) Vote(opId string) int {
	if _, exists := idxInfo.OpVoteTable[opId]; !exists {
		idxInfo.OpVoteTable[opId] = 0
	}
	idxInfo.OpVoteTable[opId]++
	idxInfo.TotalVoteNum++
	return idxInfo.OpVoteTable[opId]
}

func (idxInfo *IdxInfo) GetTotalVoteNum() int {
	return idxInfo.TotalVoteNum
}

func (idxInfo *IdxInfo) SetFastOpId(opId string) {
	idxInfo.FastOpId = opId
}

func (idxInfo *IdxInfo) GetFastOpId() string {
	return idxInfo.FastOpId
}

func (idxInfo *IdxInfo) isFast() bool {
	return idxInfo.FastOpId != ""
}

func (idxInfo *IdxInfo) getOpId(count int) (string, bool) {
	for opId, vN := range idxInfo.OpVoteTable {
		if vN >= count {
			return opId, true
		}
	}
	return "", false
}

func (idxInfo *IdxInfo) String() string {
	ret := ""
	for opId, vN := range idxInfo.OpVoteTable {
		ret += opId + "->" + strconv.Itoa(vN) + ";"
	}
	ret += "(" + strconv.Itoa(idxInfo.TotalVoteNum) + "," + idxInfo.FastOpId + ")"
	return ret
}
