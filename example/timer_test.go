package example

import (
	"testing"

	"github.com/Ak-Army/timer"
)

func TestTimer(t *testing.T) {
	te := timerExample{}
	mockTimer := timer.NewMockedTimer()
	go func() {
		te.timerExample()
	}()
	timerExample1 := mockTimer.GetTimer("timerExample1")
	timerExample2 := mockTimer.GetTimer("timerExample2")
	timerExample3 := mockTimer.GetTimer("timerExample3")
	timerExample1.TickAndWaitForReset()
	if te.tick1 != 1 {
		t.Error("tick1 is wrong")
	}
	timerExample2.TickAndWaitForReset()
	if te.tick1 != 1 {
		t.Error("tick1 is wrong")
	}
	if te.tick2 != 1 {
		t.Error("tick2 is wrong")
	}
	timerExample2.TickAndWaitForReset()
	if te.tick1 != 1 {
		t.Error("tick1 is wrong")
	}
	if te.tick2 != 2 {
		t.Error("tick2 is wrong")
	}
	timerExample3.TickAndWaitForStop()
	if te.tick1 != 1 {
		t.Error("tick1 is wrong")
	}
	if te.tick2 != 2 {
		t.Error("tick2 is wrong")
	}
	if te.tick3 != 1 {
		t.Error("tick3 is wrong")
	}
}
