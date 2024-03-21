package providers

import "github.com/polygon-io/client-go/rest/models"

type ProviderName string

const (
	PolygonIo ProviderName = "POLYGON_IO"
)

type TickerDescription struct {
	FullName string
	Currency string
	Icon     string
}

type TickerPrices struct {
	Open      float64
	Close     float64
	High      float64
	Average   float64
	Low       float64
	Timestamp models.Millis
}

type TickerItem struct {
	TickerId    string
	Provider    ProviderName
	Description TickerDescription
	Prices      []TickerPrices
	UpdatedAt   int64
}
