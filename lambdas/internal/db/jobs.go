package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"jon-richards.com/stock-app/internal/providers"
)

var jobTableName = "stock-app_Job"

type ActionTypes string

const (
	NewStockItem       ActionTypes = "NEW_STOCK_ITEM"
	PopulateItemPrices ActionTypes = "POPULATE_ITEM_PRICES"
	UpdateAllPrices    ActionTypes = "UPDATE_ALL_PRICES"
	// todo
	// UpdateDividends
	// ???
)

type JobInput struct {
	Provider providers.ProviderName
	Type     ActionTypes
	TickerId string
}

type JobInputPayload = map[string]string

type JobItem struct {
	JobId string
	JobInput
}

func (db DatabaseRepository) InsertJob(jobInput JobInput) error {
	var err error

	av, err := attributevalue.MarshalMap(jobInput)

	if err != nil {
		return err
	}
	input := dynamodb.PutItemInput{
		Item:      av,
		TableName: &jobTableName,
	}

	_, err = db.svc.PutItem(context.TODO(), &input)

	return err
}

func (db DatabaseRepository) InsertJobs(jobInputs []JobInput) error {
	var err error

	writeRequests := make([]types.WriteRequest, len(jobInputs))

	for i, jobInput := range jobInputs {
		job := JobItem{
			JobId:    uuid.NewString(),
			JobInput: jobInput,
		}

		av, err := attributevalue.MarshalMap(job)
		if err != nil {
			break
		}
		writeRequests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{Item: av},
		}
	}

	if err != nil {
		return err
	}

	input := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			jobTableName: writeRequests,
		},
	}
	_, err = db.svc.BatchWriteItem(context.TODO(), &input)

	return err
}

func (db DatabaseRepository) FindJobByProvider(provider providers.ProviderName) (*JobItem, error) {
	keyEx := expression.Key("Provider").Equal(expression.Value(provider))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		return nil, err
	}

	input := dynamodb.QueryInput{
		TableName:                 aws.String(jobTableName),
		IndexName:                 aws.String("ProviderIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int32(1),
	}

	result, err := db.svc.Query(context.TODO(), &input)

	if err != nil {
		return nil, err
	}

	item := result.Items[0]

	if item == nil {
		return nil, nil
	}

	job := new(JobItem)

	err = attributevalue.UnmarshalMap(item, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (db DatabaseRepository) DeleteJob(job *JobItem) error {
	var err error

	id, err := attributevalue.Marshal(job.JobId)

	if err != nil {
		return err
	}

	request := dynamodb.DeleteItemInput{
		TableName: &jobTableName,
		Key: map[string]types.AttributeValue{
			"JobId": id,
		},
	}

	_, err = db.svc.DeleteItem(context.TODO(), &request)

	return err
}
