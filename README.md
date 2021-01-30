# Timer

A mockable timer. Ability to test the timer/tick function in an easy way.


# Usage

Timer usage is as simple as the built-in timer.Timer:
```go
package main

import (
	"fmt"
	"time"

	"github.com/Ak-Army/timer"
)

func timerExample() {
	t := timer.NewTimer("timerExample", time.Second)
	endTime := timer.NewTimer("endTime", 5*time.Second)
	defer t.Stop()
	defer endTime.Stop()
	for {
		select {
		case <-t.C():
			fmt.Println("tick")
			t.Reset(time.Second)
		case <-endTime.C():
			fmt.Println("end")
			return
        }
    } 
}
```

On test, just create a NewMockedTimer. Start the function which will create a timer, than you can 
collect it from teh mockTimer instance with the GetTimer() function call:

```go
package main


import (
	"testing"

	"github.com/Ak-Army/timer"
)

func TestTimer(t *testing.T) {
	mockTimer := timer.NewMockedTimer()
	go func() {
		timerExample()
	}()
	timerExample := mockTimer.GetTimer("timerExample")
	endTime := mockTimer.GetTimer("endTime")
	timerExample.TickAndWaitForReset()
	timerExample.Tick()
	timerExample.WaitForReset()
	endTime.TickAndWaitForStop()
}
```
Ticker usage is simple too:
```go
package main

import (
	"fmt"
	"time"

	"github.com/Ak-Army/timer"
)

func tickerExample() {
	t := timer.NewTicker("tickerExample", time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C():
			fmt.Println("tick")
		case <-timer.Tick("endTicker", 5 * time.Second):
			fmt.Println("end")
			return
        }
    } 
}
```

Test ticking is as simple as timer tests:

```go
package main


import (
	"testing"

	"github.com/Ak-Army/timer"
)

func TestTimer(t *testing.T) {
	mockTimer := timer.NewMockedTickers()
	go func() {
		tickerExample()
	}()
	tickerExample := mockTimer.GetTicker("tickerExample")
	endTick := mockTimer.GetTicker("endTicker")
	tickerExample.Tick()
	tickerExample.Tick()
	endTick.TickAndWaitForStop()
}
```
