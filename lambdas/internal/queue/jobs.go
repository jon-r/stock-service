package queue

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type QueueMessage struct {
	Group string
	Delay int
}

func (queue QueueRepository) SendDelayedEvents(queueItems []QueueMessage) error {
	var err error
	var queueUrl = os.Getenv("SQS_QUEUE_URL")

	messageRequests := make([]*sqs.SendMessageBatchRequestEntry, len(queueItems))
	for i, item := range queueItems {
		id := uuid.NewString()
		delay := int64(item.Delay)

		data, err := json.Marshal(item)
		if err != nil {
			break
		}

		body := string(data)

		event := sqs.SendMessageBatchRequestEntry{
			Id:           &id,
			DelaySeconds: &delay,
			MessageBody:  &body,
		}

		messageRequests[i] = &event
	}

	if err != nil {
		return err
	}

	input := sqs.SendMessageBatchInput{
		QueueUrl: &queueUrl,
		Entries:  messageRequests,
	}

	_, err = queue.svc.SendMessageBatch(&input)

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
