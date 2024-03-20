package providers

import (
	"context"
	"os"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

var client = polygon.New(os.Getenv("POLYGON_API_KEY"))

func FetchPolygonTickerDescription(tickerId string) (error, *IndexDetails) {
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
