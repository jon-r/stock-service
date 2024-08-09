package providers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/models/prices"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/array"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

type p struct {
	client *polygon.Client
}

type API interface {
	getDescription(tickerId string) (*ticker.Description, error)
	getHistoricalPrices(tickerId string) (*[]prices.TickerPrices, error)
	getDailyPrices(tickerIds []string) (*[]prices.TickerPrices, error)
}

func (p *p) getDescription(tickerId string) (*ticker.Description, error) {
	params := models.GetTickerDetailsParams{
		Ticker: tickerId,
	}

	res, err := p.client.GetTickerDetails(context.TODO(), &params)

	if err != nil {
		return nil, err
	}
	description := ticker.Description{
		Currency:   res.Results.CurrencyName,
		FullName:   res.Results.Name,
		FullTicker: strings.Join([]string{res.Results.PrimaryExchange, res.Results.Ticker}, ":"),
		Icon:       res.Results.Branding.IconURL,
	}

	return &description, nil
}

// free polygon account won't be older than 2 years, so wont get all this
var historyStart = models.Millis(time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC))

func (p *p) getHistoricalPrices(tickerId string) (*[]prices.TickerPrices, error) {
	params := models.ListAggsParams{
		Ticker:     tickerId,
		Multiplier: 1,
		Timespan:   models.Day,
		From:       historyStart,
		To:         models.Millis(time.Now()),
	}.WithOrder(models.Desc).WithAdjusted(true)

	iter := p.client.ListAggs(context.TODO(), params)

	var prices []prices.TickerPrices

	for iter.Next() {
		item := iter.Item()

		prices = append(prices, p.aggregateToPrice(item, tickerId))
	}

	return &prices, iter.Err()
}

func (p *p) getDailyPrices(tickerIds []string) (*[]prices.TickerPrices, error) {
	yesterday := models.Date(time.Now().AddDate(0, 0, -1))

	params := models.GetGroupedDailyAggsParams{
		Locale:     models.US,
		MarketType: models.Stocks,
		Date:       yesterday,
	}.WithAdjusted(true)

	res, err := p.client.GetGroupedDailyAggs(context.TODO(), params)

	if err != nil {
		return nil, err
	}

	if res.ResultsCount == 0 {
		return nil, nil
	}

	var prices []prices.TickerPrices

	for _, tickerId := range tickerIds {
		item, exists := array.Find(res.Results, func(price models.Agg) bool {
			return price.Ticker == tickerId
		})

		if exists {
			prices = append(prices, p.aggregateToPrice(item, tickerId))
		}
	}

	return &prices, nil
}

func (p *p) aggregateToPrice(item models.Agg, tickerId string) prices.TickerPrices {
	return prices.TickerPrices{
		Id:        tickerId,
		Open:      item.Open,
		Close:     item.Close,
		High:      item.High,
		Low:       item.Low,
		Timestamp: item.Timestamp,
	}
}

func newPolygonAPI(httpClient *http.Client) API {
	return &p{
		client: polygon.NewWithClient(os.Getenv("POLYGON_API_KEY"), httpClient),
	}
}
