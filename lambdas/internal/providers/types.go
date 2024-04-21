package providers

import "github.com/polygon-io/client-go/rest/models"

type ProviderName string

const (
	PolygonIo ProviderName = "POLYGON_IO"
)

type TickerDescription struct {
	FullName   string
	FullTicker string
	Currency   string
	Icon       string
}

type TickerPrices struct {
	Id        string
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Timestamp models.Millis
}

type NewTickerParams struct {
	Provider ProviderName `json:"provider"`
	TickerId string       `json:"ticker"`
}

type TickerItemStub struct {
	TickerSort string `dynamodbav:"SK"`
	Provider   ProviderName
}
