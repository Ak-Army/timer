package timer

import (
	"runtime"
	"sync"
	"time"

	"bou.ke/monkey"
)

type MockedTimers struct {
	sync.Mutex
	timers     map[string]Timer
	timerPatch *monkey.PatchGuard
}

func NewMockedTimer() *MockedTimers {
	m := &MockedTimers{
		timers: make(map[string]Timer),
	}
	m.timerPatch = monkey.Patch(NewTimer, func(name string, d time.Duration) Timer {
		m.Lock()
		defer m.Unlock()
		if _, ok := m.timers[name]; !ok {
			m.timers[name] = newMockTimer(name, d)
		}

		return m.timers[name]
	})
	return m
}

func (m *MockedTimers) Reset() {
	m.Lock()
	defer m.Unlock()
	m.timers = make(map[string]Timer)
}

func (m *MockedTimers) GetTimer(name string) *mockedTimer {
	for {
		m.Lock()
		if t, ok := m.timers[name]; ok {
			m.Unlock()
			runtime.Gosched()
			return t.(*mockedTimer)
		}
		m.Unlock()
		runtime.Gosched()
		time.Sleep(1 * time.Millisecond)
	}
}

func (m *MockedTimers) UnPatch() {
	m.timerPatch.Unpatch()
}

func newMockTimer(name string, _ time.Duration) Timer {
	return &mockedTimer{
		name:      name,
		c:         make(chan time.Time),
		resetChan: make(chan time.Duration),
		stopChan:  make(chan struct{}),
	}
}

type mockedTimer struct {
	name      string
	c         chan time.Time
	resetChan chan time.Duration
	stopChan  chan struct{}
}

func (mt mockedTimer) C() <-chan time.Time {
	return mt.c
}

func (mt *mockedTimer) Reset(d time.Duration) bool {
	mt.resetChan <- d
	return false
}

func (mt *mockedTimer) Stop() bool {
	close(mt.resetChan)
	close(mt.stopChan)
	return true
}

func (mt *mockedTimer) Tick() {
	mt.c <- time.Now()
}

func (mt *mockedTimer) TickAndWaitForReset() {
	mt.c <- time.Now()
	mt.WaitForReset()
}

func (mt *mockedTimer) TickAndWaitForStop() {
	mt.c <- time.Now()
	mt.WaitForStop()
}

func (mt *mockedTimer) TickAndForResetMultiple(i int) {
	for ; i > 0; i-- {
		mt.TickAndWaitForReset()
	}
}

func (mt *mockedTimer) WaitForReset() time.Duration {
	return <-mt.resetChan
}

func (mt *mockedTimer) WaitForStop() {
	<-mt.stopChan
}
