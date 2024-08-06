package db

import "github.com/jon-r/stock-service/lambdas/internal/providers"

type KeyType string

const (
	KeyTicker         KeyType = "T#"
	KeyTickerPrice    KeyType = "P#"
	KeyTickerId       KeyType = "T#"
	KeyTickerDividend KeyType = "D#"

	KeyUser        KeyType = "U#"
	KeyUserTicker  KeyType = "T#"
	KeyUserTxEvent KeyType = "E#"
)

type StocksTableItem struct {
	Id   string `dynamodbav:"PK"`
	Sort string `dynamodbav:"SK"`
}

type TickerItem struct {
	StocksTableItem
	Provider    providers.ProviderName
	Description providers.TickerDescription
}

type PriceItem struct {
	StocksTableItem
	Price providers.TickerPrices
	Date  string `dynamodbav:"DT"`
}
