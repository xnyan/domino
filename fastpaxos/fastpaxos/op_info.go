package fastpaxos

import (
	"strconv"
)

type OpInfo struct {
	IdxVoteTable map[string]int // idx --> count
	TotalVoteNum int
	FastIdx      string
}

func NewOpInfo() *OpInfo {
	return &OpInfo{
		IdxVoteTable: make(map[string]int),
		TotalVoteNum: 0,
		FastIdx:      INVALID_IDX,
	}
}

func (opInfo *OpInfo) Vote(idx string) int {
	if _, exists := opInfo.IdxVoteTable[idx]; !exists {
		opInfo.IdxVoteTable[idx] = 0
	}
	opInfo.IdxVoteTable[idx]++
	opInfo.TotalVoteNum++
	return opInfo.IdxVoteTable[idx]
}

func (opInfo *OpInfo) GetTotalVoteNum() int {
	return opInfo.TotalVoteNum
}

func (opInfo *OpInfo) SetFastIdx(idx string) {
	opInfo.FastIdx = idx
}

func (opInfo *OpInfo) GetFastIdx() string {
	return opInfo.FastIdx
}

// This function can only be called after the fast path is determined to fail or succeed
func (opInfo *OpInfo) IsFast() bool {
	return opInfo.FastIdx != INVALID_IDX
}

// Returns true if vote number is not less than the given num
func (opInfo *OpInfo) IsVoteNum(num int) bool {
	if opInfo.TotalVoteNum >= num {
		return true
	}
	return false
}

func (opInfo *OpInfo) getIdx(count int) (string, bool) {
	for idx, vN := range opInfo.IdxVoteTable {
		if vN >= count {
			return idx, true
		}
	}
	return INVALID_IDX, false
}

func (opInfo *OpInfo) String() string {
	ret := ""
	for idx, vN := range opInfo.IdxVoteTable {
		ret += idx + "->" + strconv.Itoa(vN) + ";"
	}
	ret += "(" + strconv.Itoa(opInfo.TotalVoteNum) + "," + opInfo.FastIdx + ")"
	return ret
}
