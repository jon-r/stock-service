package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
	"jon-richards.com/stock-app/internal/providers"
	"jon-richards.com/stock-app/internal/scheduler"
)

var eventsService = scheduler.NewEventsService()
var queueService = jobs.NewQueueService()

func pollSqsQueue(ctx context.Context) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	queueTicker := time.NewTicker(10 * time.Second)
	tickerTimeout, err := strconv.Atoi(os.Getenv("TICKER_TIMEOUT"))

	if err != nil {
		tickerTimeout = 5
	}

	jobList, attempts := checkForNewJobs(ctx, 0)
	sortJobs(jobList)

	go func() {
		log.Infoln("Started polling...")
		emptyResponses := 0

		for {
			select {
			case <-done:
				log.Infoln("Finished polling")
				return
			case <-queueTicker.C:
				// 1. poll to get all items in queue
				jobList, attempts = checkForNewJobs(ctx, attempts)

				// 2. if queue is empty, disable the event rule and end the function
				emptyResponses = shutDownWhenEmpty(ctx, jobList, emptyResponses)

				// 3. group queue jobs by provider
				sortJobs(jobList)
			}
		}
	}()

	// 4. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	go invokeWorkerTicker(ctx, providers.PolygonIo, providers.PolygonIoDelay)

	// 5. Switch off after 5min
	time.Sleep(time.Duration(tickerTimeout) * time.Minute)
	done <- true
}

func main() {
	lambda.Start(pollSqsQueue)
}
