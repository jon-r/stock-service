package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

var jobTableName = "stock-app_Job"

type JobInput struct {
	Name  string
	Group string
}

type JobItem struct {
	JobId string
	Name  string
	Group string
}

func (db DatabaseRepository) InsertJobs(jobInputs []JobInput) error {
	var err error

	writeRequests := make([]*dynamodb.WriteRequest, len(jobInputs))

	for i, jobInput := range jobInputs {
		job := JobItem{
			JobId: uuid.NewString(),
			Name:  jobInput.Name,
			Group: jobInput.Group,
		}

		av, err := dynamodbattribute.MarshalMap(job)
		if err != nil {
			break
		}
		writeRequests[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: av},
		}
	}

	if err != nil {
		return err
	}

	input := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			jobTableName: writeRequests,
		},
	}
	_, err = db.svc.BatchWriteItem(&input)

	return err
}

func (db DatabaseRepository) FindJobByGroup(group string) (*JobItem, error) {
	keyEx := expression.Key("Group").Equal(expression.Value(group))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		return nil, err
	}

	input := dynamodb.QueryInput{
		TableName:                 aws.String(jobTableName),
		IndexName:                 aws.String("GroupIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := db.svc.Query(&input)

	if err != nil {
		return nil, err
	}

	item := result.Items[0]

	if item == nil {
		return nil, nil
	}

	job := new(JobItem)

	err = dynamodbattribute.UnmarshalMap(item, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (db DatabaseRepository) DeleteJob(job *JobItem) error {
	var err error

	request := dynamodb.DeleteItemInput{
		TableName: &jobTableName,
		Key: map[string]*dynamodb.AttributeValue{
			"JobId": {S: aws.String(job.JobId)},
		},
	}

	_, err = db.svc.DeleteItem(&request)

	return err
}
