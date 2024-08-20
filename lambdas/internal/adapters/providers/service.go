package providers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/models/prices"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/polygon-io/client-go/rest/models"
)

type Service interface {
	GetDescription(providerName provider.Name, tickerId string) (*ticker.Description, error)
	GetHistoricalPrices(providerName provider.Name, tickerId string) (*[]prices.TickerPrices, error)
	GetDailyPrices(providerName provider.Name, tickerIds []string) (*[]prices.TickerPrices, error)
}

type api struct {
	polygon API
}

func (api *api) GetDescription(providerName provider.Name, tickerId string) (*ticker.Description, error) {
	switch providerName {
	case provider.PolygonIo:
		return api.polygon.getDescription(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", providerName)
	}
}

func (api *api) GetHistoricalPrices(providerName provider.Name, tickerId string) (*[]prices.TickerPrices, error) {
	switch providerName {
	case provider.PolygonIo:
		return api.polygon.getHistoricalPrices(tickerId)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", providerName)
	}
}

func (api *api) GetDailyPrices(providerName provider.Name, tickerIds []string) (*[]prices.TickerPrices, error) {
	switch providerName {
	case provider.PolygonIo:
		return api.polygon.getDailyPrices(tickerIds)
	default:
		return nil, fmt.Errorf("incorrect provider name: %v", providerName)
	}
}

// free accounts won't be older than 2 years, so wont get all this
var startDate, _ = time.Parse(time.DateOnly, "2021-12-01")
var historyStart = models.Millis(startDate)

func NewService(httpClient *http.Client, c clock.Clock) Service {
	return &api{
		polygon: newPolygonAPI(httpClient, c),
	}
}
