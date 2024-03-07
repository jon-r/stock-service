package main

import (
	"fmt"

	"jon-richards.com/stock-app/internal/db"
)

func handleJobAction(job db.JobInput) error {
	switch job.Type {
	case db.NewStockItem:
		return newStock(job.Provider, job.TickerId)
	default:
		return fmt.Errorf("invalid action type = %v", job.Type)
	}
}
