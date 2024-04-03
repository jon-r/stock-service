package db

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
	"jon-richards.com/stock-app/internal/providers"
)

var tableName = aws.String(os.Getenv("DB_TICKERS_TABLE_NAME"))

func (db DatabaseRepository) NewTickerItem(log *zap.SugaredLogger, params providers.NewTickerParams) error {
	var err error

	ticker := TickerItem{
		Provider: params.Provider,
	}
	ticker.SetKey(KeyTicker, params.TickerId, KeyTickerId, params.TickerId)

	av, err := attributevalue.MarshalMap(ticker)

	log.Infow("add item",
		"original", ticker,
		"item", av,
	)

	if err != nil {
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: db.StocksTableName,
		Item:      av,
	}

	_, err = db.svc.PutItem(context.TODO(), &input)

	return err
}

func (db DatabaseRepository) SetTickerDescription(log *zap.SugaredLogger, tickerId string, description *providers.TickerDescription) error {
	var err error

	var item = StocksTableItem{}
	item.SetKey(KeyTicker, tickerId, KeyTickerId, tickerId)

	if err != nil {
		return err
	}

	update := expression.Set(expression.Name("Description"), expression.Value(*description))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return err
	}

	input := dynamodb.UpdateItemInput{
		TableName:                 tableName,
		Key:                       item.GetKey(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}

	log.Infow("Update item",
		"item", item,
		"key", item.GetKey(),
		"input", input,
	)

	_, err = db.svc.UpdateItem(context.TODO(), &input)

	return err
}

//func (db DatabaseRepository) AddTickerItemValue(tickerId string, name string, value interface{}) error {
//	var err error
//
//	key, err := attributevalue.MarshalMap(map[string]string{"TickerId": tickerId})
//
//	if err != nil {
//		return err
//	}
//
//	update := expression.Add(expression.Name(name), expression.Value(value))
//	// update.Set(expression.Name("UpdatedAt"), expression.Value(time.Now().UnixMilli()))
//	expr, err := expression.NewBuilder().WithUpdate(update).Build()
//
//	if err != nil {
//		return err
//	}
//
//	input := dynamodb.UpdateItemInput{
//		TableName:                 tableName,
//		Key:                       key,
//		ExpressionAttributeNames:  expr.Names(),
//		ExpressionAttributeValues: expr.Values(),
//		UpdateExpression:          expr.Update(),
//	}
//	_, err = db.svc.UpdateItem(context.TODO(), &input)
//
//	return err
//}

//func (db DatabaseRepository) SetTickerDescription(tickerId string, description providers.TickerDescription) error {
//	return db.SetTickerItemValue(tickerId, "Description", description)
//}

//func (db DatabaseRepository) SetTickerHistoricalPrices(tickerId string, prices []providers.TickerPrices) error {
//	// todo STK-96 cant ADD to map :(
//	//  redo this with binary set instead of map (this feels like best option) <- [][]byte will cnvert to binary set
//	//     https://www.golinuxcloud.com/golang-base64-encode/
//	//  alternatively read, then set?
//	//  OR whole new table for prices?
//	return db.SetTickerItemValue(tickerId, "Prices", prices)
//}

//func (db DatabaseRepository) UpdateTickerDailyPrices(tickerId string, prices []providers.TickerPrices) error {
//	return db.AddTickerItemValue(tickerId, "Prices", prices)
//}

//func (db DatabaseRepository) GetAllTickers() ([]providers.TickerItemStub, error) {
//	var tickers []providers.TickerItemStub
//	var err error
//	var response *dynamodb.ScanOutput
//
//	projEx := expression.NamesList(
//		expression.Name("TickerId"), expression.Name("Provider"),
//	)
//	expr, err := expression.NewBuilder().WithProjection(projEx).Build()
//
//	if err != nil {
//		return nil, err
//	}
//
//	scanPaginator := dynamodb.NewScanPaginator(db.svc, &dynamodb.ScanInput{
//		TableName:                 tableName,
//		ExpressionAttributeNames:  expr.Names(),
//		ExpressionAttributeValues: expr.Values(),
//		FilterExpression:          expr.Filter(),
//		ProjectionExpression:      expr.Projection(),
//	})
//	for scanPaginator.HasMorePages() {
//		response, err = scanPaginator.NextPage(context.TODO())
//		if err != nil {
//			break
//		} else {
//			var tickerPage []providers.TickerItemStub
//			err = attributevalue.UnmarshalListOfMaps(response.Items, &tickerPage)
//
//			if err != nil {
//				break
//			} else {
//				tickers = append(tickers, tickerPage...)
//			}
//		}
//	}
//
//	return tickers, err
//}
