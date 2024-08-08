package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/jobs_old"
	"github.com/jon-r/stock-service/lambdas/internal/logging_old"
	"github.com/jon-r/stock-service/lambdas/internal/providers_old"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler_old"
	"github.com/jon-r/stock-service/lambdas/internal/types_old"
)

type DataTickerHandler struct {
	types_old.ServiceHandler
	Clock clock.Clock
	done  chan bool
}

func (handler DataTickerHandler) handleQueuedJobs(ctx context.Context) {
	// todo this might not work?
	if handler.LogService == nil {
		handler.LogService = logging_old.NewLogger(ctx)
	}
	defer handler.LogService.Sync()

	// 1. get all queued items
	go handler.checkForJobs()

	// 2. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	go handler.invokeWorkerTicker(providers_old.PolygonIo, providers_old.PolygonIoDelay)

	// 3. Switch off after TICKER_TIMEOUT min
	tickerTimeout, timeErr := strconv.Atoi(os.Getenv("TICKER_TIMEOUT"))
	if timeErr != nil {
		tickerTimeout = 5
	}

	handler.Clock.Sleep(time.Duration(tickerTimeout) * time.Minute)
	handler.done <- true

	//return nil // todo have the goroutines send the error here?
}

var serviceHandler = types_old.ServiceHandler{
	QueueService:  jobs_old.NewQueueService(jobs_old.CreateSqsClient()),
	EventsService: scheduler_old.NewEventsService(scheduler_old.CreateEventClients()),
	NewUuid:       uuid.NewString,
}

func main() {
	handler := DataTickerHandler{
		ServiceHandler: serviceHandler,
		Clock:          clock.New(),
		done:           make(chan bool),
	}
	lambda.Start(handler.handleQueuedJobs)
}
