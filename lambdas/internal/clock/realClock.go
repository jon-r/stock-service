package clock

import (
	"log"
	"time"
)

type Ticker struct {
	C <-chan time.Time
	// todo what else?
}

type Clock interface {
	Sleep(d time.Duration)
	Ticker(d time.Duration) *Ticker
}

type realClock struct{}

func (c *realClock) Sleep(d time.Duration) { time.Sleep(d) }
func (c *realClock) Ticker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{C: t.C}
}

func RealClock() Clock {
	return &realClock{}
}

/********** */

type timer struct {
	nextTick time.Time
	Tick     func()
}

type Mock struct {
	virtualTime time.Time

	timers []timer
}

func (m *Mock) AdvanceTime(d time.Duration) {
	t := m.virtualTime.Add(d)

	log.Println("TIMERS")

	m.timers[0].Tick()

	// todo check for tickers+sleepers
	//for i, timer := range m.timers {
	//	log.Printf("TIMER %d", i)
	//	//if timer.nextTick.After(t) {
	//	timer.Tick()
	//	//}
	//}

	m.virtualTime = t
}

func (m *Mock) Ticker(d time.Duration) *Ticker {
	//t := time.NewTicker(d)

	ticker := Ticker{
		//C
	}

	//time := timer{
	//	nextTick: time.Time{},
	//	Tick: func() {
	//
	//	},
	//}

	return &ticker
}
func (m *Mock) Sleep(d time.Duration) {
	done := make(chan struct{})
	log.Println("sleep added")

	sleepTimer := timer{
		nextTick: m.virtualTime.Add(d),
		Tick: func() {
			log.Println("sleep triggered")
			close(done)
		},
	}

	m.timers = append(m.timers, sleepTimer)

	log.Println("waiting for done")
	<-done
	log.Println("sleep done")
}

func MockClock() *Mock {
	return &Mock{
		virtualTime: time.Unix(0, 0),
	}
}

var (
	// type checking
	_ Clock = &Mock{}
)
