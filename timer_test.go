package timer

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestAfterFunc(t *testing.T) {
	i := 10
	c := make(chan bool)
	var f func()
	f = func() {
		i--
		if i >= 0 {
			time.AfterFunc(0, f)
			time.Sleep(1 * time.Second)
		} else {
			c <- true
		}
	}

	AfterFunc("testAfterFunc", 0, f)
	<-c
}

func TestTimerStopStress(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			timer := AfterFunc("TestTimerStopStress", 2*time.Second, func() {
				t.Errorf("timer %d was not stopped", i)
			})
			time.Sleep(1 * time.Second)
			timer.SafeStop()
		}(i)
	}
	wg.Wait()
}

func TestTimerStopDoneFunc(t *testing.T) {
	done := make(chan struct{})
	timer := AfterFunc("TestTimerStopStress", 0, func() {
		close(done)
	})
	<-done
	timer.SafeStop()
	timer.SafeStop()
}

func TestTimerStopLongRunning(t *testing.T) {
	done := make(chan struct{})
	timer := AfterFunc("TestTimerStopStress", 0, func() {
		done <- struct{}{}
		time.Sleep(5 * time.Millisecond)
		close(done)
	})
	<-done
	time.Sleep(time.Millisecond)
	timer.SafeStop()
	timer.SafeStop()
	<-done
}

func benchmark(b *testing.B, bench func(n int)) {

	// Create equal number of garbage timers on each P before starting
	// the benchmark.
	var wg sync.WaitGroup
	garbageAll := make([][]*time.Timer, runtime.GOMAXPROCS(0))
	for i := range garbageAll {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			garbage := make([]*time.Timer, 1<<15)
			for j := range garbage {
				garbage[j] = time.AfterFunc(time.Hour, nil)
			}
			garbageAll[i] = garbage
		}(i)
	}
	wg.Wait()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bench(1000)
		}
	})
	b.StopTimer()

	for _, garbage := range garbageAll {
		for _, t := range garbage {
			t.Stop()
		}
	}
}

func BenchmarkAfterFunc(b *testing.B) {
	benchmark(b, func(n int) {
		c := make(chan bool)
		var f func()
		f = func() {
			n--
			if n >= 0 {
				time.AfterFunc(0, f)
			} else {
				c <- true
			}
		}

		AfterFunc("BenchmarkAfterFunc", 0, f)
		<-c
	})
}

func BenchmarkStartStop(b *testing.B) {
	benchmark(b, func(n int) {
		timers := make([]Timer, n)
		for i := 0; i < n; i++ {
			timers[i] = AfterFunc("BenchmarkStartStop", time.Hour, nil)
		}

		for i := 0; i < n; i++ {
			timers[i].SafeStop()
		}
	})
}
