package timer

import (
	"runtime"
	"sync"
	"time"

	"bou.ke/monkey"
)

type MockedTickers struct {
	sync.Mutex
	tickers    map[string]Ticker
	timerPatch *monkey.PatchGuard
}

func NewMockedTickers() *MockedTickers {
	m := &MockedTickers{
		tickers: make(map[string]Ticker),
	}
	m.timerPatch = monkey.Patch(NewTicker, func(name string, d time.Duration) Ticker {
		m.Lock()
		defer m.Unlock()
		if _, ok := m.tickers[name]; !ok {
			m.tickers[name] = newMockTicker(name, d)
		}

		return m.tickers[name]
	})
	return m
}

func (m *MockedTickers) Reset() {
	m.Lock()
	defer m.Unlock()
	m.tickers = make(map[string]Ticker)
}

func (m *MockedTickers) GetTicker(name string) *mockTicker {
	for {
		m.Lock()
		if t, ok := m.tickers[name]; ok {
			m.Unlock()
			runtime.Gosched()
			return t.(*mockTicker)
		}
		m.Unlock()
		runtime.Gosched()
		time.Sleep(1 * time.Millisecond)
	}
}

func (m *MockedTickers) UnPatch() {
	m.timerPatch.Unpatch()
}

func newMockTicker(name string, _ time.Duration) Ticker {
	return &mockTicker{
		name:     name,
		c:        make(chan time.Time),
		stopChan: make(chan struct{}),
	}
}

type mockTicker struct {
	name     string
	c        chan time.Time
	stopChan chan struct{}
}

func (mt mockTicker) C() <-chan time.Time {
	return mt.c
}

func (mt *mockTicker) Stop() {
	close(mt.stopChan)
}

func (mt *mockTicker) Tick() {
	mt.c <- time.Now()
	runtime.Gosched()
}

func (mt *mockTicker) TickAndWaitForStop() {
	mt.c <- time.Now()
	<-mt.stopChan
}

func (mt *mockTicker) WaitForStop() {
	<-mt.stopChan
}
