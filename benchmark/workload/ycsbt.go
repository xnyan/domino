package workload

import (
	"strconv"
)

// Acknowledgement: this implementation derives from the YCSB+T workload.
// Currently, it only supports a fixed number of read-modify-write operations per transaction.
// TODO implements the full YCSB+T workload, like varying the probability of reads/writes.
type YcsbtWorkload struct {
	*AbstractWorkload

	read_num_per_txn  int
	write_num_per_txn int
}

func NewYcsbtWorkload(
	workload *AbstractWorkload,
	ycsbt_read_num_per_txn int,
	ycsbt_write_num_per_txn int,
) *YcsbtWorkload {
	ycsbt := &YcsbtWorkload{
		AbstractWorkload:  workload,
		read_num_per_txn:  ycsbt_read_num_per_txn,
		write_num_per_txn: ycsbt_write_num_per_txn,
	}

	return ycsbt
}

// Generates a txn. This function is currently not thread-safe
func (ycsbt *YcsbtWorkload) GenTxn() *Txn {
	ycsbt.txnCount++
	txnId := strconv.FormatInt(ycsbt.txnCount, 10)
	return ycsbt.buildTxn(txnId, ycsbt.read_num_per_txn, ycsbt.write_num_per_txn)
}

func (ycsbt *YcsbtWorkload) String() string {
	return "YcsbtWorkload"
}
