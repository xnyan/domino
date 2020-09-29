package fastpaxos

import (
	"container/heap"
)

type PriorityQueue struct {
	mh MinHeap
}

func NewPriorityQueue() *PriorityQueue {
	q := &PriorityQueue{
		mh: make(MinHeap, 0),
	}
	heap.Init(&(q.mh))
	return q
}

func (q *PriorityQueue) Peek() *ClientProposal {
	p := q.mh.Peek()
	if p == nil {
		return nil
	}
	return p.(*ClientProposal)
}

func (q *PriorityQueue) Pop() *ClientProposal {
	return heap.Pop(&(q.mh)).(*ClientProposal)
}

func (q *PriorityQueue) Push(p *ClientProposal) {
	heap.Push(&(q.mh), p)
}

// A min-heap implements heap.Interface
type MinHeap []*ClientProposal

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	// Pop to give the element with the oldest (smallest) timestamp
	if h[i].Timestamp < h[j].Timestamp {
		return true
	} else if h[i].Timestamp == h[j].Timestamp {
		return h[i].ClientId < h[j].ClientId
	} else {
		return false
	}
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(p interface{}) {
	*h = append(*h, p.(*ClientProposal))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	p := old[n-1]
	*h = old[0 : n-1]
	return p
}

func (h MinHeap) Peek() interface{} {
	if len(h) > 0 {
		return h[0]
	}
	return nil
}
