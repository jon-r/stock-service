package db

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"jon-richards.com/stock-app/internal/providers"
)

type TickerDescription struct {
	FullName string
	Currency string
	Icon     string
}

type TickerPrices struct {
	Open      float32
	Close     float32
	High      float32
	Average   float32
	Low       float32
	Timestamp int64
}

type TickerItem struct {
	TickerId    string
	Provider    providers.ProviderName
	Description TickerDescription
	Prices      []TickerPrices

	UpdatedAt int64
}

var tableName = os.Getenv("DB_TICKERS_TABLE_NAME")

func (db DatabaseRepository) NewTickerItem(provider providers.ProviderName, tickerId string) error {
	var err error

	// todo check if tickerId already exists and error if it does (maybe remove updated in cdk?)
	ticker := TickerItem{
		TickerId:  tickerId,
		Provider:  provider,
		UpdatedAt: time.Now().UnixMilli(),
	}
	av, err := attributevalue.MarshalMap(ticker)

	if err != nil {
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	}

	_, err = db.svc.PutItem(context.TODO(), &input)

	return err
}
