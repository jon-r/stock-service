package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/db"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
)

var dbService = db.NewDatabaseService()

var queueService = jobs.NewQueueService()

func handleRequest(ctx context.Context, event jobs.JobAction) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	var err error

	// 1. handle action
	err = handleJobAction(event)

	if err == nil {
		log.Infoln("Job completed",
			"jobId", event.JobId,
		)
		return // job done
	}

	var queueErr error
	// 2. if action failed <= 3 times, or new queue actions after last, add to the queue
	if event.Attempts <= 3 {
		log.Warnw("failed to process event, re-adding it to queue",
			"jobId", event.JobId,
			"error", err,
		)
		queueErr = retryFailedJob(event)

		if queueErr == nil {
			return
		}
	}

	var failReason error
	if queueErr != nil {
		failReason = queueErr
	} else {
		failReason = err
	}

	// 3. if action failed 3 times, or was not able to relist it, put in DLQ
	log.Errorw("failed to process event, adding it to dlq",
		"jobId", event.JobId,
		"error", failReason,
	)
	queueErr = queueService.AddJobToDLQ(event, failReason)

	if queueErr != nil {
		log.Fatalw("Failed to add item to DL",
			"error", err,
		)
	}
}

func main() {
	lambda.Start(handleRequest)
}
