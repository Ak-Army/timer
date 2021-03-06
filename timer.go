package timer

import (
	"time"
)

type Timer interface {
	Stop() bool
	SafeStop()
	Reset(d time.Duration) bool
	SafeReset(d time.Duration)
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

func (t namedTimer) SafeStop() {
	t.Timer.Stop()
	select {
	case <-t.Timer.C:
	default:
	}
}

func (t namedTimer) SafeReset(d time.Duration) {
	t.SafeStop()
	t.Timer.Reset(d)
}

func After(name string, d time.Duration) <-chan time.Time {
	return NewTimer(name, d).C()
}
