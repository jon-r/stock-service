package db

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"jon-richards.com/stock-app/internal/providers"
)

var tableName = aws.String(os.Getenv("DB_TICKERS_TABLE_NAME"))

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
		TableName: tableName,
		Item:      av,
	}

	_, err = db.svc.PutItem(context.TODO(), &input)

	return err
}

func (db DatabaseRepository) SetTickerItemValue(tickerId string, name string, value interface{}) error {
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
		TableName:                 tableName,
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, err = db.svc.UpdateItem(context.TODO(), &input)

	return err
}

func (db DatabaseRepository) AddTickerItemValue(tickerId string, name string, value interface{}) error {
	var err error

	key, err := attributevalue.MarshalMap(map[string]string{"TickerId": tickerId})

	if err != nil {
		return err
	}

	update := expression.Add(expression.Name(name), expression.Value(value))
	update.Set(expression.Name("UpdatedAt"), expression.Value(time.Now().UnixMilli()))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return err
	}

	input := dynamodb.UpdateItemInput{
		TableName:                 tableName,
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, err = db.svc.UpdateItem(context.TODO(), &input)

	return err
}

func (db DatabaseRepository) SetTickerDescription(tickerId string, description *providers.TickerDescription) error {
	return db.SetTickerItemValue(tickerId, "Description", *description)
}

func (db DatabaseRepository) SetTickerHistoricalPrices(tickerId string, prices *[]providers.TickerPrices) error {
	return db.SetTickerItemValue(tickerId, "Prices", *prices)
}

func (db DatabaseRepository) UpdateTickerDailyPrices(tickerId string, prices *providers.TickerPrices) error {
	return db.AddTickerItemValue(tickerId, "Prices", *prices)
}

func (db DatabaseRepository) GetAllTickers() ([]providers.TickerItemStub, error) {
	var tickers []providers.TickerItemStub
	var err error
	var response *dynamodb.ScanOutput

	projEx := expression.NamesList(
		expression.Name("TickerId"), expression.Name("Provider"),
	)
	expr, err := expression.NewBuilder().WithProjection(projEx).Build()

	if err != nil {
		return nil, err
	}

	scanPaginator := dynamodb.NewScanPaginator(db.svc, &dynamodb.ScanInput{
		TableName:                 tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})
	for scanPaginator.HasMorePages() {
		response, err = scanPaginator.NextPage(context.TODO())
		if err != nil {
			break
		} else {
			var tickerPage []providers.TickerItemStub
			err = attributevalue.UnmarshalListOfMaps(response.Items, &tickerPage)

			if err != nil {
				break
			} else {
				tickers = append(tickers, tickerPage...)
			}
		}
	}

	return tickers, err
}
