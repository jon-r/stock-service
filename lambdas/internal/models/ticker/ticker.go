package ticker

import (
	"os"
)

func NewTickerEntity(params NewTickerParams) *Entity {
	entity := &Entity{
		Provider: params.Provider,
	}
	entity.SetKey(KeyTicker, params.TickerId, KeyTickerId, params.TickerId)

	return entity
}

func (t *Entity) TableName() string {
	return os.Getenv("STOCK_TICKER_TABLE_NAME")
}
