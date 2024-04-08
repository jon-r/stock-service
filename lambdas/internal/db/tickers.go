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

// todo maybe move this elsewhere? also can the generic be added?
func mapPricesToStockItems(prices []providers.TickerPrices, tickerId string) []PriceItem {
	priceItems := make([]PriceItem, len(prices))

	for i, price := range prices {
		date, _ := price.Timestamp.MarshalJSON()
		priceItem := PriceItem{
			Price: price,
			Date:  string(date),
		}
		priceItem.SetKey(KeyTicker, tickerId, KeyTickerPrice, string(date))

		priceItems[i] = priceItem
	}

	return priceItems
}

func (db DatabaseRepository) SetTickerHistoricalPrices(log *zap.SugaredLogger, tickerId string, prices *[]providers.TickerPrices) error {
	var err error
	var item map[string]types.AttributeValue

	priceItems := mapPricesToStockItems(*prices, tickerId)

	written := 0
	batchSize := 25
	start := 0
	end := start + batchSize

	for start < len(priceItems) {
		var writeReqs []types.WriteRequest
		if end > len(priceItems) {
			end = len(priceItems)
		}
		for _, price := range priceItems[start:end] {
			item, err = attributevalue.MarshalMap(price)
			if err != nil {
				log.Warnw("Couldn't marshal item for batch writing",
					"item", price,
					"error", err,
				)
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
			log.Warnw("Couldn't add a batch of table tableItems to the table",
				"table", *db.StocksTableName,
				"req", writeReqs,
				"error", err,
			)
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	if written > 0 {
		log.Infof("Inserted %d tableItems to table %s", written, *db.StocksTableName)
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
