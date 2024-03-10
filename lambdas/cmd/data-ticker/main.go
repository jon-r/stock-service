package main

import (
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/providers"
)

var eventsService = jobs.NewEventsService()
var queueService = jobs.NewQueueService()

func pollProvider(provider providers.ProviderName, interval time.Duration) {
	tk := time.NewTicker(interval)

	for range tk.C {

	}
}

func pollSqsQueue() {
	var err error

	//for provider, timer := range providers.SettingsTimers {
	//	go pollProvider(provider, time.Duration(timer) * time.Second)
	//}

	log.Println("Hello world")

	// 1. poll to get all items in queue
	jobs, err := queueService.GetJobsFromQueue()

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
	lambda.Start(pollSqsQueue)
}
