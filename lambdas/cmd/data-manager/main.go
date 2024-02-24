package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"log"

	"jon-richards.com/stock-app/db"
	"jon-richards.com/stock-app/queue"
	s "jon-richards.com/stock-app/session"
)

var awsSession = s.NewAwsSession()
var dbService = db.NewDatabaseService(*awsSession)
var queueService = queue.NewQueueService(*awsSession)

var fakeJobs = []db.JobInput{
	{Name: "Phoebe", Group: "long", Speed: 5},
	{Name: "Harley", Group: "long", Speed: 5},
	{Name: "Bandit", Group: "long", Speed: 5},
	{Name: "Delilah", Group: "long", Speed: 5},
	{Name: "Tiger", Group: "long", Speed: 5},
	{Name: "Panda", Group: "long", Speed: 5},

	{Name: "Whiskey", Group: "short", Speed: 3},
	{Name: "Jasper", Group: "short", Speed: 3},
	{Name: "Belle", Group: "short", Speed: 3},
	{Name: "Shelby", Group: "short", Speed: 3},
	{Name: "Zara", Group: "short", Speed: 3},
	{Name: "Bruno", Group: "short", Speed: 3},
}

// todo have the delay and everything as some sort of config file?
var fakeQueueEvents = []queue.QueueMessage{
	{Group: "long", Delay: 5},
	{Group: "short", Delay: 3},
}

func loadData() {
	var err error

	err = dbService.InsertJobs(fakeJobs)

	if err != nil {
		log.Fatalf("Error adding data to DB: %s", err)
	} else {
		log.Println("Successfully added items to DB")
	}

	err = queueService.InsertDelayedEvents(fakeQueueEvents)

	if err != nil {
		log.Fatalf("Error adding items to Queue: %s", err)
	} else {
		log.Println("Successfully added items to Queue")
	}
}

func main() {
	lambda.Start(loadData)
}
