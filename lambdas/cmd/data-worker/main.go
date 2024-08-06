package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"github.com/jon-r/stock-service/lambdas/internal/providers"
	"github.com/jon-r/stock-service/lambdas/internal/types"
)

type DataWorkerHandler struct {
	types.ServiceHandler
	ProviderService providers.ProviderService
}

func (handler DataWorkerHandler) doJob(job jobs.JobAction) error {
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

func (handler DataWorkerHandler) handleJobEvent(ctx context.Context, event jobs.JobAction) error {
	// todo this might not work?
	if handler.LogService == nil {
		handler.LogService = logging.NewLogger(ctx)
	}
	defer handler.LogService.Sync()

	var err error

	// 1. handle action
	handler.LogService.Infow("Attempt to do job",
		"job", event,
	)
	err = handler.doJob(event)

	if err == nil {
		handler.LogService.Infoln("Job completed",
			"jobId", event.JobId,
		)
		return nil // job done
	}

	handler.LogService.Warnw("failed to process event, re-adding it to queue",
		"jobId", event.JobId,
		"error", err,
	)

	// 2. if action failed or new queue actions after last, try again
	queueErr := handler.QueueService.RetryJob(event, err.Error(), handler.NewUuid)

	if queueErr != nil {
		handler.LogService.Fatalw("Failed to add item to DLQ",
			"jobId", event.JobId,
			"error", queueErr,
		)
		return queueErr
	}

	return err
}

var serviceHandler = types.ServiceHandler{
	QueueService: jobs.NewQueueService(jobs.CreateSqsClient()),
	DbService:    db.NewDatabaseService(db.CreateDatabaseClient()),
}

func main() {
	handler := DataWorkerHandler{ServiceHandler: serviceHandler, ProviderService: providers.NewProviderService()}
	lambda.Start(handler.handleJobEvent)
}
