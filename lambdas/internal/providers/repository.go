package providers

import "fmt"

func FetchTickerDescription(provider ProviderName, tickerId string) (*TickerDescription, error) {
	switch provider {
	case PolygonIo:
		return fetchPolygonTickerDescription(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}

func FetchTickerHistoricalPrices(provider ProviderName, tickerId string) (*[]TickerPrices, error) {
	switch provider {
	case PolygonIo:
		return fetchPolygonTickerPrices(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}

func FetchTickerDailyPrices(provider ProviderName, tickerIds []string) (*map[string]TickerPrices, error) {
	switch provider {
	case PolygonIo:
		return fetchPolygonDailyPrices(tickerIds)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}
