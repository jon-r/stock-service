package db

import (
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
