package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"jon-richards.com/stock-app/db"
	"jon-richards.com/stock-app/queue"
)

var dbService = db.NewDatabaseService()
var queueService = queue.NewQueueService()

var fakeJobs = []db.JobInput{
	{QueueGroup: "long", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Phoebe"}},
	{QueueGroup: "long", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Harley"}},
	{QueueGroup: "long", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Bandit"}},
	{QueueGroup: "long", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Delilah"}},
	{QueueGroup: "long", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Tiger"}},
	{QueueGroup: "long", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Panda"}},

	{QueueGroup: "short", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Whiskey"}},
	{QueueGroup: "short", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Jasper"}},
	{QueueGroup: "short", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Belle"}},
	{QueueGroup: "short", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Shelby"}},
	{QueueGroup: "short", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Zara"}},
	{QueueGroup: "short", Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Bruno"}},
}

var fakeQueueEvents = []queue.Message{
	{QueueGroup: "long"},
	{QueueGroup: "short"},
}

func handleRequest() {
	var err error

	err = dbService.InsertJobs(fakeJobs)

	if err != nil {
		log.Fatalf("Error adding data to DB: %s", err)
	} else {
		log.Println("Successfully added items to DB")
	}

	err = queueService.SendDelayedEvents(fakeQueueEvents)

	if err != nil {
		log.Fatalf("Error adding items to Queue: %s", err)
	} else {
		log.Println("Successfully added items to Queue")
	}
}

func main() {
	lambda.Start(handleRequest)
}
