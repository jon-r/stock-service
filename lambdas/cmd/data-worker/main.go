package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/db_old"
	"github.com/jon-r/stock-service/lambdas/internal/jobs_old"
	"github.com/jon-r/stock-service/lambdas/internal/logging_old"
	"github.com/jon-r/stock-service/lambdas/internal/providers_old"
	"github.com/jon-r/stock-service/lambdas/internal/types_old"
)

type DataWorkerHandler struct {
	types_old.ServiceHandler
	ProviderService providers_old.ProviderService
}

func (handler DataWorkerHandler) doJob(job jobs_old.JobAction) error {
	switch job.Type {
	case jobs_old.LoadTickerDescription:
		return handler.setTickerDescription(job.Provider, job.TickerId)
	case jobs_old.LoadHistoricalPrices:
		return handler.setTickerHistoricalPrices(job.Provider, job.TickerId)
	case jobs_old.UpdatePrices:
		return handler.updateTickerPrices(job.Provider, strings.Split(job.TickerId, ","))

	// TODO STK-86
	// jobs_old.LoadTickerIcon

	// TODO STK-88
	// jobs_old.UpdateDividends
	// jobs_old.LoadHistoricalDividends

	default:
		return fmt.Errorf("invalid action type = %v", job.Type)
	}
}

func (handler DataWorkerHandler) handleJobEvent(ctx context.Context, event jobs_old.JobAction) error {
	// todo this might not work?
	if handler.LogService == nil {
		handler.LogService = logging_old.NewLogger(ctx)
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

var serviceHandler = types_old.ServiceHandler{
	QueueService: jobs_old.NewQueueService(jobs_old.CreateSqsClient()),
	DbService:    db_old.NewDatabaseService(db_old.CreateDatabaseClient()),
}

func main() {
	handler := DataWorkerHandler{ServiceHandler: serviceHandler, ProviderService: providers_old.NewProviderService()}
	lambda.Start(handler.handleJobEvent)
}
