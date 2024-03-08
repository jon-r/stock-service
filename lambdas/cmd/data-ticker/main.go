package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/jobs"
)

var eventsService = jobs.NewEventsService()

func handleRequest() {
	var err error

	log.Println("Hello world")

	// 1. get all items in queue

	// 2. if queue is empty, disable the event rule and end the function
	err = eventsService.StopTickerScheduler()

	if err != nil {
		log.Fatalf("Failed to stop event ticker, %v", err)
	} else {
		log.Println("Stopped event ticker")
	}

	// 3. group queue jobs by provider

	// 4 for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
}

func main() {
	lambda.Start(handleRequest)
}
