package providers

import (
	"context"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"os"
)

var client = polygon.New(os.Getenv("POLYGON_API_KEY"))

func GetPolygonIndexDetails(tickerId string) (error, *IndexDetails) {
	params := models.GetTickerDetailsParams{
		Ticker: tickerId,
	}

	res, err := client.GetTickerDetails(context.TODO(), &params)

	if err != nil {
		return err, nil
	}
	details := IndexDetails{
		Currency: res.Results.CurrencyName,
		FullName: res.Results.Name,
	}

	return nil, &details
}

// https://github.com/polygon-io/client-go/blob/master/rest/example/stocks/ticker-details/main.go
