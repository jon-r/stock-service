package queue

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type QueueMessage struct {
	Group string
	Delay int
}

func (queue QueueRepository) InsertDelayedEvents(queueItems []QueueMessage) error {
	var err error
	var queueUrl = os.Getenv("SQS_QUEUE_URL")

	messageRequests := make([]*sqs.SendMessageBatchRequestEntry, len(queueItems))
	for i, item := range queueItems {
		id := uuid.NewString()
		delay := int64(item.Delay)

		data, err := json.Marshal(item)
		if err != nil {
			log.Printf("could not convert queue item to queue message")
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
