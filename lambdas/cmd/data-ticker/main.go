package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/clock"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"github.com/jon-r/stock-service/lambdas/internal/types"
)

type DataTickerHandler struct {
	types.ServiceHandler
	Clock clock.Clock
}

func (handler DataTickerHandler) pollSqsQueue(ctx context.Context) {
	// todo this might not work?
	if handler.LogService == nil {
		handler.LogService = logging.NewLogger(ctx)
	}
	defer handler.LogService.Sync()

	//queueTicker := handler.Timer.NewTicker(10 * time.Second)
	tickerTimeout, timeErr := strconv.Atoi(os.Getenv("TICKER_TIMEOUT"))

	if timeErr != nil {
		tickerTimeout = 5
	}

	//go handler.checkForJobs()

	// todo keep this for now in case it stops working
	//jobList, attempts := handler.receiveNewJobs(0)
	//sortJobs(jobList)
	//
	//go func() {
	//	handler.LogService.Infoln("Started polling...")
	//	emptyResponses := 0
	//
	//	for {
	//		select {
	//		case <-done:
	//			handler.LogService.Infoln("Finished polling")
	//			return
	//		case <-queueTicker.C:
	//			// 1. poll to get all items in queue
	//			jobList, attempts = handler.receiveNewJobs(attempts)
	//
	//			// 2. if queue is empty, disable the event rule and end the function
	//			emptyResponses = handler.shutDownWhenEmpty(jobList, emptyResponses)
	//
	//			// 3. group queue jobs by provider
	//			sortJobs(jobList)
	//		}
	//	}
	//}()

	// 4. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	//go handler.invokeWorkerTicker(providers.PolygonIo, providers.PolygonIoDelay)

	// 5. Switch off after 5min
	handler.LogService.Infoln("started")
	handler.Clock.Sleep(time.Duration(tickerTimeout) * time.Minute)
	handler.LogService.Infoln("DONE?")
	done <- true

	//return nil // todo have the goroutines send the error here?
}

var serviceHandler = types.ServiceHandler{
	QueueService:  jobs.NewQueueService(jobs.CreateSqsClient()),
	EventsService: scheduler.NewEventsService(scheduler.CreateEventClients()),
	NewUuid:       uuid.NewString,
}

func main() {
	handler := DataTickerHandler{ServiceHandler: serviceHandler, Clock: clock.RealClock()}
	lambda.Start(handler.pollSqsQueue)
}
