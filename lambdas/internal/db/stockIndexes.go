package db

import (
	"jon-richards.com/stock-app/internal/providers"
)

// todo this file to be replaced with tickers.go

//var stockTableName = os.Getenv("DB_TICKERS_TABLE_NAME")

//type StockItem_OLD struct {
//	StockIndexId string
//	Name         string
//	Provider     providers.ProviderName
//	Image        string
//	UpdatedAt    int64
//}

type StockItemProperties struct {
	FullName string
	Currency string
	// Icon string todo
}

type StockPrices struct {
	Open      float32
	Close     float32
	High      float32
	Average   float32
	Low       float32
	Timestamp int64
}

type StockItem struct {
	StockIndexId string
	Provider     providers.ProviderName
	Properties   StockItemProperties
	Prices       []StockPrices
	UpdatedAt    int64
}

//func (db DatabaseRepository) GetItems() ([]StockItem_OLD, error) {
//	var err error
//
//	projEx := expression.NamesList(
//		expression.Name("StockIndexId"), expression.Name("Name"), expression.Name("Provider"))
//
//	expr, err := expression.NewBuilder().WithProjection(projEx).Build()
//
//	// todo need to paginate? or loop till all found
//	// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go#L248
//	input := dynamodb.QueryInput{
//		TableName:            &stockTableName,
//		ProjectionExpression: expr.Projection(),
//	}
//
//	result, err := db.svc.Query(context.TODO(), &input)
//
//	if err != nil {
//		return nil, err
//	}
//
//	var items []StockItem_OLD
//	err = attributevalue.UnmarshalListOfMaps(result.Items, items)
//	if err != nil {
//		return nil, err
//	}
//
//	return items, nil
//}

//func (db DatabaseRepository) NewStockItem(provider providers.ProviderName, tickerId string, properties StockItemProperties) error {
//	var err error
//
//	stock := StockItem{
//		StockIndexId: tickerId,
//		Properties:   properties,
//		Provider:     provider,
//		Prices:       nil,
//		UpdatedAt:    time.Now().UnixMilli(),
//	}
//	av, err := attributevalue.MarshalMap(stock)
//
//	if err != nil {
//		return err
//	}
//
//	input := dynamodb.PutItemInput{
//		TableName: &stockTableName,
//		Item:      av,
//	}
//
//	_, err = db.svc.PutItem(context.TODO(), &input)
//
//	return err
//}

// todo this should be renamed to update stock item based
//func (db DatabaseRepository) UpsertStockItem(res *providers.DogApiRes, jobItem *JobItem) error {
//	var err error
//
//	// todo will need to have the ID as part of the job db item, so it can be updated instead of created new!!
//	stock := StockItem_OLD{
//		StockIndexId: uuid.NewString(),
//		Name:         jobItem.Payload["Name"],
//		Provider:     jobItem.Provider,
//		Image:        res.Message,
//		UpdatedAt:    time.Now().Unix(),
//	}
//	av, err := attributevalue.MarshalMap(stock)
//
//	if err != nil {
//		return err
//	}
//
//	input := dynamodb.PutItemInput{
//		TableName: &stockTableName,
//		Item:      av,
//	}
//
//	_, err = db.svc.PutItem(context.TODO(), &input)
//
//	return err
//}
