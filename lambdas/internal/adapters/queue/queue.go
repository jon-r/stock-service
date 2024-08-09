package queue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jon-r/stock-service/lambdas/internal/utils/array"
)

type queue struct {
	client *sqs.Client
	idGen  NewIdFunc
}

type Broker interface {
	SendMessage(queueUrl string, message interface{}) (*sqs.SendMessageOutput, error)
	SendMessages(queueUrl string, messages interface{}) (*sqs.SendMessageBatchOutput, error)
	ReceiveMessages(queueUrl string) (*[]types.Message, error)
	DeleteMessage(queueUrl string, messageId string) (*sqs.DeleteMessageOutput, error)
}

func (q *queue) SendMessage(queueUrl string, message interface{}) (*sqs.SendMessageOutput, error) {
	var err error

	data, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	messageRequest := sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(string(data)),
	}

	return q.client.SendMessage(context.TODO(), &messageRequest)
}

func (q *queue) SendMessages(queueUrl string, messages interface{}) (*sqs.SendMessageBatchOutput, error) {
	var err error
	var data []byte

	slice := array.UnpackArray(messages)

	messageRequests := make([]types.SendMessageBatchRequestEntry, len(slice))
	for i, message := range slice {
		id := q.idGen()

		data, err = json.Marshal(message)

		if err != nil {
			return nil, err
		}

		messageRequests[i] = types.SendMessageBatchRequestEntry{
			Id:          aws.String(id),
			MessageBody: aws.String(string(data)),
		}
	}

	input := sqs.SendMessageBatchInput{
		QueueUrl: aws.String(queueUrl),
		Entries:  messageRequests,
	}

	return q.client.SendMessageBatch(context.TODO(), &input)
}

func (q *queue) ReceiveMessages(queueUrl string) (*[]types.Message, error) {
	input := sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
	}

	res, err := q.client.ReceiveMessage(context.TODO(), &input)

	if err != nil {
		return nil, err
	}

	return &res.Messages, nil
}

func (q *queue) DeleteMessage(queueUrl string, messageId string) (*sqs.DeleteMessageOutput, error) {
	input := sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(messageId),
	}

	return q.client.DeleteMessage(context.TODO(), &input)
}

func NewBroker(config aws.Config, idGen NewIdFunc) Broker {
	return &queue{
		client: sqs.NewFromConfig(config),
		idGen:  idGen,
	}
}
