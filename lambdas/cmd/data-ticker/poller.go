package main

import (
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/models/job"
)

func (h *handler) pollJobsQueue() {
	var err error
	var jobList *[]job.Job

	emptyResponses := 0
	failedAttempts := 0

	ticker := h.Clock.Ticker(10 * time.Second)

	h.Log.Debug("begin poll jobs queue")

	for {
		select {
		case <-h.done:
			h.Log.Debug("finished polling jobs queue")
			ticker.Stop()
			return
		case <-ticker.C:
			h.Log.Debugw("tick!")

			// 1. poll to get all items in queue
			jobList, err = h.Jobs.ReceiveJobs()

			// if queue errors too many times , disable the event rule and stop the ticker
			if err != nil {
				h.Log.Warnw("failed to receive jobs", "err", err)
				failedAttempts++

				if failedAttempts > 5 {
					h.Log.Errorf("aborting after %d failed attempts", failedAttempts)
					h.Jobs.StopScheduledRule()
					h.done <- true
				}
			} else {
				failedAttempts = 0
			}

			// 3. if queue is empty too many times, disable the event rule
			if len(*jobList) == 0 {
				h.Log.Debug("no jobs received")
				emptyResponses++

				if emptyResponses == 6 {
					h.Log.Info("no new jobs received in 60 seconds, disabling scheduler")
					h.Jobs.StopScheduledRule()
				}
			} else {
				emptyResponses = 0
			}

			h.addJobsToQueues(jobList)
		}
	}

}
