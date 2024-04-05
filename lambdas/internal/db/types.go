package db

import "jon-richards.com/stock-app/internal/providers"

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
	Description providers.TickerDescription `dynomodbav:",omitemptyelem"`
}
