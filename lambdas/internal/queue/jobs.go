package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type Message struct {
	QueueGroup string
}

//type QueueMessage struct {
//	Group string
//	Delay int
//}

type GroupSettings struct {
	Delay int32
	Url   string
}

type ParsedMessage struct {
	Message
	GroupSettings
}

var GroupSettingsList = map[string]GroupSettings{
	"slow": {Delay: 7, Url: "https://dog.ceo/api/breeds/image/random"},
	"fast": {Delay: 4, Url: "https://dog.ceo/api/breeds/image/random"},
}

func (queue QueueRepository) SendDelayedEvents(queueMessages []Message) error {
	var err error
	var queueUrl = os.Getenv("SQS_QUEUE_URL")

	messageRequests := make([]types.SendMessageBatchRequestEntry, len(queueMessages))
	for i, message := range queueMessages {
		id := uuid.NewString()

		settings, ok := GroupSettingsList[message.QueueGroup]
		if !ok {
			err = fmt.Errorf("no settings available for group %v", message.QueueGroup)
		}

		data, err := json.Marshal(message)
		if err != nil {
			break
		}

		body := string(data)

		fmt.Printf("settings: %v", settings)

		// fixme delay is not happening
		event := types.SendMessageBatchRequestEntry{
			Id:           &id,
			DelaySeconds: settings.Delay,
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

func ParseQueueEvent(event events.SQSEvent) (*ParsedMessage, error) {
	var message Message

	err := json.Unmarshal([]byte(event.Records[0].Body), &message)

	if err != nil {
		return nil, err
	}

	settings, ok := GroupSettingsList[message.QueueGroup]

	if !ok {
		err = fmt.Errorf("no settings available for group %v", message.QueueGroup)
		return nil, err
	}

	var parsedMessage = ParsedMessage{
		message,
		settings,
	}

	return &parsedMessage, nil
}
