package providers

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/models/prices"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/array"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

type p struct {
	clock  clock.Clock
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

func (p *p) getHistoricalPrices(tickerId string) (*[]prices.TickerPrices, error) {
	params := models.ListAggsParams{
		Ticker:     tickerId,
		Multiplier: 1,
		Timespan:   models.Day,
		From:       historyStart,
		To:         models.Millis(p.clock.Now()),
	}.WithOrder(models.Desc).WithAdjusted(true)

	iter := p.client.ListAggs(context.TODO(), params)

	var pricesList []prices.TickerPrices

	for iter.Next() {
		item := iter.Item()

		pricesList = append(pricesList, p.aggregateToPrice(item, tickerId))
	}

	return &pricesList, iter.Err()
}

func (p *p) getDailyPrices(tickerIds []string) (*[]prices.TickerPrices, error) {
	yesterday := models.Date(p.clock.Now().AddDate(0, 0, -1))

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

	var pricesList []prices.TickerPrices

	for _, tickerId := range tickerIds {
		item, exists := array.Find(res.Results, func(price models.Agg) bool {
			return price.Ticker == tickerId
		})

		if exists {
			pricesList = append(pricesList, p.aggregateToPrice(item, tickerId))
		}
	}

	return &pricesList, nil
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

func newPolygonAPI(httpClient *http.Client, c clock.Clock) API {
	return &p{
		clock:  c,
		client: polygon.NewWithClient(os.Getenv("POLYGON_API_KEY"), httpClient),
	}
}
