package main

import (
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/models/job"
)

func (h *handler) pollJobsQueue() {
	var err error
	var jobList *[]job.Job
	//var done chan bool

	emptyResponses := 0
	failedAttempts := 0

	ticker := h.clock.Ticker(10 * time.Second)

	for {
		select {
		case <-h.done:
			h.Log.Debug("finished polling jobs queue")
			ticker.Stop()
			return
		case <-ticker.C:
			jobList, err = h.Jobs.ReceiveJobs()

			if err != nil {
				h.Log.Warnw("failed to receive jobs", "err", err)
				failedAttempts++

				if failedAttempts > 5 {
					h.Log.Errorf("aborting after %d failed attempts", failedAttempts)
					h.done <- true
				}
			} else {
				failedAttempts = 0
			}

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

	//jobs, err := h.Jobs.ReceiveJobs()
	//if err != nil {
	//	pollAttempts += 1
	//}
}
