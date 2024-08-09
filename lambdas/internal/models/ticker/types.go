package ticker

import (
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

type Description struct {
	FullName   string
	FullTicker string
	Currency   string
	Icon       string
}

type Entity struct {
	db.EntityBase
	Provider    provider.Name
	Description Description
}

type EntityStub struct {
	TickerSort string `dynamodbav:"SK"`
	Provider   provider.Name
}

const (
	KeyTicker   db.KeyType = "T#"
	KeyTickerId db.KeyType = "T#"
)

type NewTickerParams struct {
	Provider provider.Name `json:"provider"`
	TickerId string        `json:"ticker"`
}
