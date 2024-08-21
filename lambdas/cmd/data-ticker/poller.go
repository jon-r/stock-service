package main

import "time"

const pollInterval = 10 * time.Second

func (h *handler) pollJobsQueue() {
	h.Log.Debugln("begin poll jobs queue")
	ticker := h.Clock.Ticker(pollInterval)

	for {
		select {
		case <-h.queueManager.done:
			h.Log.Debugln("finished polling jobs queue")
			ticker.Stop()
			return
		case <-ticker.C:
			h.Log.Debugln("tick!")
			h.checkForJobs()
		}
	}
}

func (h *handler) checkForJobs() {
	// 1. poll to get all items in queue
	jobList, err := h.Jobs.ReceiveJobs()

	// 2. if queue errors too many times , disable the event rule and stop the ticker
	if err != nil {
		h.Log.Warnw("failed to receive jobs", "err", err)
		h.queueManager.failedAttempts++

		if h.queueManager.failedAttempts > 5 {
			h.Log.Errorf("aborting after %d failed attempts", h.queueManager.failedAttempts)
			h.Jobs.StopScheduledRule()
			h.queueManager.done <- true
		}

		return
	} else {
		h.queueManager.failedAttempts = 0
	}

	// 3. if queue is empty too many times, disable the event rule
	if len(*jobList) == 0 {
		h.Log.Debugln("no jobs received")
		h.queueManager.emptyResponses++

		if h.queueManager.emptyResponses == 6 {
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
