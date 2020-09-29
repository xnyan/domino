package fastpaxos

import (
	"testing"

	"domino/rpc"
)

func TestPriority(t *testing.T) {
	q := NewPriorityQueue()

	q.Push(&ClientProposal{
		Op:        &rpc.Operation{Id: "1"},
		Timestamp: 1,
		ClientId:  "1",
	})

	q.Push(&ClientProposal{
		Op:        &rpc.Operation{Id: "3"},
		Timestamp: 1,
		ClientId:  "3",
	})

	q.Push(&ClientProposal{
		Op:        &rpc.Operation{Id: "2"},
		Timestamp: 1,
		ClientId:  "2",
	})

	q.Push(&ClientProposal{
		Op:        &rpc.Operation{Id: "0"},
		Timestamp: 0,
		ClientId:  "2",
	})

	q.Push(&ClientProposal{
		Op:        &rpc.Operation{Id: "6"},
		Timestamp: 6,
		ClientId:  "1",
	})

	q.Push(&ClientProposal{
		Op:        &rpc.Operation{Id: "5"},
		Timestamp: 5,
		ClientId:  "1",
	})

	q.Push(&ClientProposal{
		Op:        &rpc.Operation{Id: "4"},
		Timestamp: 4,
		ClientId:  "1",
	})

	p := q.Peek()
	if p.Op.Id != "0" {
		t.Errorf("Expects 0 but %s", p.Op.Id)
	}

	p = q.Pop()
	if p.Op.Id != "0" {
		t.Errorf("Expects 0 but %s", p.Op.Id)
	}

	p = q.Peek()
	if p.Op.Id != "1" {
		t.Errorf("Expects 1 but %s", p.Op.Id)
	}

	p = q.Pop()
	if p.Op.Id != "1" {
		t.Errorf("Expects 1 but %s", p.Op.Id)
	}

	p = q.Peek()
	if p.Op.Id != "2" {
		t.Errorf("Expects 2 but %s", p.Op.Id)
	}

	p = q.Pop()
	if p.Op.Id != "2" {
		t.Errorf("Expects 2 but %s", p.Op.Id)
	}

	p = q.Peek()
	if p.Op.Id != "3" {
		t.Errorf("Expects 3 but %s", p.Op.Id)
	}

	p = q.Pop()
	if p.Op.Id != "3" {
		t.Errorf("Expects 3 but %s", p.Op.Id)
	}

	p = q.Peek()
	if p.Op.Id != "4" {
		t.Errorf("Expects 4 but %s", p.Op.Id)
	}

	p = q.Pop()
	if p.Op.Id != "4" {
		t.Errorf("Expects 4 but %s", p.Op.Id)
	}

	p = q.Peek()
	if p.Op.Id != "5" {
		t.Errorf("Expects 5 but %s", p.Op.Id)
	}

	p = q.Pop()
	if p.Op.Id != "5" {
		t.Errorf("Expects 5 but %s", p.Op.Id)
	}

	p = q.Peek()
	if p.Op.Id != "6" {
		t.Errorf("Expects 6 but %s", p.Op.Id)
	}

	p = q.Pop()
	if p.Op.Id != "6" {
		t.Errorf("Expects 6 but %s", p.Op.Id)
	}

	p = q.Peek()
	if p != nil {
		t.Errorf("Expects nil but %v", p)
	}
}
