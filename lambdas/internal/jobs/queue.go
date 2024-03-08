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
	svc *sqs.Client
}

func NewQueueService() *QueueRepository {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &QueueRepository{
		svc: sqs.NewFromConfig(sdkConfig),
	}
}

func (queue QueueRepository) AddJobsToQueue(jobs []JobAction) error {
	var err error

	queueUrl := os.Getenv("SQS_QUEUE_URL")

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
		QueueUrl: &queueUrl,
		Entries:  messageRequests,
	}

	_, err = queue.svc.SendMessageBatch(context.TODO(), &input)

	return err
}
