package queue

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type QueueMessage struct {
	Group string
	Delay int
}

func (queue QueueRepository) SendDelayedEvents(queueItems []QueueMessage) error {
	var err error
	var queueUrl = os.Getenv("SQS_QUEUE_URL")

	messageRequests := make([]types.SendMessageBatchRequestEntry, len(queueItems))
	for i, item := range queueItems {
		id := uuid.NewString()
		delay := int32(item.Delay)

		data, err := json.Marshal(item)
		if err != nil {
			break
		}

		body := string(data)

		event := types.SendMessageBatchRequestEntry{
			Id:           &id,
			DelaySeconds: delay,
			MessageBody:  &body,
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

func ParseQueueEvent(event events.SQSEvent) (*QueueMessage, error) {
	var message QueueMessage

	err := json.Unmarshal([]byte(event.Records[0].Body), &message)

	if err != nil {
		return nil, err
	}

	return &message, nil
}
