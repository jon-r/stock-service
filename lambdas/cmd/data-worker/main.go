package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/config"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/providers"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/prices"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
	"go.uber.org/zap/zapcore"
)

type dataWorkerHandler interface {
	HandleRequest(ctx context.Context, job job.Job) error
}

type handler struct {
	tickers tickers.Controller
	jobs    jobs.Controller
	prices  prices.Controller
	log     logger.Logger
}

func newHandler() dataWorkerHandler {
	cfg := config.GetAwsConfig()
	log := logger.NewLogger(zapcore.InfoLevel)

	// todo once tests split up, some of this can be moved to the controller
	providersService := providers.NewService(nil)
	queueBroker := queue.NewBroker(cfg, nil)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, nil, nil, log)
	tickersCtrl := tickers.NewController(dbRepository, providersService, log)
	pricesCtrl := prices.NewController(dbRepository, providersService, log)

	return &handler{tickersCtrl, jobsCtrl, pricesCtrl, log}
}

func (h *handler) HandleRequest(ctx context.Context, j job.Job) error {
	// todo look at zap docs to see if this can be done better
	h.log = h.log.LoadLambdaContext(ctx)

	// 1. handle action
	err := h.doJob(j)

	if err == nil {
		h.log.Infoln("job completed", "jobId", j.JobId)
		return nil // job done
	}

	// 2. if action failed or new queue actions after last, try again
	h.log.Warnw("failed to process event, re-adding it to queue",
		"jobId", j.JobId,
		"error", err,
	)

	queueErr := h.jobs.RetryJob(j, err.Error())

	if queueErr != nil {
		h.log.Errorw("Failed to add item to DLQ",
			"jobId", j.JobId,
			"error", queueErr,
		)
		return queueErr
	}

	return err
}

var serviceHandler = newHandler()

func main() {
	lambda.Start(serviceHandler.HandleRequest)
}
