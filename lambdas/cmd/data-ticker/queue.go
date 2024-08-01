package main

import (
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/jobs"
)

var done = make(chan bool)

func (handler DataTickerHandler) checkForJobs() {
	queueTicker := handler.Clock.Ticker(10 * time.Second)

	var jobList *[]jobs.JobQueueItem

	emptyResponses := 0

	jobList, attempts := handler.receiveNewJobs(0)
	sortJobs(jobList)

	handler.LogService.Infoln("Started polling...")

	for {
		select {
		case <-done:
			handler.LogService.Infoln("Finished polling")
			return
		case <-queueTicker.C:
			handler.LogService.Infoln("TICK?")
			// 1. poll to get all items in queue
			jobList, attempts = handler.receiveNewJobs(attempts)

			// 2. if queue is empty, disable the event rule and end the function
			emptyResponses = handler.shutDownWhenEmpty(jobList, emptyResponses)

			// 3. group queue jobs by provider
			sortJobs(jobList)
		}
	}
}

func (handler DataTickerHandler) receiveNewJobs(attempts int) (*[]jobs.JobQueueItem, int) {
	handler.LogService.Infoln("attempt to receive jobs...")
	jobList, err := handler.QueueService.ReceiveJobs()

	handler.LogService.Infow("jobs TEMP", "joblist", jobList)

	if err != nil {
		attempts += 1
		handler.LogService.Warnw("Failed to get queue items",
			"attempts", attempts,
			"error", err,
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
