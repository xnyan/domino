package common

import (
	"time"
)

type Clock interface {
	NextTime(t int64) int64
	PrevTime(t int64) int64
	GetClockTime() int64
}

////////////////////////////////////
//System clock in ns
type SysNanoClock struct {
}

func NewSysNanoClock() *SysNanoClock {
	return &SysNanoClock{}
}

func (c *SysNanoClock) NextTime(t int64) int64 {
	return t + 1
}

func (c *SysNanoClock) PrevTime(t int64) int64 {
	return t - 1
}

func (c *SysNanoClock) GetClockTime() int64 {
	return time.Now().UnixNano()
}

////////////////////////////////////
//System clock in us
type SysMicroClock struct {
}

func NewSysMicroClock() *SysMicroClock {
	return &SysMicroClock{}
}

func (c *SysMicroClock) NextTime(t int64) int64 {
	return t + 1000
}

func (c *SysMicroClock) PrevTime(t int64) int64 {
	return t - 1000
}

func (c *SysMicroClock) GetClockTime() int64 {
	return time.Now().UnixNano() / 1000 * 1000
}

////////////////////////////////////
//System clock in ms
type SysMilliClock struct {
}

func NewSysMilliClock() *SysMilliClock {
	return &SysMilliClock{}
}

func (c *SysMilliClock) NextTime(t int64) int64 {
	return t + 1000000
}

func (c *SysMilliClock) PrevTime(t int64) int64 {
	return t - 1000000
}

func (c *SysMilliClock) GetClockTime() int64 {
	return time.Now().UnixNano() / 1000000 * 1000000
}

////////////////////////////////////
//Fixed clock, not thread-safe
type FixedStepClock struct {
	base int64
	step int64
}

func NewFixedStepClock(base, step int64) *FixedStepClock {
	return &FixedStepClock{base: base, step: step}
}

func (c *FixedStepClock) NextTime(t int64) int64 {
	return t + c.step
}

func (c *FixedStepClock) PrevTime(t int64) int64 {
	return t - c.step
}

func (c *FixedStepClock) GetClockTime() int64 {
	c.base += c.step
	return c.base
}
