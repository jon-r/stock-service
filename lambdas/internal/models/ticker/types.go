package ticker

import (
	"strings"

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

func (e *Entity) GetTickerId() string {
	tickerId, _ := strings.CutPrefix(e.EntityBase.Sort, string(KeyTickerId))
	return tickerId
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
