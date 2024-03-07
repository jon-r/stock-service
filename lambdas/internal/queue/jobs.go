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

	"jon-richards.com/stock-app/internal/providers"
)

type Message struct {
	Provider providers.ProviderName
}

//type QueueMessage struct {
//	Provider string
//	Delay int
//}

//type Settings struct {
//	Delay int32
//	Url   string
//}

type ParsedMessage struct {
	Message
	providers.Settings
}

//var GroupSettingsList = map[string]Settings{
//	"slow": {Delay: 7, Url: "https://dog.ceo/api/breeds/image/random"},
//	"fast": {Delay: 4, Url: "https://dog.ceo/api/breeds/image/random"},
//}

func (queue QueueRepository) SendDelayedEvents(queueMessages []Message) error {
	var err error
	var queueUrl = os.Getenv("SQS_QUEUE_URL")

	messageRequests := make([]types.SendMessageBatchRequestEntry, len(queueMessages))
	for i, message := range queueMessages {
		id := uuid.NewString()

		settings := providers.GetSettings(message.Provider)

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

	settings := providers.GetSettings(message.Provider)

	var parsedMessage = ParsedMessage{
		message,
		settings,
	}

	return &parsedMessage, nil
}
