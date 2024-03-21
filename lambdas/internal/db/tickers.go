package db

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"jon-richards.com/stock-app/internal/providers"
)

var tableName = os.Getenv("DB_TICKERS_TABLE_NAME")

func (db DatabaseRepository) NewTickerItem(provider providers.ProviderName, tickerId string) error {
	var err error

	ticker := providers.TickerItem{
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

func (db DatabaseRepository) UpdateTickerItem(tickerId string, name string, value interface{}) error {
	var err error

	key, err := attributevalue.MarshalMap(map[string]string{"TickerId": tickerId})

	if err != nil {
		return err
	}

	update := expression.Set(expression.Name(name), expression.Value(value))
	update.Set(expression.Name("UpdatedAt"), expression.Value(time.Now().UnixMilli()))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return err
	}

	input := dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, err = db.svc.UpdateItem(context.TODO(), &input)

	return err
}

func (db DatabaseRepository) SetTickerDescription(tickerId string, description *providers.TickerDescription) error {
	return db.UpdateTickerItem(tickerId, "Description", *description)
}

func (db DatabaseRepository) SetTickerHistoricalPrices(tickerId string, prices *[]providers.TickerPrices) error {
	return db.UpdateTickerItem(tickerId, "Prices", *prices)
}
