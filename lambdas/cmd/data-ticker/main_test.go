package main

import (
	"context"
	"testing"
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/clock"
	"github.com/jon-r/stock-service/lambdas/internal/testutil"
)

func TestPollSqsQueue(t *testing.T) {
	t.Run("NoErrors", pollSqsQueueNoErrors)
}

//type mockTimer struct{}
//
//func (mockTimer) Sleep(d time.Duration) { /* do nothing */ }
//func (mockTimer) NewTicker(d time.Duration) *time.Ticker {
//	ticker := time.NewTicker(time.Second)
//
//	return ticker
//}

// fixme not sure how to mock this

func pollSqsQueueNoErrors(t *testing.T) {
	stubber, mockServiceHandler := testutil.EnterTest(nil)
	mockClock := clock.MockClock()

	mockHandler := DataTickerHandler{
		ServiceHandler: *mockServiceHandler,
		Clock:          mockClock,
	}

	//expectedQueueInput := &sqs.ReceiveMessageInput{
	//	QueueUrl:            aws.String(""),
	//	MaxNumberOfMessages: 10,
	//	WaitTimeSeconds:     5,
	//}
	//queueResponse := &sqs.ReceiveMessageOutput{
	//	Messages: []types.Message{
	//		{
	//			ReceiptHandle: aws.String("message1"),
	//			Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`),
	//		},
	//		{
	//			ReceiptHandle: aws.String("message2"),
	//			Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":0}`),
	//		},
	//	},
	//}
	//stubber.Add(testutil.StubSqsReceiveMessages(expectedQueueInput, queueResponse, nil))
	//// todo loop add this multiple times (to time the ticker out)
	//stubber.Add(testutil.StubSqsReceiveMessages(
	//	expectedQueueInput,
	//	&sqs.ReceiveMessageOutput{Messages: []types.Message{}},
	//	nil,
	//))

	// todo:
	//  - mock invoke
	//  - mock delete message
	//
	testDone := make(chan struct{})
	go func() {
		// todo grab errors
		mockHandler.pollSqsQueue(context.TODO())
		//if err != nil {
		//
		//}
		time.AfterFunc(time.Second, func() {
			close(testDone)
		})
	}()

	mockClock.AdvanceTime(10 * time.Minute)

	<-testDone

	testutil.Assert(stubber, nil, nil, t)
}
