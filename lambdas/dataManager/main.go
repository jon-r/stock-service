package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type JobItem struct {
	JobId string `json:"JobId"`
	Name  string `json:"Name"`
	Group string `json:"Group"`
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

func namesToDbJobs(names []string, delay int, group string) ([]*dynamodb.WriteRequest, QueueEvent) {
	requests := make([]*dynamodb.WriteRequest, len(names))
	for i, name := range names {
		job := JobItem{
			JobId: uuid.NewString(),
			Name:  name,
			Group: group,
		}

		av, err := dynamodbattribute.MarshalMap(job)
		if err != nil {
			log.Fatalf("Error marshalling new job item: %s", err)
		}
		requests[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		}
	}

	queueItem := QueueEvent{
		Group: group,
		Delay: delay,
	}

	return requests, queueItem
}

// todo split this file up. move sqs + dynamo db helpers to own files
func jobsToSQS(queueItems []QueueEvent) ([]*sqs.SendMessageBatchRequestEntry, error) {
	requests := make([]*sqs.SendMessageBatchRequestEntry, len(queueItems))
	for i, item := range queueItems {
		id := uuid.NewString()
		delay := int64(item.Delay)

		data, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}

		body := string(data)

		event := &sqs.SendMessageBatchRequestEntry{
			Id:           &id,
			DelaySeconds: &delay,
			MessageBody:  &body,
		}

		requests[i] = event
	}

	return requests, nil
}

func fakeInputs() ([]*dynamodb.WriteRequest, []QueueEvent) {
	names1 := []string{
		"Phoebe",
		"Harley",
		"Bandit",
		"Delilah",
		"Tiger",
		"Panda",
	}
	names2 := []string{
		"Whiskey",
		"Jasper",
		"Belle",
		"Shelby",
		"Zara",
		"Bruno",
	}

	// todo these values are the rate limit for api requests
	jobs1, queue1 := namesToDbJobs(names1, 5, "slow")
	jobs2, queue2 := namesToDbJobs(names2, 3, "fast")

	return append(jobs1, jobs2...), []QueueEvent{queue1, queue2}
}

func loadData() {
	jobs, queue := fakeInputs()

	tableName := "stock-app_Job"

	dynamoDbInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			tableName: jobs,
		},
	}

	_, err := dbService.BatchWriteItem(dynamoDbInput)
	if err != nil {
		log.Fatalf("Error calling dynamodb.WriteItem: %s", err)
	} else {
		log.Println("Successfully added items to table " + tableName)
	}

	queueUrl := os.Getenv("SQS_QUEUE_URL")
	entries, err := jobsToSQS(queue)
	if err != nil {
		log.Fatalf("Error preparing sqs data: %s", err)
	}

	sqsInput := &sqs.SendMessageBatchInput{
		QueueUrl: &queueUrl,
		Entries:  entries,
	}

	_, err = sqsService.SendMessageBatch(sqsInput)
	if err != nil {
		log.Fatalf("Error calling sqs.SendMessage: %s", err)
	} else {
		log.Println("Successfully added items to queue")
	}
}

func main() {
	lambda.Start(loadData)
}
