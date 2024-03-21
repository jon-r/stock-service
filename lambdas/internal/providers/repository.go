package providers

import (
	"fmt"
)

type IndexDetails struct {
	TickerId string
	Currency string
	FullName string
	// icon
}

func FetchTickerDescription(provider ProviderName, tickerId string) (error, *TickerDescription) {
	switch provider {
	case PolygonIo:
		return FetchPolygonTickerDescription(tickerId)
	default:
		return fmt.Errorf("incorrect provider name: %v", provider), nil
	}
}

func FetchTickerHistoricalPrices(provider ProviderName, tickerId string) (error, *[]TickerPrices) {
	switch provider {
	case PolygonIo:
		return FetchPolygonTickerPrices(tickerId)
	default:
		return fmt.Errorf("incorrect provider name: %v", provider), nil
	}
}
