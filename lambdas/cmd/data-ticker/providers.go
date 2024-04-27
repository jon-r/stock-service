package main

import (
	"context"
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
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

func invokeWorkerTicker(ctx context.Context, provider providers.ProviderName, delay providers.SettingsDelay) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

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
					log.Infow("Invoking Job",
						"job", job,
					)
					err = eventsService.InvokeWorker(job.Action)
					if err != nil {
						log.Warnw("Failed to Invoke Worker",
							"error", err,
						)

						err = queueService.RetryJob(job.Action, err.Error())
					}

					err = queueService.DeleteJob(job.RecieptHandle)
					if err != nil {
						log.Warnw("Failed to delete Job from queue",
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
