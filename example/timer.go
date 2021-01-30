package example

import (
	"time"

	"github.com/Ak-Army/timer"
)

type timerExample struct {
	tick1 int
	tick2 int
	tick3 int
}

func (t *timerExample) timerExample() {
	t1 := timer.NewTimer("timerExample1", time.Second)
	t2 := timer.NewTimer("timerExample2", 3*time.Second)
	t3 := timer.NewTimer("timerExample3", 5*time.Second)
	defer t1.Stop()
	defer t2.Stop()
	defer t3.Stop()
	for {
		select {
		case <-t1.C():
			t.tick1++
			t1.Reset(time.Second)
		case <-t2.C():
			t.tick2++
			t2.Reset(3 * time.Second)
		case <-t3.C():
			t.tick3++
			return
		}
	}
}
