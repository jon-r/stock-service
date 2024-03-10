package jobs

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type QueueRepository struct {
	svc      *sqs.Client
	queueUrl string
}

func NewQueueService() *QueueRepository {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &QueueRepository{
		svc:      sqs.NewFromConfig(sdkConfig),
		queueUrl: os.Getenv("SQS_QUEUE_URL"),
	}
}

func (queue QueueRepository) AddJobsToQueue(jobs []JobAction) error {
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

	input := sqs.SendMessageBatchInput{
		QueueUrl: &queue.queueUrl,
		Entries:  messageRequests,
	}

	_, err = queue.svc.SendMessageBatch(context.TODO(), &input)

	return err
}

func (queue QueueRepository) GetJobsFromQueue() (*[]JobAction, error) {
	input := sqs.ReceiveMessageInput{
		QueueUrl:            &queue.queueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     10,
	}

	result, err := queue.svc.ReceiveMessage(context.TODO(), &input)

	if err != nil {
		return nil, err
	}

	jobs := make([]JobAction, len(result.Messages))

	for i, message := range result.Messages {
		err = json.Unmarshal([]byte(*message.Body), &jobs[i])
	}

	return &jobs, nil
}
