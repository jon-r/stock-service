package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"jon-richards.com/stock-app/internal/db"
	"jon-richards.com/stock-app/internal/queue"
)

var dbService = db.NewDatabaseService()
var queueService = queue.NewQueueService()

// todo custom event from data ticker (prob be similar to job table?)
func handleRequest(ctx context.Context, event events.SQSEvent) {
	var err error

	message, err := queue.ParseQueueEvent(event)

	if err != nil {
		log.Fatalf("Error parsing queue event: %s", err)
	} else {
		log.Printf("Handling event: %s", message.Provider)
	}

	job, err := dbService.FindJobByProvider(message.Provider)

	if err != nil {
		log.Fatalf("Error getting item: %s", err)
	} else if job == nil {
		// no more items to fetch
		return
	} else {
		log.Printf("Job: %s", job.JobId)
	}

	err = handleJobAction(job.JobInput)

	if err != nil {
		log.Fatalf("Error with action '%s': %s", job.JobInput.Type, err)
	} else {
		log.Printf("Successfully ran action %v", job.JobInput)
	}

	err = dbService.DeleteJob(job)

	if err != nil {
		log.Fatalf("Error calling dynamodb.DeleteItem: %s", err)
	}

	err = queueService.SendDelayedEvents([]queue.Message{message.Message})

	if err != nil {
		log.Fatalf("Error adding item to Queue: %s", err)
	} else {
		log.Println("Successfully added item to Queue")
	}
}

func main() {
	lambda.Start(handleRequest)
}
