package common

import (
	"container/heap"
)

type PriorityQueueInt64 struct {
	mh MinHeapInt64
}

func NewPriorityQueueInt64() *PriorityQueueInt64 {
	q := &PriorityQueueInt64{
		mh: make(MinHeapInt64, 0),
	}
	heap.Init(&(q.mh))
	return q
}

func (q *PriorityQueueInt64) Peek() (int64, bool) {
	return q.mh.Peek()
}

func (q *PriorityQueueInt64) Pop() int64 {
	return heap.Pop(&(q.mh)).(int64)
}

func (q *PriorityQueueInt64) Push(p int64) {
	heap.Push(&(q.mh), p)
}

func (q *PriorityQueueInt64) Len() int {
	return q.mh.Len()
}

func (q *PriorityQueueInt64) Cap() int {
	return q.mh.Cap()
}

// A min-heap implements heap.Interface
type MinHeapInt64 []int64

func (h MinHeapInt64) Len() int {
	return len(h)
}

func (h MinHeapInt64) Cap() int {
	return cap(h)
}

func (h MinHeapInt64) Less(i, j int) bool {
	if h[i] < h[j] {
		return true
	}
	return false
}

func (h MinHeapInt64) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeapInt64) Push(p interface{}) {
	*h = append(*h, p.(int64))
}

func (h *MinHeapInt64) Pop() interface{} {
	old := *h
	n := len(old)
	p := old[n-1]
	*h = old[0 : n-1]
	return p
}

func (h MinHeapInt64) Peek() (int64, bool) {
	if len(h) > 0 {
		return h[0], true
	}
	return 0, false
}
