package dynamic

import (
	"time"
)

type TimeManager interface {
	GetCurrentTime() int64
	PrevT(t int64) int64
	NextT(t int64) int64
}

type realClockTm struct {
}

func NewRealClockTm() TimeManager {
	return &realClockTm{}
}

func (tm *realClockTm) GetCurrentTime() int64 {
	return time.Now().UnixNano()
}

func (tm *realClockTm) PrevT(t int64) int64 {
	return t - 1
}

func (tm *realClockTm) NextT(t int64) int64 {
	return t + 1
}
