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

func (t *namedTimer) C() <-chan time.Time {
	return t.Timer.C
}

func (t *namedTimer) SafeStop() {
	t.Timer.Stop()
	select {
	case <-t.Timer.C:
	default:
	}
}

func (t *namedTimer) SafeReset(d time.Duration) {
	t.SafeStop()
	t.Timer.Reset(d)
}

func After(name string, d time.Duration) <-chan time.Time {
	return NewTimer(name, d).C()
}

type afterFuncTimer struct {
	Timer
	fn     func()
	stopCh chan struct{}
	done   chan struct{}
}

func (t *afterFuncTimer) SafeStop() {
	select {
	case <-t.stopCh:
		return
	default:
	}
	close(t.stopCh)
	<-t.done
	t.Timer.SafeStop()
}

func (t *afterFuncTimer) SafeReset(d time.Duration) {
	t.SafeStop()
	t.Timer.Reset(d)
	t.start()
}

func (t *afterFuncTimer) start() {
	t.stopCh = make(chan struct{})
	t.done = make(chan struct{})
	go func() {
		select {
		case <-t.C():
			close(t.stopCh)
			t.fn()
		case <-t.stopCh:
		}
		close(t.done)
	}()
}

func AfterFunc(name string, d time.Duration, fn func()) Timer {
	timer := &afterFuncTimer{
		Timer: NewTimer(name, d).(*namedTimer),
		fn:    fn,
	}
	timer.start()

	return timer
}
