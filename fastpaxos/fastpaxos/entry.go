package fastpaxos

import (
	"strconv"

	"domino/fastpaxos/rpc"
)

//// Log entry struct

const (
	ENTRY_INIT          = 0
	ENTRY_FAST_ACCEPTED = 1
	ENTRY_SLOW_ACCEPTED = 2
	ENTRY_COMMITTED     = 3
)

type Entry struct {
	op     *rpc.Operation
	status int
	timestamp int64 // added by @skoya76
}

func (entry *Entry) SetOp(op *rpc.Operation) {
	entry.op = op
}

func (entry *Entry) GetOp() *rpc.Operation {
	return entry.op
}

func (entry *Entry) SetFastAccepted() {
	entry.status = ENTRY_FAST_ACCEPTED
}

func (entry *Entry) SetSlowAccepted() {
	entry.status = ENTRY_SLOW_ACCEPTED
}

func (entry *Entry) IsSlowAccepted() bool {
	return entry.status == ENTRY_SLOW_ACCEPTED

}

func (entry *Entry) SetCommitted() {
	entry.status = ENTRY_COMMITTED
}

func (entry *Entry) IsCommitted() bool {
	return entry.status == ENTRY_COMMITTED
}

func (entry *Entry) String() string {
	return strconv.Itoa(entry.status) + ", " +
		entry.op.Id + ", " +
		entry.op.Type + ", " +
		entry.op.Key + ", " +
		entry.op.Val
}

func (entry *Entry) SetStartDuration(t int64){
	entry.timestamp = t
}

func (entry *Entry) GetStartDuration() int64{
	return entry.timestamp
}