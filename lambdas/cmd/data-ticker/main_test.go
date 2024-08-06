package main

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
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

	receiveQueueEvent(stubber, []types.Message{
		{
			ReceiptHandle: aws.String("message1"),
			Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`),
		},
		{
			ReceiptHandle: aws.String("message2"),
			Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":0}`),
		},
	})

	receiveQueueEvent(stubber, []types.Message{})
	invokeWorkerEvent(stubber, `{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`)
	deleteQueueEvent(stubber, "message1")
	receiveQueueEvent(stubber, []types.Message{})
	invokeWorkerEvent(stubber, `{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":0}`)
	deleteQueueEvent(stubber, "message2")
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	disableRuleEvent(stubber)
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})
	receiveQueueEvent(stubber, []types.Message{})

	//receiveQueueEvent(stubber, []types.Message{})

	testDone := make(chan bool)
	go func() {
		// todo grab errors
		mockHandler.handleQueuedJobs(context.TODO())

		testDone <- true
	}()

	// todo this log needs to be here or the tests breaks. not sure why??
	mockHandler.LogService.Debugln("fast forward 10min")
	mockClock.Add(10 * time.Minute)

	stubber.Clear() // clear any lingering poll events

	<-testDone

	testutil.Assert(stubber, nil, nil, t)
}

func receiveQueueEvent(stubber *testtools.AwsmStubber, messages []types.Message) {
	expectedQueueInput := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String("SQS_QUEUE_URL"),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
	}
	queueResponse := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}
	stubber.Add(testutil.StubSqsReceiveMessages(expectedQueueInput, queueResponse, nil))
}
func deleteQueueEvent(stubber *testtools.AwsmStubber, messageId string) {
	stubber.Add(testutil.StubSqsDeleteMessage("SQS_QUEUE_URL", messageId, nil))
}

func invokeWorkerEvent(stubber *testtools.AwsmStubber, payloadJson string) {
	stubber.Add(testutil.StubLambdaInvoke("LAMBDA_WORKER_NAME", []byte(payloadJson), nil))
}
func disableRuleEvent(stubber *testtools.AwsmStubber) {
	stubber.Add(testutil.StubEventbridgeDisableRule("EVENTBRIDGE_RULE_NAME", nil))
}
