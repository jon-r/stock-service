package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
	"jon-richards.com/stock-app/internal/providers"
)

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
		TableName:                 db.StocksTableName,
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

func (db DatabaseRepository) SetTickerHistoricalPrices(log *zap.SugaredLogger, tickerId string, prices []providers.TickerPrices) error {
	// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go#L182

	var err error
	var item map[string]types.AttributeValue

	written := 0
	batchSize := 25
	start := 0
	end := start + batchSize
	// todo split this up to be one fn that makes the data, and one that batch inserts it
	for start < len(prices) {
		var writeReqs []types.WriteRequest
		if end > len(prices) {
			end = len(prices)
		}
		for _, price := range prices[start:end] {
			date, _ := price.Timestamp.MarshalJSON()
			priceItem := PriceItem{
				Price: price,
				Date:  string(date),
			}
			priceItem.SetKey(KeyTicker, tickerId, KeyTickerPrice, string(date))

			item, err = attributevalue.MarshalMap(price)
			if err != nil {
				log.Warnw("Couldn't marshal price for batch writing",
					"price", price.Timestamp,
					"error", err)
			} else {
				writeReqs = append(writeReqs, types.WriteRequest{
					PutRequest: &types.PutRequest{Item: item},
				})
			}
		}
		_, err = db.svc.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{*db.StocksTableName: writeReqs},
		})
		if err != nil {
			log.Warnw("Couldn't add a batch of movies to the table",
				"table", *db.StocksTableName,
				"error", err,
			)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	if written > 0 {
		log.Infof("Inserted %d items to table %s", written, *db.StocksTableName)
	}

	return err
}

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
