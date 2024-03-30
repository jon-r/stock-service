package jobs

import "jon-richards.com/stock-app/internal/providers"

type JobTypes string

const (
	LoadTickerDescription   JobTypes = "LOAD_TICKER_DESCRIPTION"
	LoadHistoricalPrices    JobTypes = "LOAD_HISTORICAL_PRICES"
	LoadHistoricalDividends JobTypes = "LOAD_HISTORICAL_DIVIDENDS"

	LoadTickerIcon  JobTypes = "LOAD_TICKER_ICON"
	UpdatePrices    JobTypes = "UPDATE_PRICES"
	UpdateDividends JobTypes = "UPDATE_DIVIDENDS"
	// ???
)

type JobAction struct {
	JobId    string
	Provider providers.ProviderName
	Type     JobTypes
	TickerId string
	Attempts int
}

type JobErrorItem struct {
	JobAction
	ErrorReason string
}

type JobQueueItem struct {
	RecieptHandle string
	Action        JobAction
}
