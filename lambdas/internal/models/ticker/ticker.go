package ticker

import "os"

func NewTickerEntity(params NewTickerParams) *Entity {
	entity := &Entity{
		Provider: params.Provider,
	}
	entity.SetKey(KeyTicker, params.TickerId, KeyTickerId, params.TickerId)

	return entity
}

func TableName() string {
	return os.Getenv("DB_STOCKS_TABLE_NAME")
}
