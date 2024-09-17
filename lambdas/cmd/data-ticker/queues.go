package main

import (
	"context"
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

const pollInterval = 10 * time.Second
const maxFailedAttempts = 5
const maxEmptyResponses = int(time.Minute / pollInterval)

func (h *handler) pollJobsQueue(ctx context.Context, cancel context.CancelFunc) {
	h.Log.Debugln("begin poll jobs queue")
	h.pollUntilCancelled(ctx, func() {
		h.checkForJobs(cancel)
	}, pollInterval)
	h.Log.Debugln("finished polling jobs queue")
}

func (h *handler) checkForJobs(cancel context.CancelFunc) {
	// 1. poll to get all items in queue
	jobList, err := h.Jobs.ReceiveJobs()

	// 2. if queue errors too many times , disable the event rule and stop the ticker
	if err != nil {
		h.Log.Warnw("failed to receive jobs", "err", err)
		h.queueManager.failedAttempts++

		if h.queueManager.failedAttempts == maxFailedAttempts {
			h.Log.Errorf("aborting after %d failed attempts", h.queueManager.failedAttempts)
			h.Jobs.StopScheduledRule()
			cancel()
		}

		return
	} else {
		h.queueManager.failedAttempts = 0
	}

	// 3. if queue is empty too many times, disable the event rule (but keep the ticker running)
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

	// 4. assign the jobs by provider
	for _, j := range *jobList {
		h.queueManager.queues[j.Provider] <- j
	}
}

func (h *handler) pollProviderQueue(ctx context.Context, providerName provider.Name) {
	h.Log.Debugln("begin poll provider queue")
	interval := provider.GetRequestsPerMin()[providerName]
	h.pollUntilCancelled(ctx, func() {
		h.invokeNextJob(providerName)
	}, interval)
	h.Log.Debugln("finished polling provider jobs")
}

func (h *handler) invokeNextJob(providerName provider.Name) {
	select {
	case j, ok := <-h.queueManager.queues[providerName]:
		if ok {
			h.Log.Debugw("processing job", "job", j)
			// todo send error back to the handler
			h.Jobs.InvokeWorker(j)
		}
		// else no jobs
	default:
		h.Log.Debugw("no job to process", "provider", providerName)
	}
}
