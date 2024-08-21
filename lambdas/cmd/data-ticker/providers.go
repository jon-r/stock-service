package main

import (
	"context"

	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

func (h *handler) addJobsToQueues(jobList *[]job.Job) {
	for _, j := range *jobList {
		h.queueManager.queues[j.Provider] <- j
	}
}

func (h *handler) pollProviderQueue(ctx context.Context, providerName provider.Name) {
	interval := provider.GetRequestsPerMin()[providerName]
	ticker := h.Clock.Ticker(interval)

	for {
		select {
		case <-ctx.Done():
			h.Log.Debugln("finished polling provider jobs")
			ticker.Stop()
			return
		case <-ticker.C:
			h.invokeNextJob(providerName)
		}
	}
}

func (h *handler) invokeNextJob(providerName provider.Name) {
	select {
	case j, ok := <-h.queueManager.queues[providerName]:
		if ok {
			h.Log.Debugw("processing job", "job", j)
			// todo send error back to the handler
			h.Jobs.InvokeWorker(j)
		} else {
			// else no jobs
			h.Log.Debugw("no job to process", "provider", providerName)
		}
	default:
		// nothing in queue (keep this default, else the queue channel is blocking)
	}
}
