package common

import (
	"sync"
)

type Future struct {
	value interface{}
	cv    *sync.Cond
}

func NewFuture() *Future {
	f := &Future{value: nil}
	l := sync.Mutex{}
	f.cv = sync.NewCond(&l)
	return f
}

// Blocks if the value is not available yet
func (f *Future) GetValue() interface{} {
	f.cv.L.Lock()
	for f.value == nil {
		f.cv.Wait()
	}
	f.cv.L.Unlock()
	return f.value
}

// Sets up the value and wakes up all of the listeners
// The value CANNOT be nil, otherwise the listeners will be in deadlock
func (f *Future) SetValue(v interface{}) {
	f.cv.L.Lock()
	f.value = v
	f.cv.L.Unlock()

	f.cv.Broadcast()
}
