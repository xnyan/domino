package workload

import (
	"strconv"
)

// This workload always generates a transaction that has the same reads and writes
type OneTxnWorkload struct {
	*AbstractWorkload

	read_num  int
	write_num int
}

func NewOneTxnWorkload(
	workload *AbstractWorkload,
	one_txn_read_num int,
	one_txn_write_num int,
) *OneTxnWorkload {
	oneTxn := &OneTxnWorkload{
		AbstractWorkload: workload,
		read_num:         one_txn_read_num,
		write_num:        one_txn_write_num,
	}

	return oneTxn
}

// Generates a txn. This function is currently not thread-safe
func (oneTxn *OneTxnWorkload) GenTxn() *Txn {
	oneTxn.txnCount++
	txnId := strconv.FormatInt(oneTxn.txnCount, 10)

	txn := &Txn{
		TxnId:     txnId,
		ReadKeys:  make([]string, 0),
		WriteData: make(map[string]string),
	}

	// read keys
	for i := 0; i < oneTxn.read_num; i++ {
		txn.ReadKeys = append(txn.ReadKeys, oneTxn.KeyList[i])
	}
	// write keys
	for i := 0; i < oneTxn.write_num; i++ {
		//txn.WriteData[oneTxn.KeyList[i]] = oneTxn.KeyList[i]
		txn.WriteData[oneTxn.KeyList[i]] = oneTxn.getVal()
	}

	return txn
}
