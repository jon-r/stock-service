package main

import (
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/providers"
)

var providerQueues = map[providers.ProviderName]chan jobs.JobQueueItem{
	providers.PolygonIo: make(chan jobs.JobQueueItem, 20),
}

func sortJobs(jobList *[]jobs.JobQueueItem) {
	for _, job := range *jobList {
		providerQueues[job.Action.Provider] <- job
	}
}

func (handler DataTickerHandler) invokeWorkerTicker(provider providers.ProviderName, delay providers.SettingsDelay) {
	var err error

	duration := time.Duration(delay) * time.Second
	ticker := time.NewTicker(duration)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			select {
			case job, ok := <-providerQueues[provider]:
				if ok {
					handler.log.Infow("Invoking Job",
						"job", job,
					)
					err = handler.eventsService.InvokeWorker(job.Action)
					if err != nil {
						handler.log.Warnw("Failed to Invoke Worker",
							"error", err,
						)

						err = handler.queueService.RetryJob(job.Action, err.Error())
					}

					err = handler.queueService.DeleteJob(job.RecieptHandle)
					if err != nil {
						handler.log.Warnw("Failed to delete Job from queue",
							"error", err,
						)
					}
				}
			default:
				// no jobs for this provider
			}
		}
	}
}
