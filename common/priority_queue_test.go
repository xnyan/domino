package common

import (
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	q := NewPriorityQueueInt64()
	if v, ok := q.Peek(); ok {
		t.Errorf("v %d, ok %t", v, ok)
	}
}
