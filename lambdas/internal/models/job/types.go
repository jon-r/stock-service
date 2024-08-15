package job

import (
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

type Types string

const (
	LoadTickerDescription   Types = "LOAD_TICKER_DESCRIPTION"
	LoadHistoricalPrices    Types = "LOAD_HISTORICAL_PRICES"
	LoadHistoricalDividends Types = "LOAD_HISTORICAL_DIVIDENDS"

	LoadTickerIcon  Types = "LOAD_TICKER_ICON"
	LoadDailyPrices Types = "UPDATE_PRICES"
	UpdateDividends Types = "UPDATE_DIVIDENDS"
	// ???
)

type Job struct {
	ReceiptId *string `json:"-"`
	JobId     string
	Provider  provider.Name
	Type      Types
	TickerId  string
	Attempts  int
}

type FailedJob struct {
	Job
	ErrorReason string
}
