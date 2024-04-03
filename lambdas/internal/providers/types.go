package providers

import "github.com/polygon-io/client-go/rest/models"

type ProviderName string

const (
	PolygonIo ProviderName = "POLYGON_IO"
)

type TickerDescription struct {
	FullName    string
	FullTicker  string
	Currency    string
	Icon        string
	Description string
}

type TickerPrices struct {
	Open      float64
	Close     float64
	High      float64
	Average   float64 // todo remove average
	Low       float64
	Timestamp models.Millis
}

type NewTickerParams struct {
	Provider ProviderName `json:"provider"`
	TickerId string       `json:"ticker"`
}

type TickerItemStub struct {
	TickerId string
	Provider ProviderName
}
