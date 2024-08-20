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

type queueManager struct {
	done           chan bool
	emptyResponses int
	failedAttempts int
	queues         map[provider.Name]chan job.Job
}

type handler struct {
	*handlers.LambdaHandler
	Clock clock.Clock

	queueManager queueManager
}

func newHandler(lambdaHandler *handlers.LambdaHandler, c clock.Clock) *handler {
	return &handler{
		LambdaHandler: lambdaHandler,
		Clock:         c,
		queueManager: queueManager{
			done:           make(chan bool),
			queues:         map[provider.Name]chan job.Job{provider.PolygonIo: make(chan job.Job, 20)},
			emptyResponses: 0,
			failedAttempts: 0,
		},
	}
}

var dataTickerHandler = newHandler(handlers.NewLambdaHandler(), clock.New())

func (h *handler) HandleRequest(ctx context.Context) {
	// todo look at zap docs to see if this can be done better. its not passing context to controllers
	h.Log = h.Log.LoadLambdaContext(ctx)
	defer h.Log.Sync()

	// 1. get all queued items
	go h.pollJobsQueue()

	// 2. for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
	go h.pollProviderQueue(provider.PolygonIo)

	tickerTimeout, timeErr := strconv.Atoi(os.Getenv("TICKER_TIMEOUT"))
	if timeErr != nil {
		tickerTimeout = 5
	}

	// 3. Switch off after TICKER_TIMEOUT minutes
	h.Clock.Sleep(time.Duration(tickerTimeout) * time.Minute)
	h.queueManager.done <- true

	// return nil // todo have the goroutines send the error here?
}

func main() {
	lambda.Start(dataTickerHandler.HandleRequest)
}
