package clock

import (
	"time"
)

type timer struct {
	nextTick time.Time
	Tick     func(d time.Time)
}

type Mock struct {
	virtualTime time.Time
	timers      []timer
}

func (m *Mock) AdvanceTime(d time.Duration) {
	gosched()

	t := m.virtualTime.Add(d)

	for _, timer := range m.timers {
		if t.After(timer.nextTick) {
			timer.Tick(t)
		}
	}

	m.virtualTime = t
}

func (m *Mock) NewTicker(interval time.Duration) *Ticker {
	ch := make(chan time.Time, 1)
	ticker := Ticker{c: ch, C: ch}
	ticker.Stop = func() {
		ticker.isStopped = true
	}

	prevTime := m.virtualTime

	tickTimer := timer{
		nextTick: prevTime.Add(interval),
		Tick: func(d time.Time) {
			ticks := int(d.Sub(prevTime) / interval)

			for i := 1; i <= ticks; i++ {
				if !ticker.isStopped {
					d = d.Add(interval)
					ticker.c <- d
				}
			}
		},
	}

	m.timers = append(m.timers, tickTimer)

	return &ticker
}

func (m *Mock) Sleep(d time.Duration) {
	done := make(chan bool)

	sleepTimer := timer{
		nextTick: m.virtualTime.Add(d),
		Tick: func(d time.Time) {
			done <- true
		},
	}

	m.timers = append(m.timers, sleepTimer)

	<-done
}

func MockClock() *Mock {
	return &Mock{virtualTime: time.Unix(0, 0)}
}

// Sleep momentarily so that other goroutines can process.
func gosched() { time.Sleep(1 * time.Millisecond) }

var (
	// type checking
	_ Clock = &Mock{}
)
