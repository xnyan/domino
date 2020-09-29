package common

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	q.Push(1)
	// Testing iteration
	itr := q.Iterator()
	for itr.HasNext() {
		n, ok := itr.Next()
		if !ok {
			t.Errorf("Expect true but %t", ok)
		}
		if n.(int) != 1 {
			t.Errorf("Expect 1 but %d", n)
		}
	}

	q.Push(2)
	q.Push(3)
	itr = q.Iterator()
	n, _ := itr.Next()
	if n.(int) != 1 {
		t.Errorf("Expect 1 but %d", n)
	}
	n, _ = itr.Next()
	if n.(int) != 2 {
		t.Errorf("Expect 2 but %d", n)
	}

	n, _ = itr.Next()
	if n.(int) != 3 {
		t.Errorf("Expect 3 but %d", n)
	}

	_, ok := itr.Next()
	if ok {
		t.Errorf("Expect faslt but %t", ok)
	}
}
