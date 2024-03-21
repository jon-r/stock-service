package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/logging"
)

//var dbService = db.NewDatabaseService()
//var queueService = queue.NewQueueService()

//var fakeJobs = []db.JobInput{
//	{Provider: providers.Fast, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Phoebe"}},
//	{Provider: providers.Fast, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Harley"}},
//	{Provider: providers.Fast, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Bandit"}},
//	{Provider: providers.Fast, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Delilah"}},
//	{Provider: providers.Fast, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Tiger"}},
//	{Provider: providers.Fast, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Panda"}},
//
//	{Provider: providers.Slow, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Whiskey"}},
//	{Provider: providers.Slow, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Jasper"}},
//	{Provider: providers.Slow, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Belle"}},
//	{Provider: providers.Slow, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Shelby"}},
//	{Provider: providers.Slow, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Zara"}},
//	{Provider: providers.Slow, Action: "INSERT_DATA_TEST", Payload: db.JobInputPayload{"Name": "Bruno"}},
//}

//var fakeQueueEvents = []queue.Message{
//	{Provider: providers.Fast},
//	{Provider: providers.Slow},
//}

func handleRequest(ctx context.Context) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	log.Infoln("Hello world")

	//var err error
	//
	//err = dbService.InsertJobs(fakeJobs)
	//
	//if err != nil {
	//	logging.Fatalf("Error adding data to DB: %s", err)
	//} else {
	//	logging.Println("Successfully added items to DB")
	//}
	//
	//err = queueService.SendDelayedEvents(fakeQueueEvents)
	//
	//if err != nil {
	//	logging.Fatalf("Error adding items to Queue: %s", err)
	//} else {
	//	logging.Println("Successfully added items to Queue")
	//}
}

func main() {
	lambda.Start(handleRequest)
}
