package main

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
)

func TestPollSqsQueue(t *testing.T) {
	t.Run("NoErrors", handleRequestNoErrors)
}

func handleRequestNoErrors(t *testing.T) {
	stubber, ctx := test.Enter()
	mockClock := clock.NewMock()

	mockHandler := newHandler(
		handlers.NewMock(*stubber.SdkConfig),
		mockClock,
	)

	expectedQueueInput := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String("SQS_QUEUE_URL"),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
	}
	queueResponse := &sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{
				ReceiptHandle: aws.String("message1"),
				Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`),
			},
		},
	}
	stubber.Add(test.StubSqsReceiveMessages(expectedQueueInput, queueResponse, nil))

	payloadJson := `{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`
	stubber.Add(test.StubLambdaInvoke("LAMBDA_WORKER_NAME", []byte(payloadJson), nil))
	stubber.Add(test.StubSqsDeleteMessage("SQS_QUEUE_URL", "message1", nil))

	// todo grab errors
	go mockHandler.HandleRequest(ctx)

	// have to add one second at a time otherwise the mock tickers all stack up
	for range [15]int{} {
		mockClock.Add(1 * time.Second)
	}

	test.AssertAll(t, stubber, nil, nil)
}
