package providers

import (
	"context"
	"os"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"jon-richards.com/stock-app/internal/db"
)

var client = polygon.New(os.Getenv("POLYGON_API_KEY"))

func FetchPolygonTickerDescription(tickerId string) (error, *db.TickerDescription) {
	params := models.GetTickerDetailsParams{
		Ticker: tickerId,
	}

	res, err := client.GetTickerDetails(context.TODO(), &params)

	if err != nil {
		return err, nil
	}
	details := db.TickerDescription{
		Currency: res.Results.CurrencyName,
		FullName: res.Results.Name,
		Icon:     res.Results.Branding.IconURL,
	}

	return nil, &details
}

var oldestHistoricalPoint = models.Millis(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

func FetchPolygonTickerPrices(tickerId string) (error, *[]db.TickerPrices) {
	params := models.ListAggsParams{
		Ticker:     tickerId,
		Multiplier: 1,
		Timespan:   "day",
		From:       oldestHistoricalPoint,
		To:         models.Millis(time.Now()),
	}.WithOrder(models.Desc).WithAdjusted(true)

	iter := client.ListAggs(context.TODO(), params)

	var prices []db.TickerPrices

	for iter.Next() {
		item := iter.Item()

		prices = append(prices, db.TickerPrices{
			Open:      item.Open,
			Close:     item.Close,
			High:      item.High,
			Low:       item.Low,
			Timestamp: item.Timestamp,
		})
	}

	if iter.Err() != nil {
		return iter.Err(), nil
	}

	return nil, &prices
}

// https://github.com/polygon-io/client-go/blob/master/rest/example/stocks/ticker-details/main.go
