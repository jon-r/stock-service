package main

import (
	"fmt"

	"jon-richards.com/stock-app/internal/jobs"
)

func handleJobAction(job jobs.JobAction) error {
	switch job.Type {
	case jobs.NewTickerItem:
		return setTickerDescription(job.Provider, job.TickerId)
	default:
		return fmt.Errorf("invalid action type = %v", job.Type)
	}
}
