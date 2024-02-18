package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type JobItem struct {
	JobId string `json:"JobId"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type StockItem struct {
	StockIndexId string `json:"StockIndexId"`
	Name         string `json:"name"`
	Group        string `json:"group"`
	Image        string `json:"image"`
	UpdatedAt    string `json:"updatedAt"`
}

type DogApiRes struct {
	status  string
	message string
}

type QueueEvent struct {
	Group string
	Delay int
}

var awsSession = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var dbService = dynamodb.New(awsSession)
var sqsService = sqs.New(awsSession)

func getJobItemByGroup(group string) (*JobItem, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("stock-app_Job"),
		Key: map[string]*dynamodb.AttributeValue{
			"Group": {
				S: aws.String(group),
			},
		},
	}

	result, err := dbService.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	job := new(JobItem)
	err = dynamodbattribute.UnmarshalMap(result.Item, job)

	if err != nil {
		return nil, err
	}

	return job, nil
}

func setDbItems(res *DogApiRes, jobItem *JobItem) (*dynamodb.BatchWriteItemInput, error) {
	stock := StockItem{
		StockIndexId: uuid.NewString(),
		Name:         jobItem.Name,
		Group:        jobItem.Group,
		Image:        res.message,
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}
	av, err := dynamodbattribute.MarshalMap(stock)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"stock-app_Job": {&dynamodb.WriteRequest{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: map[string]*dynamodb.AttributeValue{
						"JobId": {
							S: aws.String(jobItem.JobId),
						},
					},
				},
			},
			},
			"stock-app_StockIndex": {
				&dynamodb.WriteRequest{
					PutRequest: &dynamodb.PutRequest{
						Item: av,
					},
				},
			},
			// todo message would come after success or on errors. maybe not needed if cloudwatch covers it?
		},
	}

	return input, nil
}

func fetchItem() (*DogApiRes, error) {
	res, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data DogApiRes
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func jobsToSQS(queueItem QueueEvent) (*sqs.SendMessageInput, error) {
	delay := int64(queueItem.Delay)

	data, err := json.Marshal(queueItem)
	if err != nil {
		return nil, err
	}

	body := string(data)

	input := &sqs.SendMessageInput{
		DelaySeconds: &delay,
		MessageBody:  &body,
	}

	return input, nil
}

func handleRequest(ctx context.Context, event events.SQSEvent) {
	var queueItem QueueEvent
	err := json.Unmarshal([]byte(event.Records[0].Body), &queueItem)
	if err != nil {
		log.Fatalf("Error parsing queue event: %s", err)
	} else {
		log.Printf("Handling event: %s", queueItem.Group)
	}

	job, err := getJobItemByGroup(queueItem.Group)
	if err != nil {
		log.Fatalf("Error getting item: %s", err)
	} else {
		log.Printf("Job: %s", job.Name)
	}

	if job == nil {
		// no more items to fetch
		return
	}

	res, err := fetchItem()
	if err != nil {
		log.Fatalf("Error calling http.get: %s", err)
	}

	dynamoDbInput, err := setDbItems(res, job)
	if err != nil {
		log.Fatalf("Error preparing dynamoDb dataa: %s", err)
	}

	_, err = dbService.BatchWriteItem(dynamoDbInput)
	if err != nil {
		log.Fatalf("Error calling dynamodb.WriteItem: %s", err)
	} else {
		log.Println("Successfully added items to tables")
	}

	sqsInput, err := jobsToSQS(queueItem)
	if err != nil {
		log.Fatalf("Error preparing sqs event")
	}

	_, err = sqsService.SendMessage(sqsInput)
	if err != nil {
		log.Fatalf("Error calling sqs.SendMessage: %s", err)
	} else {
		log.Println("Successfully added items to queue")
	}
}

func main() {
	lambda.Start(handleRequest)
}
