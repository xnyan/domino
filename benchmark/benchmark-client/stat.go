package main

import (
	"time"
)

// Txn execution statistics
type TxnStat struct {
	TxnStatId      string
	ReadKeyList    []string
	WriteKeyValMap map[string]string
	ExecHistory    []*TxnExecutionStat // A list of execution statuses for the transaction.
}

func NewTxnStat(
	txnStatId string,
	readKeyList []string,
	writeKeyValMap map[string]string,
) *TxnStat {
	return &TxnStat{
		TxnStatId:      txnStatId,
		ReadKeyList:    readKeyList,
		WriteKeyValMap: writeKeyValMap,
		ExecHistory:    make([]*TxnExecutionStat, 0),
	}
}

type TxnExecutionStat struct {
	TxnId     string // txn execution id. A transaction needs to change its txn id for retrying
	BeginTime *time.Time
	EndTime   *time.Time
	IsUseFp   bool
	IsAccept  bool
	IsFast    bool
}

func (execStat *TxnExecutionStat) ExecutionDuration() time.Duration {
	return execStat.EndTime.Sub(*(execStat.BeginTime))
}
