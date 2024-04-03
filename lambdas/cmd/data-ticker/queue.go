package main

import (
	"context"

	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
)

var done = make(chan bool)

func checkForNewJobs(ctx context.Context, attempts int) (*[]jobs.JobQueueItem, int) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	log.Infoln("attempt to receive jobs...")
	jobList, err := queueService.ReceiveJobs()

	if err != nil {
		attempts += 1
		log.Warnw("Failed to get queue items",
			"attempts", attempts,
		)
	} else {
		attempts = 0
	}

	if attempts >= 6 {
		err = eventsService.StopTickerScheduler()
		if err != nil {
			log.Errorw("Failed to stop scheduler",
				"error", err,
			)
		}
		log.Fatalln("Aborting after too many failed attempts")
	}

	return jobList, attempts
}

func shutDownWhenEmpty(ctx context.Context, jobList *[]jobs.JobQueueItem, emptyResponses int) int {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	if len(*jobList) == 0 {
		emptyResponses += 1
	} else {
		// todo remove this
		log.Infof("Found %v jobs", len(*jobList))
		emptyResponses = 0
	}

	if emptyResponses == 6 {
		log.Infoln("No new jobs received in 60 seconds, disabling scheduler")
		err := eventsService.StopTickerScheduler()
		if err != nil {
			log.Errorw("Failed to stop scheduler",
				"error", err,
			)
		}
	}

	return emptyResponses
}
