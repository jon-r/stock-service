package jobs

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type QueueRepository struct {
	svc      *sqs.Client
	QueueUrl string
	DLQUrl   string
}

func NewQueueService() *QueueRepository {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &QueueRepository{
		svc:      sqs.NewFromConfig(sdkConfig),
		QueueUrl: os.Getenv("SQS_QUEUE_URL"),
		DLQUrl:   os.Getenv("SQS_DLQ_URL"),
	}
}

func (queue QueueRepository) AddJobs(jobs []JobAction) error {
	var err error

	messageRequests := make([]types.SendMessageBatchRequestEntry, len(jobs))
	for i, message := range jobs {
		id := uuid.NewString()

		data, err := json.Marshal(message)
		if err != nil {
			break
		}

		body := string(data)

		event := types.SendMessageBatchRequestEntry{
			Id:          &id,
			MessageBody: &body,
		}

		messageRequests[i] = event
	}

	if err != nil {
		return err
	}

	log.Printf("Queue url: %v\n", queue.QueueUrl)
	log.Printf("attempt to add items %v\n", messageRequests)

	input := sqs.SendMessageBatchInput{
		QueueUrl: aws.String(queue.QueueUrl),
		Entries:  messageRequests,
	}

	res, err := queue.svc.SendMessageBatch(context.TODO(), &input)

	log.Printf("queue res: %+v\n", res)

	return err
}

func (queue QueueRepository) ReceiveJobs() (*[]JobQueueItem, error) {
	input := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queue.QueueUrl),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
	}

	log.Println("attempt to receive items...")

	result, err := queue.svc.ReceiveMessage(context.TODO(), &input)

	if err != nil {
		return nil, err
	}

	jobs := make([]JobQueueItem, len(result.Messages))

	for i, message := range result.Messages {
		queueItem := JobQueueItem{
			RecieptHandle: *message.ReceiptHandle,
		}
		err = json.Unmarshal([]byte(*message.Body), &queueItem.Action)
		jobs[i] = queueItem
	}

	return &jobs, nil
}

func (queue QueueRepository) DeleteJob(receiptHandle string) error {
	input := sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue.QueueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := queue.svc.DeleteMessage(context.TODO(), &input)

	return err
}

func (queue QueueRepository) AddJobToDLQ(job JobAction, failReason error) error {
	var err error

	data, err := json.Marshal(JobErrorItem{
		JobAction:   job,
		ErrorReason: failReason,
	})
	if err != nil {
		return err
	}

	body := string(data)

	input := sqs.SendMessageInput{
		QueueUrl:    aws.String(queue.DLQUrl),
		MessageBody: &body,
	}

	_, err = queue.svc.SendMessage(context.TODO(), &input)

	return err
}
