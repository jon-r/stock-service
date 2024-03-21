package main

import (
	"fmt"

	"jon-richards.com/stock-app/internal/jobs"
)

func handleJobAction(job jobs.JobAction) error {
	switch job.Type {
	case jobs.LoadTickerDescription:
		return setTickerDescription(job.Provider, job.TickerId)
	case jobs.LoadHistoricalPrices:
		return setTickerHistoricalPrices(job.Provider, job.TickerId)

	default:
		return fmt.Errorf("invalid action type = %v", job.Type)
	}
}

func retryFailedJob(job jobs.JobAction) error {
	job.Attempts += 1
	return queueService.AddJobs([]jobs.JobAction{job})
}
