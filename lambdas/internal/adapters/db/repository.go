package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type database struct {
	client *dynamodb.Client
}

type Repository interface {
	HealthCheck() bool
	Create(tableName string, entity interface{}) (*dynamodb.PutItemOutput, error)
	// todo can type key?
	Update(tableName string, key map[string]types.AttributeValue, update expression.Expression) (*dynamodb.UpdateItemOutput, error)
	AddMany(tableName string, entities interface{}) (int, error)
	GetMany(tableName string, query expression.Expression) ([]map[string]types.AttributeValue, error)
}

// https://github.com/minhtran241/dynamodb-go-crud/blob/main/internal/repository/adapter/adapter.go

func (db *database) HealthCheck() bool {
	_, err := db.client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	return err == nil
}

func (db *database) Create(tableName string, entity interface{}) (*dynamodb.PutItemOutput, error) {
	av, err := attributevalue.MarshalMap(entity)

	if err != nil {
		return nil, err
	}

	input := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	return db.client.PutItem(context.TODO(), &input)
}

func (db *database) Update(tableName string, key map[string]types.AttributeValue, update expression.Expression) (*dynamodb.UpdateItemOutput, error) {
	input := dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		ExpressionAttributeNames:  update.Names(),
		ExpressionAttributeValues: update.Values(),
		UpdateExpression:          update.Update(),
	}

	return db.client.UpdateItem(context.TODO(), &input)
}

func (db *database) AddMany(tableName string, entities interface{}) (int, error) {
	var err error
	var data map[string]types.AttributeValue
	slice := unpackArray(entities)

	items := make([]map[string]types.AttributeValue, len(slice))
	for i, entity := range slice {
		data, err = attributevalue.MarshalMap(entity)
		items[i] = data
	}

	if err != nil {
		return 0, err
	}

	written := 0
	batchSize := 25
	start := 0
	end := start + batchSize

	for start < len(items) {
		var writeReqs []types.WriteRequest
		if end > len(items) {
			end = len(items)
		}
		for _, item := range items[start:end] {
			writeReqs = append(writeReqs, types.WriteRequest{
				PutRequest: &types.PutRequest{Item: item},
			})
		}
		_, err = db.client.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{tableName: writeReqs},
		})
		if err != nil {
			return 0, err
		} else {
			written += len(writeReqs)
		}
		start = end
		end += batchSize
	}

	return written, nil
}

func (db *database) GetMany(tableName string, query expression.Expression) ([]map[string]types.AttributeValue, error) {
	var items []map[string]types.AttributeValue
	var err error
	var response *dynamodb.ScanOutput

	scanPaginator := dynamodb.NewScanPaginator(db.client, &dynamodb.ScanInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  query.Names(),
		ExpressionAttributeValues: query.Values(),
		FilterExpression:          query.Filter(),
		ProjectionExpression:      query.Projection(),
	})

	for scanPaginator.HasMorePages() {
		response, err = scanPaginator.NextPage(context.TODO())
		if err != nil {
			break
		} else {
			items = append(items, response.Items...)
		}
	}

	return items, err
}

func NewRepository(config aws.Config) Repository {
	return &database{
		client: dynamodb.NewFromConfig(config),
	}
}
