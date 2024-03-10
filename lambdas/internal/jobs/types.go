package jobs

import "jon-richards.com/stock-app/internal/providers"

type JobTypes string

const (
	NewTickerItem     JobTypes = "NEW_TICKER_ITEM"
	LoadTickerHistory JobTypes = "LOAD_TICKER_HISTORY"
	UpdateAllPrices   JobTypes = "UPDATE_ALL_PRICES"
	UpdateDividends   JobTypes = "UPDATE_DIVIDENDS"
	// ???
)

type JobAction struct {
	Provider providers.ProviderName
	Type     JobTypes
	TickerId string
	Attempts int
}
