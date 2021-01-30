package timer

import (
	"time"
)

type Timer interface {
	Stop() bool
	Reset(d time.Duration) bool
	C() <-chan time.Time
}

type namedTimer struct {
	*time.Timer
	name string
}

//go:noinline
func NewTimer(name string, d time.Duration) Timer {
	return &namedTimer{
		Timer: time.NewTimer(d),
		name:  name,
	}
}

func (t namedTimer) C() <-chan time.Time {
	return t.Timer.C
}

func After(name string, d time.Duration) <-chan time.Time {
	return NewTimer(name, d).C()
}
