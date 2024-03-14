package main

import (
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/providers"
)

var eventsService = jobs.NewEventsService()
var queueService = jobs.NewQueueService()

func pollSqsQueue() {
	queueTicker := time.NewTicker(10 * time.Second)

	jobList := checkForNewJobs()
	sortJobs(jobList)

	go func() {
		for {
			select {
			case <-done:
				log.Println("Finished polling")
				return
			case <-queueTicker.C:
				// 1. poll to get all items in queue
				jobList = checkForNewJobs()

				// 2. if queue is empty, disable the event rule and end the function
				shutDownWhenEmpty(jobList)

				// 3. group queue jobs by provider
				sortJobs(jobList)
			}
		}
	}()

	// 4. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	go invokeWorkerTicker(providers.PolygonIo, providers.PolygonIoDelay)

	// 5. Switch off after 5min
	time.Sleep(5 * time.Minute)
}

func main() {
	lambda.Start(pollSqsQueue)
}
