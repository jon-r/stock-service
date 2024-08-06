package main

import (
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/jobs"
)

func (handler DataTickerHandler) checkForJobs() {
	queueTicker := handler.Clock.Ticker(10 * time.Second)

	var jobList *[]jobs.JobQueueItem

	emptyResponses := 0

	jobList, attempts := handler.receiveNewJobs(0)
	allocateJobs(jobList)

	handler.LogService.Infoln("Started polling...")

	for {
		select {
		case <-handler.done:
			handler.LogService.Infoln("Finished polling")
			queueTicker.Stop()
			return
		case <-queueTicker.C:
			// 1. poll to get all items in queue
			jobList, attempts = handler.receiveNewJobs(attempts)

			// 2. if queue is empty, disable the event rule and end the function
			emptyResponses = handler.shutDownWhenEmpty(jobList, emptyResponses)

			// 3. group queue jobs by provider
			allocateJobs(jobList)
		}
	}
}

func (handler DataTickerHandler) receiveNewJobs(attempts int) (*[]jobs.JobQueueItem, int) {
	handler.LogService.Infof("attempt to receive jobs...")
	jobList, err := handler.QueueService.ReceiveJobs()

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
