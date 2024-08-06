package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"github.com/jon-r/stock-service/lambdas/internal/types"
)

type DataTickerHandler struct {
	types.ServiceHandler
	Clock clock.Clock
	done  chan bool
}

func (handler DataTickerHandler) handleQueuedJobs(ctx context.Context) {
	// todo this might not work?
	if handler.LogService == nil {
		handler.LogService = logging.NewLogger(ctx)
	}
	defer handler.LogService.Sync()

	// 1. get all queued items
	go handler.checkForJobs()

	// 2. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	// fixme reenable this
	//go handler.invokeWorkerTicker(providers.PolygonIo, providers.PolygonIoDelay)

	// 3. Switch off after TICKER_TIMEOUT min
	tickerTimeout, timeErr := strconv.Atoi(os.Getenv("TICKER_TIMEOUT"))
	if timeErr != nil {
		tickerTimeout = 5
	}

	handler.LogService.Debugln("waiting %v mins", tickerTimeout)
	handler.Clock.Sleep(time.Duration(tickerTimeout) * time.Minute)
	handler.done <- true
	handler.LogService.Debugln("DONE?")

	//return nil // todo have the goroutines send the error here?
}

var serviceHandler = types.ServiceHandler{
	QueueService:  jobs.NewQueueService(jobs.CreateSqsClient()),
	EventsService: scheduler.NewEventsService(scheduler.CreateEventClients()),
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
