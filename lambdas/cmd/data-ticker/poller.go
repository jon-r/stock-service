package main

import (
	"context"
	"time"
)

const pollInterval = 10 * time.Second
const maxFailedAttempts = 5
const maxEmptyResponses = 6

func (h *handler) pollJobsQueue(ctx context.Context, cancel context.CancelFunc) {
	h.Log.Debugln("begin poll jobs queue")
	ticker := h.Clock.Ticker(pollInterval)

	for {
		select {
		case <-ctx.Done():
			h.Log.Debugln("finished polling jobs queue")
			ticker.Stop()
			return
		case <-ticker.C:
			h.checkForJobs(cancel)
		}
	}
}

func (h *handler) checkForJobs(cancel context.CancelFunc) {
	// 1. poll to get all items in queue
	jobList, err := h.Jobs.ReceiveJobs()

	// 2. if queue errors too many times , disable the event rule and stop the ticker
	if err != nil {
		h.Log.Warnw("failed to receive jobs", "err", err)
		h.queueManager.failedAttempts++

		if h.queueManager.failedAttempts > maxFailedAttempts {
			h.Log.Errorf("aborting after %d failed attempts", h.queueManager.failedAttempts)
			h.Jobs.StopScheduledRule()
			cancel()
			return
		}

		return
	} else {
		h.queueManager.failedAttempts = 0
	}

	// 3. if queue is empty too many times, disable the event rule
	if len(*jobList) == 0 {
		h.Log.Debugln("no jobs received")
		h.queueManager.emptyResponses++

		if h.queueManager.emptyResponses == maxEmptyResponses {
			h.Log.Infoln("no new jobs received in 60 seconds, disabling scheduler")
			h.Jobs.StopScheduledRule()
		}

		return
	} else {
		h.queueManager.emptyResponses = 0
	}

	// 4. sort the jobs by provider
	for _, j := range *jobList {
		h.queueManager.queues[j.Provider] <- j
	}
}
