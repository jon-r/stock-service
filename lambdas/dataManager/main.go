package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"log"
	"os"
)

type JobEvent struct {
	JobId string `json:"JobId"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type QueueEvent struct {
	Group string
	Delay int
}

func namesToDbJobs(names []string, delay int, group string) ([]*dynamodb.WriteRequest, QueueEvent) {
	requests := make([]*dynamodb.WriteRequest, len(names))
	for i, name := range names {
		job := JobEvent{
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

// todo split this file up. move sqs + dynmo db helpers to own files
func jobsToSQS(queueItem []QueueEvent) []*sqs.SendMessageBatchRequestEntry {
	requests := make([]*sqs.SendMessageBatchRequestEntry, len(queueItem))
	for i, item := range queueItem {
		id := uuid.NewString()
		delay := int64(item.Delay)
		event := &sqs.SendMessageBatchRequestEntry{
			Id:           &id,
			DelaySeconds: &delay,
			MessageBody:  &item.Group,
		}

		requests[i] = event
	}

	return requests
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

	jobs1, queue1 := namesToDbJobs(names1, 60, "slow")
	jobs2, queue2 := namesToDbJobs(names2, 25, "fast")

	return append(jobs1, jobs2...), []QueueEvent{queue1, queue2}
}

func loadData(ctx context.Context) {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dbService := dynamodb.New(awsSession)
	sqsService := sqs.New(awsSession)

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

	// todo use event bridge instead of queue
	queueUrl := os.Getenv("SQS_QUEUE_URL")
	sqsInput := &sqs.SendMessageBatchInput{
		QueueUrl: &queueUrl,
		Entries:  jobsToSQS(queue),
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
