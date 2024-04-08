package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"jon-richards.com/stock-app/internal/db"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
)

var dbService = db.NewDatabaseService()

var queueService = jobs.NewQueueService()

func handleJobAction(log *zap.SugaredLogger, job jobs.JobAction) error {
	switch job.Type {
	case jobs.LoadTickerDescription:
		return setTickerDescription(log, job.Provider, job.TickerId)
	case jobs.LoadHistoricalPrices:
		return setTickerHistoricalPrices(log, job.Provider, job.TickerId)

	// TODO STK-112
	//case jobs.UpdatePrices:
	//	return updateTickerPrices(ctx, job.Provider, strings.Split(job.TickerId, ","))

	// TODO STK-86
	// jobs.LoadTickerIcon

	// TODO STK-88
	// jobs.UpdateDividends
	// jobs.LoadHistoricalDividends

	default:
		return fmt.Errorf("invalid action type = %v", job.Type)
	}
}

func handleRequest(ctx context.Context, event jobs.JobAction) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	var err error

	// 1. handle action
	log.Infow("Attempt to do job",
		"job", event,
	)
	err = handleJobAction(log, event)

	if err == nil {
		log.Infoln("Job completed",
			"jobId", event.JobId,
		)
		return // job done
	}

	log.Warnw("failed to process event, re-adding it to queue",
		"jobId", event.JobId,
		"error", err,
	)

	// 2. if action failed or new queue actions after last, try again
	queueErr := queueService.RetryJob(event, err.Error())

	if queueErr != nil {
		log.Fatalw("Failed to add item to DLQ",
			"jobId", event.JobId,
			"error", queueErr,
		)
	}
}

func main() {
	lambda.Start(handleRequest)
}
