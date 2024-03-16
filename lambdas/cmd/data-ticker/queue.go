package main

import (
	"log"
	"os"

	"jon-richards.com/stock-app/internal/jobs"
)

var done = make(chan bool)

func checkForNewJobs(attempts int) (*[]jobs.JobQueueItem, int) {
	jobList, err := queueService.ReceiveJobs()

	if err != nil {
		attempts += 1
		log.Printf("Failed to get queue items = %v, attempts made = %v\n", err, attempts)
	} else {
		attempts = 0
	}

	if attempts >= 6 {
		err = eventsService.StopTickerScheduler()
		if err != nil {
			log.Printf("Failed to stop scheduler = %v\n", err)
		}
		log.Fatalf("Aborting after too many failed attempts")
	}

	return jobList, attempts
}

func shutDownWhenEmpty(jobList *[]jobs.JobQueueItem, emptyResponses int) int {
	if len(*jobList) == 0 {
		emptyResponses += 1
	} else {
		emptyResponses = 0
	}

	if emptyResponses == 6 {
		log.Println("No new jobs received in 60 seconds, disabling scheduler")
		err := eventsService.StopTickerScheduler()
		if err != nil {
			log.Printf("Failed to stop scheduler = %v", err)
		}
	}
	if emptyResponses == 12 {
		log.Println("No new jobs in 120seconds shutting down")
		os.Exit(0)
	}

	return emptyResponses
}
