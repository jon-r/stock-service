package providers_old

import (
	"github.com/polygon-io/client-go/rest/models"
)

type mockProviderService struct{}

func (p *mockProviderService) FetchTickerDescription(provider ProviderName, tickerId string) (*TickerDescription, error) {
	return &TickerDescription{
		FullName:   "Full name " + tickerId,
		FullTicker: "Ticker:" + tickerId,
		Currency:   "GBP",
		Icon:       "Icon:" + string(provider) + "/" + tickerId,
	}, nil
}

func (p *mockProviderService) FetchTickerHistoricalPrices(provider ProviderName, tickerId string) (*[]TickerPrices, error) {
	return &[]TickerPrices{
		{
			Id:        tickerId + ":" + string(provider),
			Open:      10,
			Close:     20,
			High:      30,
			Low:       5,
			Timestamp: models.Millis{},
		},
		{
			Id:        tickerId + ":" + string(provider),
			Open:      20,
			Close:     30,
			High:      35,
			Low:       15,
			Timestamp: models.Millis{},
		},
	}, nil
}

func (p *mockProviderService) FetchTickerDailyPrices(provider ProviderName, tickerIds []string) (*[]TickerPrices, error) {
	var prices []TickerPrices

	for _, tickerId := range tickerIds {
		prices = append(prices, TickerPrices{
			Id:        tickerId + ":" + string(provider),
			Open:      40,
			Close:     50,
			High:      55,
			Low:       35,
			Timestamp: models.Millis{},
		})
	}

	return &prices, nil
}

func NewMockProviderService() ProviderService {
	return &mockProviderService{}
}
