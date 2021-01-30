package example

import (
	"time"

	"github.com/Ak-Army/timer"
)

type tickerExample struct {
	tick1 int
	tick2 int
	tick3 int
}

func (t *tickerExample) timerExample() {
	t1 := timer.NewTicker("tickerExample1", time.Second)
	t2 := timer.NewTicker("tickerExample2", 3*time.Second)
	t3 := timer.NewTicker("tickerExample3", 5*time.Second)
	defer t1.Stop()
	defer t2.Stop()
	defer t3.Stop()
	for {
		select {
		case <-t1.C():
			t.tick1++
		case <-t2.C():
			t.tick2++
		case <-t3.C():
			t.tick3++
			return
		}
	}
}
