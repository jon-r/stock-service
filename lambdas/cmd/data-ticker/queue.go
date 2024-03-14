package main

import (
	"log"
	"os"

	"jon-richards.com/stock-app/internal/jobs"
)

var done = make(chan bool)

func checkForNewJobs() *[]jobs.JobQueueItem {
	attempts := 0
	jobList, err := queueService.ReceiveJobs()

	if err != nil {
		attempts += 1
		log.Printf("Failed to get queue items = %v", err)
	} else {
		attempts = 0
	}

	if attempts >= 6 {
		err = eventsService.StopTickerScheduler()
		if err != nil {
			log.Printf("Failed to stop scheduler = %v", err)
		}
		log.Fatalf("Aborting after too many failed attempts")
	}

	return jobList
}

func shutDownWhenEmpty(jobList *[]jobs.JobQueueItem) {
	emptyResponses := 0
	if len(*jobList) == 0 {
		emptyResponses += 1
	} else {
		emptyResponses = 0
	}

	if emptyResponses == 6 {
		log.Println("No jobs received in 60 seconds, disabling scheduler")
		err := eventsService.StopTickerScheduler()
		if err != nil {
			log.Printf("Failed to get stop scheduler = %v", err)
		}
	}
	if emptyResponses == 12 {
		log.Printf("No jobs in 120seconds shutting down")
		os.Exit(0)
	}
}
