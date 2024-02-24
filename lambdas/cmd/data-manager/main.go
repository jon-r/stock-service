package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"jon-richards.com/stock-app/db"
	"jon-richards.com/stock-app/queue"
	s "jon-richards.com/stock-app/session"
)

var awsSession = s.NewAwsSession()
var dbService = db.NewDatabaseService(*awsSession)
var queueService = queue.NewQueueService(*awsSession)

var fakeJobs = []db.JobInput{
	{Name: "Phoebe", Group: "long"},
	{Name: "Harley", Group: "long"},
	{Name: "Bandit", Group: "long"},
	{Name: "Delilah", Group: "long"},
	{Name: "Tiger", Group: "long"},
	{Name: "Panda", Group: "long"},

	{Name: "Whiskey", Group: "short"},
	{Name: "Jasper", Group: "short"},
	{Name: "Belle", Group: "short"},
	{Name: "Shelby", Group: "short"},
	{Name: "Zara", Group: "short"},
	{Name: "Bruno", Group: "short"},
}

// todo have the delay and everything as some sort of config file?
var fakeQueueEvents = []queue.QueueMessage{
	{Group: "long", Delay: 5},
	{Group: "short", Delay: 3},
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
