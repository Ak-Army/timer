package example

import (
	"testing"

	"github.com/Ak-Army/timer"
)

func TestTicker(t *testing.T) {
	te := tickerExample{}
	mockTicker := timer.NewMockedTickers()
	go func() {
		te.timerExample()
	}()
	timerExample1 := mockTicker.GetTicker("tickerExample1")
	timerExample2 := mockTicker.GetTicker("tickerExample2")
	timerExample3 := mockTicker.GetTicker("tickerExample3")
	timerExample1.Tick()
	if te.tick1 != 1 {
		t.Error("tick1 is wrong")
	}
	timerExample2.Tick()
	if te.tick1 != 1 {
		t.Error("tick1 is wrong")
	}
	if te.tick2 != 1 {
		t.Error("tick2 is wrong")
	}
	timerExample2.Tick()
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
