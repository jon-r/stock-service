package main

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/testutil"
)

func TestPollSqsQueue(t *testing.T) {
	t.Run("NoErrors", pollSqsQueueNoErrors)
}

func pollSqsQueueNoErrors(t *testing.T) {
	stubber, mockServiceHandler := testutil.EnterTest(nil)
	mockClock := clock.NewMock()

	mockHandler := DataTickerHandler{
		ServiceHandler: *mockServiceHandler,
		Clock:          mockClock,
		done:           make(chan bool),
	}

	expectedQueueInput := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(""),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
	}
	queueResponse := &sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{
				ReceiptHandle: aws.String("message1"),
				Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`),
			},
			{
				ReceiptHandle: aws.String("message2"),
				Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":0}`),
			},
		},
	}
	stubber.Add(testutil.StubSqsReceiveMessages(expectedQueueInput, queueResponse, nil))

	// time based on how many times this gets triggered before the main function stops running
	for i := 1; i <= 6; i++ {
		stubber.Add(testutil.StubSqsReceiveMessages(
			expectedQueueInput,
			&sqs.ReceiveMessageOutput{Messages: []types.Message{}},
			nil,
		))
	}

	stubber.Add(testutil.StubEventbridgeDisableRule("EVENTBRIDGE_RULE_NAME", nil))

	for i := 1; i <= 7; i++ {
		stubber.Add(testutil.StubSqsReceiveMessages(
			expectedQueueInput,
			&sqs.ReceiveMessageOutput{Messages: []types.Message{}},
			nil,
		))
	}

	// todo:
	//  - mock invoke
	//  - mock delete message
	//
	testDone := make(chan bool)
	go func() {
		// todo grab errors
		mockHandler.handleQueuedJobs(context.TODO())

		testDone <- true
	}()

	mockClock.Add(10 * time.Minute)

	<-testDone

	testutil.Assert(stubber, nil, nil, t)
}
