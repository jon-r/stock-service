package main

import (
	"log"
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/clock"
)

func DummyTimedFunction(timer clock.Clock) {
	log.Println("started")

	finished := make(chan bool)

	go func() {
		done = make(chan bool)
		counter := 0

		queueTicker := timer.NewTicker(4 * time.Second)

		for {
			select {
			case <-queueTicker.C:
				counter++
				log.Printf("ticked %v", counter)

				if counter > 4 {
					queueTicker.Stop()
					finished <- true
				}
			}
		}
	}()

	<-finished

	log.Println("finished")
}
