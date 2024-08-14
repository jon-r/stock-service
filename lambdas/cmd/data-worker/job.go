package main

import (
	"fmt"
	"strings"

	"github.com/jon-r/stock-service/lambdas/internal/models/job"
)

func (h *handler) doJob(j job.Job) error {
	h.log.Debugw("attempt to do job", "job", j)
	switch j.Type {
	case job.LoadTickerDescription:
		return h.tickers.LoadDescription(j.Provider, j.TickerId)
	case job.LoadHistoricalPrices:
		return h.prices.LoadHistoricalPrices(j.Provider, j.TickerId)
	case job.LoadDailyPrices:
		return h.prices.LoadDailyPrices(j.Provider, strings.Split(j.TickerId, ","))
	// TODO STK-86
	// jobs.LoadTickerIcon

	// TODO STK-88
	// jobs.LoadDailyDividends
	// jobs.LoadHistoricalDividends

	default:
		return fmt.Errorf("invalid action type = %v", j.Type)
	}
}
