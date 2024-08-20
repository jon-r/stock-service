package queue

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	stubber := testtools.NewStubber()
	idGen := func() string { return "TEST_ID" }
	client := NewBroker(*stubber.SdkConfig, idGen)

	t.Run("SendMessage", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "SendMessage",
			Input: &sqs.SendMessageInput{
				MessageBody: aws.String(`{"JobId":"1234","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"APPL","Attempts":0}`),
				QueueUrl:    aws.String("queue1"),
			},
			Output: &sqs.SendMessageOutput{},
		})

		message := job.Job{
			JobId:    "1234",
			Provider: provider.PolygonIo,
			Type:     job.LoadTickerDescription,
			TickerId: "APPL",
			Attempts: 0,
		}

		_, err := client.SendMessage("queue1", message)

		fmt.Printf("%+v", err)

		assert.Nil(t, err)
	})
	t.Run("SendMessages", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "SendMessageBatch",
			Input: &sqs.SendMessageBatchInput{
				QueueUrl: aws.String("queue1"),
				Entries: []types.SendMessageBatchRequestEntry{
					{
						Id:          aws.String("TEST_ID"),
						MessageBody: aws.String(`{"JobId":"1234","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"APPL","Attempts":0}`),
					},
					{
						Id:          aws.String("TEST_ID"),
						MessageBody: aws.String(`{"JobId":"5678","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"APPL","Attempts":0}`),
					},
				},
			},
			Output: &sqs.SendMessageBatchOutput{},
		})

		messages := []job.Job{
			{
				JobId:    "1234",
				Provider: provider.PolygonIo,
				Type:     job.LoadTickerDescription,
				TickerId: "APPL",
				Attempts: 0,
			},
			{
				JobId:    "5678",
				Provider: provider.PolygonIo,
				Type:     job.LoadHistoricalPrices,
				TickerId: "APPL",
				Attempts: 0,
			},
		}

		_, err := client.SendMessages("queue1", messages)

		fmt.Printf("%+v", err)

		assert.Nil(t, err)
	})
	t.Run("ReceiveMessages", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("queue1"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						ReceiptHandle: aws.String("abcd"),
						Body:          aws.String(`{"JobId":"1234","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"APPL","Attempts":0}`),
					},
					{
						ReceiptHandle: aws.String("abcd"),
						Body:          aws.String(`{"JobId":"5678","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"APPL","Attempts":0}`),
					},
				},
			},
		})

		res, err := client.ReceiveMessages("queue1")

		fmt.Printf("%+v", res)

		assert.Nil(t, err)
		assert.Equal(t, &[]types.Message{
			{
				ReceiptHandle: aws.String("abcd"),
				Body:          aws.String(`{"JobId":"1234","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"APPL","Attempts":0}`),
			},
			{
				ReceiptHandle: aws.String("abcd"),
				Body:          aws.String(`{"JobId":"5678","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"APPL","Attempts":0}`),
			},
		}, res)
	})

	t.Run("DeleteMessage", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("queue1"),
				ReceiptHandle: aws.String("abcd"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		_, err := client.DeleteMessage("queue1", aws.String("abcd"))

		assert.Nil(t, err)
	})
}
