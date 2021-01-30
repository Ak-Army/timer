package timer

import "time"

type Ticker interface {
	Stop()
	C() <-chan time.Time
}

type namedTicker struct {
	*time.Ticker
	name string
}

//go:noinline
func NewTicker(name string, d time.Duration) Ticker {
	return &namedTicker{
		Ticker: time.NewTicker(d),
		name:   name,
	}
}

func (t namedTicker) C() <-chan time.Time {
	return t.Ticker.C
}

func Tick(name string, d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return NewTicker(name, d).C()
}
