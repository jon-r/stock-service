package prices

import (
	"os"

	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
)

type Entity struct {
	db.EntityBase
	Prices TickerPrices
	Date   string `dynamodbav:"DT"`
}

func NewPriceEntity(prices TickerPrices) *Entity {
	date, _ := prices.Timestamp.MarshalJSON()
	entity := &Entity{
		Prices: prices,
		Date:   string(date),
	}
	entity.SetKey(ticker.KeyTicker, prices.Id, KeyTickerPrice, string(date))

	return entity
}

func (t *Entity) TableName() string {
	return os.Getenv("STOCK_TICKER_TABLE_NAME")
}
