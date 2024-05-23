package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"github.com/jon-r/stock-service/lambdas/internal/providers"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"go.uber.org/zap"
)

type DataTickerHandler struct {
	queueService  *jobs.QueueRepository
	eventsService *scheduler.EventsRepository
	log           *zap.SugaredLogger
}

func (handler DataTickerHandler) pollSqsQueue(ctx context.Context) {
	// todo this might not work?
	handler.log = logging.NewLogger(ctx)
	defer handler.log.Sync()

	queueTicker := time.NewTicker(10 * time.Second)
	tickerTimeout, err := strconv.Atoi(os.Getenv("TICKER_TIMEOUT"))

	if err != nil {
		tickerTimeout = 5
	}

	jobList, attempts := handler.checkForNewJobs(0)
	sortJobs(jobList)

	go func() {
		handler.log.Infoln("Started polling...")
		emptyResponses := 0

		for {
			select {
			case <-done:
				handler.log.Infoln("Finished polling")
				return
			case <-queueTicker.C:
				// 1. poll to get all items in queue
				jobList, attempts = handler.checkForNewJobs(attempts)

				// 2. if queue is empty, disable the event rule and end the function
				emptyResponses = handler.shutDownWhenEmpty(jobList, emptyResponses)

				// 3. group queue jobs by provider
				sortJobs(jobList)
			}
		}
	}()

	// 4. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	go handler.invokeWorkerTicker(providers.PolygonIo, providers.PolygonIoDelay)

	// 5. Switch off after 5min
	time.Sleep(time.Duration(tickerTimeout) * time.Minute)
	done <- true
}

var handler = DataTickerHandler{
	queueService:  jobs.NewQueueService(jobs.CreateSqsClient()),
	eventsService: scheduler.NewEventsService(scheduler.CreateEventClients()),
}

func main() {
	lambda.Start(handler.pollSqsQueue)
}
