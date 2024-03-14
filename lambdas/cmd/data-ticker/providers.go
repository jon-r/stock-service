package main

import (
	"log"
	"time"

	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/providers"
)

var providerQueues = map[providers.ProviderName]chan jobs.JobQueueItem{
	providers.PolygonIo: make(chan jobs.JobQueueItem, 20),
}

func sortJobs(jobList *[]jobs.JobQueueItem) {
	for _, job := range *jobList {
		providerQueues[job.Action.Provider] <- job
	}
}

func invokeWorkerTicker(provider providers.ProviderName, delay providers.SettingsDelay) {
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
					err = eventsService.InvokeWorker(job.Action)
					if err != nil {
						log.Printf("Failed to Invoke Worker = %v", err)

						updatedJob := job.Action
						updatedJob.Attempts += 1

						// put the failed item back into the queue
						err = queueService.AddJobs([]jobs.JobAction{updatedJob})
					}

					err = queueService.DeleteJob(job.RecieptHandle)
					if err != nil {
						log.Printf("Failed to delete Job from queue = %v", err)
					}
				}
			default:
				// no jobs for this provider
			}
		}
	}
}
