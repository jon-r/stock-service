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

func FetchTickerDescription(provider ProviderName, tickerId string) (*TickerDescription, error) {
	switch provider {
	case PolygonIo:
		return FetchPolygonTickerDescription(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}

func FetchTickerHistoricalPrices(provider ProviderName, tickerId string) (*[]TickerPrices, error) {
	switch provider {
	case PolygonIo:
		return FetchPolygonTickerPrices(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}
