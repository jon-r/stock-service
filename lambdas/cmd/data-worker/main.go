package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/db"
	"jon-richards.com/stock-app/internal/jobs"
)

var dbService = db.NewDatabaseService()

var queueService = jobs.NewQueueService()

// todo custom event from data ticker (prob be similar to job table?)
func handleRequest(ctx context.Context, event jobs.JobAction) {
	var err error

	// 1. handle action
	err = handleJobAction(event)

	if err == nil {
		log.Printf("Job %v completed", event.JobId)
		return // job done
	}

	var queueErr error
	// 2. if action failed <= 3 times, or new queue actions after last, add to the queue
	if event.Attempts <= 3 {
		log.Printf("failed to process event %v, readding it to queue: %v\n", event.JobId, err)
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
	log.Printf("Job %v failed %d times, adding to DQL", event.JobId, event.Attempts)
	queueErr = queueService.AddJobToDLQ(event, failReason)

	if queueErr != nil {
		log.Fatalf("Failed to add item to DLQ: %v", err)
	}
}

func main() {
	lambda.Start(handleRequest)
}
