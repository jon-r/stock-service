package providers

import (
	"github.com/jon-r/stock-service/lambdas/internal/models/prices"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/polygon-io/client-go/rest/models"
)

type mockProviderApi struct{}

func (api *mockProviderApi) GetDescription(providerName provider.Name, tickerId string) (*ticker.Description, error) {
	return &ticker.Description{
		FullName:   "Full name " + tickerId,
		FullTicker: "Ticker:" + tickerId,
		Currency:   "GBP",
		Icon:       "Icon:" + string(providerName) + "/" + tickerId,
	}, nil
}

func (api *mockProviderApi) GetHistoricalPrices(_ provider.Name, tickerId string) (*[]prices.TickerPrices, error) {
	return &[]prices.TickerPrices{
		{
			Id:        tickerId,
			Open:      10,
			Close:     20,
			High:      30,
			Low:       5,
			Timestamp: models.Millis{},
		},
		{
			Id:        tickerId,
			Open:      20,
			Close:     30,
			High:      35,
			Low:       15,
			Timestamp: models.Millis{},
		},
	}, nil
}

func (api *mockProviderApi) GetDailyPrices(_ provider.Name, tickerIds []string) (*[]prices.TickerPrices, error) {
	var p []prices.TickerPrices

	for _, tickerId := range tickerIds {
		p = append(p, prices.TickerPrices{
			Id:        tickerId,
			Open:      40,
			Close:     50,
			High:      55,
			Low:       35,
			Timestamp: models.Millis{},
		})
	}

	return &p, nil
}

func NewMock() Service {
	return &mockProviderApi{}
}
