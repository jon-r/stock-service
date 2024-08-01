package clock

import (
	"time"
)

type Ticker struct {
	C         <-chan time.Time
	c         chan time.Time
	Stop      func()
	isStopped bool
	// todo what else?
}

type Clock interface {
	Sleep(d time.Duration)
	NewTicker(d time.Duration) *Ticker
}

type realClock struct{}

func (c *realClock) Sleep(d time.Duration) { time.Sleep(d) }
func (c *realClock) NewTicker(d time.Duration) *Ticker {
	t := time.NewTicker(d)
	return &Ticker{C: t.C}
}
func RealClock() Clock { return &realClock{} }
