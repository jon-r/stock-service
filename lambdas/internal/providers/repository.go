package providers

import "fmt"

// todo redo other services like this
type ProviderService interface {
	FetchTickerDescription(provider ProviderName, tickerId string) (*TickerDescription, error)
	FetchTickerHistoricalPrices(provider ProviderName, tickerId string) (*[]TickerPrices, error)
	FetchTickerDailyPrices(provider ProviderName, tickerIds []string) (*[]TickerPrices, error)
}

type providerService struct{}

func (p *providerService) FetchTickerDescription(provider ProviderName, tickerId string) (*TickerDescription, error) {
	switch provider {
	case PolygonIo:
		return fetchPolygonTickerDescription(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}

func (p *providerService) FetchTickerHistoricalPrices(provider ProviderName, tickerId string) (*[]TickerPrices, error) {
	switch provider {
	case PolygonIo:
		return fetchPolygonTickerPrices(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}

func (p *providerService) FetchTickerDailyPrices(provider ProviderName, tickerIds []string) (*[]TickerPrices, error) {
	switch provider {
	case PolygonIo:
		return fetchPolygonDailyPrices(tickerIds)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", provider)
	}
}

func NewProviderService() ProviderService {
	return &providerService{}
}
