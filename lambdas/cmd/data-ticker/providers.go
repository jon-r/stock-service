package main

import (
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

func (h *handler) addJobsToQueues(jobList *[]job.Job) {
	for _, j := range *jobList {
		h.queueManager.queues[j.Provider] <- j
	}
}

func (h *handler) pollProviderQueue(providerName provider.Name) {
	interval := provider.GetRequestsPerMin()[providerName]
	ticker := h.Clock.Ticker(interval)

	for {
		select {
		case <-ticker.C:
			h.invokeNextJob(providerName)
		}
	}
}

func (h *handler) invokeNextJob(providerName provider.Name) {
	select {
	case j, ok := <-h.queueManager.queues[providerName]:
		h.Log.Debugw("tock!")
		if ok {
			h.Log.Debugw("processing job", "job", j)
			// todo send this error back to the handler
			h.Jobs.InvokeWorker(j)
		} // else no jobs
	}
}
