package db

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"

	"jon-richards.com/stock-app/lambdas/internal/providers"
)

var stockTableName = "stock-app_StockIndex"

type StockItem struct {
	StockIndexId string
	Name         string
	Provider     providers.ProviderName
	Image        string
	UpdatedAt    int64
}

func (db DatabaseRepository) GetItems() ([]StockItem, error) {
	var err error

	projEx := expression.NamesList(
		expression.Name("StockIndexId"), expression.Name("Name"), expression.Name("Provider"))

	expr, err := expression.NewBuilder().WithProjection(projEx).Build()

	// todo need to paginate? or loop till all found
	// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go#L248
	input := dynamodb.QueryInput{
		TableName:            &stockTableName,
		ProjectionExpression: expr.Projection(),
	}

	result, err := db.svc.Query(context.TODO(), &input)

	if err != nil {
		return nil, err
	}

	var items []StockItem
	err = attributevalue.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db DatabaseRepository) UpsertStockItem(res *providers.DogApiRes, jobItem *JobItem) error {
	var err error

	// todo will need to have the ID as part of the job db item, so it can be updated instead of created new!!
	stock := StockItem{
		StockIndexId: uuid.NewString(),
		Name:         jobItem.Payload["Name"],
		Provider:     jobItem.Provider,
		Image:        res.Message,
		UpdatedAt:    time.Now().Unix(),
	}
	av, err := attributevalue.MarshalMap(stock)

	if err != nil {
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: &stockTableName,
		Item:      av,
	}

	_, err = db.svc.PutItem(context.TODO(), &input)

	return err
}
