package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

type queueList map[provider.Name]chan job.Job

type queueManager struct {
	emptyResponses int
	failedAttempts int
	queues         queueList
}

type handler struct {
	*handlers.LambdaHandler
	Clock        clock.Clock
	queueManager queueManager
}

func newHandler(lambdaHandler *handlers.LambdaHandler, c clock.Clock) *handler {
	return &handler{
		lambdaHandler,
		c,
		queueManager{queues: queueList{
			provider.PolygonIo: make(chan job.Job, 20),
		}},
	}
}

var dataTickerHandler = newHandler(handlers.NewLambdaHandler(), clock.New())

func (h *handler) HandleRequest(ctx context.Context) {
	tickerTimeout, timeErr := strconv.Atoi(os.Getenv("TICKER_TIMEOUT"))
	if timeErr != nil {
		tickerTimeout = 5
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	h.Log.LoadContext(ctx)
	defer h.Log.Sync()

	// 1. get all queued items
	go h.pollJobsQueue(ctx, cancel)

	// 2. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	go h.pollProviderQueue(ctx, provider.PolygonIo)

	// 3. Switch off after TICKER_TIMEOUT minutes
	h.Clock.Sleep(time.Duration(tickerTimeout) * time.Minute)
	cancel()

	// return nil // todo have the goroutines send the error here?
}

func main() {
	lambda.Start(dataTickerHandler.HandleRequest)
}
