package main

import (
	"log"
	"testing"
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/clock"
)

func TestDummyTimedFunction(t *testing.T) {
	t.Run("testing", func(t *testing.T) {
		done := make(chan bool)

		//mocker := clock.RealClock()
		mocker := clock.MockClock()

		go func() {
			DummyTimedFunction(mocker)
			done <- true
		}()

		//time.AfterFunc(time.Second, func() {
		mocker.AdvanceTime(2 * time.Second)
		mocker.AdvanceTime(9 * time.Second)
		log.Println("waiting...")
		mocker.AdvanceTime(20 * time.Second)
		//})

		<-done

		log.Println("done")
	})
}
