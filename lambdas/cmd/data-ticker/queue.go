package main

import "github.com/jon-r/stock-service/lambdas/internal/jobs"

var done = make(chan bool)

func (handler DataTickerHandler) checkForNewJobs(attempts int) (*[]jobs.JobQueueItem, int) {
	handler.LogService.Infoln("attempt to receive jobs...")
	jobList, err := handler.QueueService.ReceiveJobs()

	if err != nil {
		attempts += 1
		handler.LogService.Warnw("Failed to get queue items",
			"attempts", attempts,
		)
	} else {
		attempts = 0
	}

	if attempts >= 6 {
		err = handler.EventsService.StopTickerScheduler()
		if err != nil {
			handler.LogService.Errorw("Failed to stop scheduler",
				"error", err,
			)
		}
		handler.LogService.Fatalln("Aborting after too many failed attempts")
	}

	return jobList, attempts
}

func (handler DataTickerHandler) shutDownWhenEmpty(jobList *[]jobs.JobQueueItem, emptyResponses int) int {
	if len(*jobList) == 0 {
		emptyResponses += 1
	} else {
		emptyResponses = 0
	}

	if emptyResponses == 6 {
		handler.LogService.Infoln("No new jobs received in 60 seconds, disabling scheduler")
		err := handler.EventsService.StopTickerScheduler()
		if err != nil {
			handler.LogService.Errorw("Failed to stop scheduler",
				"error", err,
			)
		}
	}

	return emptyResponses
}
