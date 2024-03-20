package providers

import (
	"fmt"

	"jon-richards.com/stock-app/internal/db"
)

type IndexDetails struct {
	TickerId string
	Currency string
	FullName string
	// icon
}

func FetchTickerDescription(provider ProviderName, tickerId string) (error, *db.TickerDescription) {
	switch provider {
	case PolygonIo:
		return FetchPolygonTickerDescription(tickerId)
	default:
		return fmt.Errorf("incorrect provider name: %v", provider), nil
	}
}

func FetchTickerHistoricalPrices(provider ProviderName, tickerId string) (error, *[]db.TickerPrices) {
	switch provider {
	case PolygonIo:
		return FetchPolygonTickerPrices(tickerId)
	default:
		return fmt.Errorf("incorrect provider name: %v", provider), nil
	}
}
