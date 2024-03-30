package providers

import (
	"context"
	"os"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/samber/lo"
)

var client = polygon.New(os.Getenv("POLYGON_API_KEY"))

func convertPolygonToPrice(item models.Agg) TickerPrices {
	return TickerPrices{
		Open:      item.Open,
		Close:     item.Close,
		High:      item.High,
		Low:       item.Low,
		Timestamp: item.Timestamp,
	}
}

func fetchPolygonTickerDescription(tickerId string) (*TickerDescription, error) {
	params := models.GetTickerDetailsParams{
		Ticker: tickerId,
	}

	res, err := client.GetTickerDetails(context.TODO(), &params)

	if err != nil {
		return nil, err
	}
	details := TickerDescription{
		Currency: res.Results.CurrencyName,
		FullName: res.Results.Name,
		Icon:     res.Results.Branding.IconURL,
	}

	return &details, nil
}

// free polygon account won't be older than 2 years, so wont get all this
var historyStart = models.Millis(time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC))

func fetchPolygonTickerPrices(tickerId string) (*[]TickerPrices, error) {
	params := models.ListAggsParams{
		Ticker:     tickerId,
		Multiplier: 1,
		Timespan:   "day",
		From:       historyStart,
		To:         models.Millis(time.Now()),
	}.WithOrder(models.Desc).WithAdjusted(true)

	iter := client.ListAggs(context.TODO(), params)

	var prices []TickerPrices

	for iter.Next() {
		item := iter.Item()

		prices = append(prices, convertPolygonToPrice(item))
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return &prices, nil
}

func fetchPolygonDailyPrices(tickerIds []string) (*map[string]TickerPrices, error) {
	params := models.GetGroupedDailyAggsParams{
		Locale:     models.US,
		MarketType: models.Stocks,
		Date:       models.Date(time.Now()),
	}.WithAdjusted(true)

	res, err := client.GetGroupedDailyAggs(context.TODO(), params)

	if err != nil {
		return nil, err
	}

	if res.ResultsCount == 0 {
		return nil, nil
	}

	prices := make(map[string]TickerPrices, len(tickerIds))
	for _, tickerId := range tickerIds {
		item, exists := lo.Find(res.Results, func(price models.Agg) bool {
			return price.Ticker == tickerId
		})
		if exists {
			prices[tickerId] = convertPolygonToPrice(item)
		}
	}

	return &prices, nil
}

// https://github.com/polygon-io/client-go/blob/master/rest/example/stocks/ticker-details/main.go
