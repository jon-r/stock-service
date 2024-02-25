package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"jon-richards.com/stock-app/db"
	"jon-richards.com/stock-app/queue"
	"jon-richards.com/stock-app/remote"
)

var dbService = db.NewDatabaseService()
var queueService = queue.NewQueueService()

func handleRequest(ctx context.Context, event events.SQSEvent) {
	var err error

	message, err := queue.ParseQueueEvent(event)

	if err != nil {
		log.Fatalf("Error parsing queue event: %s", err)
	} else {
		log.Printf("Handling event: %s", message.Group)
	}

	job, err := dbService.FindJobByGroup(message.Group)

	if err != nil {
		log.Fatalf("Error getting item: %s", err)
	} else if job == nil {
		// no more items to fetch
		return
	} else {
		log.Printf("Job: %s", job.Name)
	}

	res, err := remote.FetchDogItem()

	if err != nil {
		log.Fatalf("Error calling http.get: %s", err)
	}

	err = dbService.UpsertStockItem(res, job)

	if err != nil {
		log.Fatalf("Error calling dynamodb.WriteItem: %s", err)
	} else {
		log.Println("Successfully added items to tables")
	}

	err = dbService.DeleteJob(job)

	if err != nil {
		log.Fatalf("Error calling dynamodb.DeleteItem: %s", err)
	}

	err = queueService.SendDelayedEvents([]queue.QueueMessage{*message})

	if err != nil {
		log.Fatalf("Error adding item to Queue: %s", err)
	} else {
		log.Println("Successfully added item to Queue")
	}
}

func main() {
	lambda.Start(handleRequest)
}
