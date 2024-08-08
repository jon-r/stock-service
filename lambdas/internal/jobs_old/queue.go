package jobs_old

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type QueueRepository struct {
	svc      *sqs.Client
	QueueUrl string
	DLQUrl   string
}

func CreateSqsClient() *sqs.Client {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return sqs.NewFromConfig(sdkConfig)
}

func NewQueueService(client *sqs.Client) *QueueRepository {
	return &QueueRepository{
		svc:      client,
		QueueUrl: os.Getenv("SQS_QUEUE_URL"),
		DLQUrl:   os.Getenv("SQS_DLQ_URL"),
	}
}

func (queue QueueRepository) AddJobs(jobs []JobAction, newUuid UuidGen) error {
	var err error

	messageRequests := make([]types.SendMessageBatchRequestEntry, len(jobs))
	for i, message := range jobs {
		id := newUuid()

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

	input := sqs.SendMessageBatchInput{
		QueueUrl: aws.String(queue.QueueUrl),
		Entries:  messageRequests,
	}

	_, err = queue.svc.SendMessageBatch(context.TODO(), &input)

	return err
}

func (queue QueueRepository) RetryJob(job JobAction, failReason string, newUuid UuidGen) error {
	var err error
	updatedJob := job
	updatedJob.Attempts += 1

	if updatedJob.Attempts > 3 {
		err = queue.AddJobToDLQ(updatedJob, failReason)
	} else {
		// put the failed item back into the queue
		err = queue.AddJobs([]JobAction{updatedJob}, newUuid)
	}

	return err
}

func (queue QueueRepository) ReceiveJobs() (*[]JobQueueItem, error) {
	input := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queue.QueueUrl),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
	}

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

func (queue QueueRepository) AddJobToDLQ(job JobAction, failReason string) error {
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
