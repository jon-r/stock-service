package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"go.uber.org/zap"
)

type DataWorkerHandler struct {
	queueService *jobs.QueueRepository
	dbService    *db.DatabaseRepository
	logService   *zap.SugaredLogger
	newUuid      jobs.UuidGen
}

func (handler DataWorkerHandler) handleJobAction(job jobs.JobAction) error {
	switch job.Type {
	case jobs.LoadTickerDescription:
		return handler.setTickerDescription(job.Provider, job.TickerId)
	case jobs.LoadHistoricalPrices:
		return handler.setTickerHistoricalPrices(job.Provider, job.TickerId)
	case jobs.UpdatePrices:
		return handler.updateTickerPrices(job.Provider, strings.Split(job.TickerId, ","))

	// TODO STK-86
	// jobs.LoadTickerIcon

	// TODO STK-88
	// jobs.UpdateDividends
	// jobs.LoadHistoricalDividends

	default:
		return fmt.Errorf("invalid action type = %v", job.Type)
	}
}

func (handler DataWorkerHandler) handleRequest(ctx context.Context, event jobs.JobAction) {
	// todo this might not work?
	if handler.logService == nil {
		handler.logService = logging.NewLogger(ctx)
	}
	defer handler.logService.Sync()

	var err error

	// 1. handle action
	handler.logService.Infow("Attempt to do job",
		"job", event,
	)
	err = handler.handleJobAction(event)

	if err == nil {
		handler.logService.Infoln("Job completed",
			"jobId", event.JobId,
		)
		return // job done
	}

	handler.logService.Warnw("failed to process event, re-adding it to queue",
		"jobId", event.JobId,
		"error", err,
	)

	// 2. if action failed or new queue actions after last, try again
	queueErr := handler.queueService.RetryJob(event, err.Error(), handler.newUuid)

	if queueErr != nil {
		handler.logService.Fatalw("Failed to add item to DLQ",
			"jobId", event.JobId,
			"error", queueErr,
		)
	}
}

var handler = DataWorkerHandler{
	queueService: jobs.NewQueueService(jobs.CreateSqsClient()),
	dbService:    db.NewDatabaseService(db.CreateDatabaseClient()),
}

func main() {
	lambda.Start(handler.handleRequest)
}
