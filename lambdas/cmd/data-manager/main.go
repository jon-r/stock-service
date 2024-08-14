package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/config"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
	"go.uber.org/zap/zapcore"
)

type dataManagerHandler interface {
	HandleRequest(ctx context.Context) error
}

type handler struct {
	tickers tickers.Controller
	jobs    jobs.Controller
	log     logger.Logger
}

func newHandler() dataManagerHandler {
	cfg := config.GetAwsConfig()
	log := logger.NewLogger(zapcore.InfoLevel)
	idGen := uuid.NewString

	// todo once tests split up, some of this can be moved to the controller
	queueBroker := queue.NewBroker(cfg, idGen)
	eventsScheduler := events.NewScheduler(cfg)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, eventsScheduler, idGen, log)
	tickersCtrl := tickers.NewController(dbRepository, nil, log)

	return &handler{tickersCtrl, jobsCtrl, log}
}

func (h *handler) HandleRequest(ctx context.Context) error {
	// todo look at zap docs to see if this can be done better
	h.log = h.log.LoadLambdaContext(ctx)

	// 1. get all tickers
	tickerList, err := h.tickers.GetAll()

	if err != nil {
		return err
	}

	// 2. convert the jobs into update actions
	return h.jobs.LaunchDailyTickerJobs(tickerList)
}

var serviceHandler = newHandler()

func main() {
	lambda.Start(serviceHandler.HandleRequest)
}
