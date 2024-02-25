package db

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	"jon-richards.com/stock-app/remote"
)

var stockTableName = "stock-app_StockIndex"

type StockItem struct {
	StockIndexId string
	Name         string
	Group        string
	Image        string
	UpdatedAt    string
}

func (db DatabaseRepository) GetItems() ([]StockItem, error) {
	var err error

	projEx := expression.NamesList(
		expression.Name("StockIndexId"), expression.Name("Name"), expression.Name("Group"))

	expr, err := expression.NewBuilder().WithProjection(projEx).Build()

	// todo need to paginate? or loop till all found
	input := dynamodb.QueryInput{
		TableName:            &stockTableName,
		ProjectionExpression: expr.Projection(),
	}

	result, err := db.svc.Query(&input)

	if err != nil {
		return nil, err
	}

	var items []StockItem
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db DatabaseRepository) UpsertStockItem(res *remote.DogApiRes, jobItem *JobItem) error {
	var err error

	// todo need to have the ID as part of the job db item, so it can be updated instead of created new!!
	stock := StockItem{
		StockIndexId: uuid.NewString(),
		Name:         jobItem.Name,
		Group:        jobItem.Group,
		Image:        res.Message,
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}
	av, err := dynamodbattribute.MarshalMap(stock)

	if err != nil {
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: &stockTableName,
		Item:      av,
	}

	_, err = db.svc.PutItem(&input)

	return err
}
