package prices

import (
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/polygon-io/client-go/rest/models"
)

type TickerPrices struct {
	Id        string
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Timestamp models.Millis
}

const (
	KeyTickerPrice db.KeyType = "P#"
)
