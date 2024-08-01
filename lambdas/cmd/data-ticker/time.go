package main

import (
	"time"
)

type TimeHandler interface {
	Sleep(d time.Duration)
	NewTicker(d time.Duration) *time.Ticker
}
type t struct{}

func (t) Sleep(d time.Duration)                  { time.Sleep(d) }
func (t) NewTicker(d time.Duration) *time.Ticker { return time.NewTicker(d) }
