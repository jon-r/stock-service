package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestCheckForJobs(t *testing.T) {
	t.Run("checkForJobs NoErrors", checkForJobsNoError)
	t.Run("checkForJobs NoMessages", checkForJobsNoMessages)
	t.Run("checkForJobs Errors", checkForJobsErrors)
}

func checkForJobsNoError(t *testing.T) {
	stubber, _ := test.Enter()
	mockClock := clock.NewMock()

	mockHandler := newHandler(
		handlers.NewMock(*stubber.SdkConfig),
		mockClock,
	)

	addReceiveQueueEventStub(stubber, []types.Message{
		{
			ReceiptHandle: aws.String("message1"),
			Body:          aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`),
		},
	})

	cancelSpy := func() {}
	mockHandler.checkForJobs(cancelSpy)

	queuedEvents := <-mockHandler.queueManager.queues[provider.PolygonIo]

	assert.Equal(t, queuedEvents, job.Job{
		ReceiptId: aws.String("message1"),
		JobId:     "TEST_ID",
		Provider:  provider.PolygonIo,
		Type:      job.LoadTickerDescription,
		TickerId:  "AMZN",
		Attempts:  0,
	})

	err := stubber.VerifyAllStubsCalled()
	assert.Nil(t, err)
}

func checkForJobsNoMessages(t *testing.T) {
	stubber, _ := test.Enter()
	mockClock := clock.NewMock()

	mockHandler := newHandler(
		handlers.NewMock(*stubber.SdkConfig),
		mockClock,
	)

	for range [6]int{} {
		addReceiveQueueEventStub(stubber, []types.Message{})
	}
	addDisableRuleStub(stubber)

	cancelSpyCount := 0
	cancelSpy := func() { cancelSpyCount++ }

	mockHandler.checkForJobs(cancelSpy)

	// empty once
	assert.Equal(t, 1, mockHandler.queueManager.emptyResponses)

	for range [5]int{} {
		mockHandler.checkForJobs(cancelSpy)
	}

	// empty 6 times, disable rule triggered
	assert.Equal(t, 6, mockHandler.queueManager.emptyResponses)

	err := stubber.VerifyAllStubsCalled()
	assert.Nil(t, err)
}

func checkForJobsErrors(t *testing.T) {
	stubber, _ := test.Enter()
	mockClock := clock.NewMock()

	mockHandler := newHandler(
		handlers.NewMock(*stubber.SdkConfig),
		mockClock,
	)

	for range [5]int{} {
		stubber.Add(test.StubSqsReceiveMessages(
			nil, nil, fmt.Errorf("test error"),
		))
	}

	stubber.Add(test.StubEventbridgeDisableRule(
		"EVENTBRIDGE_RULE_NAME", nil,
	))

	cancelSpyCount := 0
	cancelSpy := func() { cancelSpyCount++ }

	mockHandler.checkForJobs(cancelSpy)

	// errored once
	assert.Equal(t, 1, mockHandler.queueManager.failedAttempts)

	for range [4]int{} {
		mockHandler.checkForJobs(cancelSpy)
	}

	// errored times, cancel triggered
	assert.Equal(t, 5, mockHandler.queueManager.failedAttempts)
	assert.Equal(t, 1, cancelSpyCount)

	err := stubber.VerifyAllStubsCalled()
	assert.Nil(t, err)
}

func TestInvokeNextJob(t *testing.T) {
	t.Run("invokeNextJob NoErrors", invokeNextJobNoErrors)
	t.Run("invokeNextJob NoJobs", invokeNextJobNoJobs)
	t.Run("invokeNextJob Errors", invokeNextJobErrors)
}

func invokeNextJobNoErrors(t *testing.T) {
	stubber, _ := test.Enter()
	mockClock := clock.NewMock()

	mockHandler := newHandler(
		handlers.NewMock(*stubber.SdkConfig),
		mockClock,
	)

	mockHandler.queueManager.queues[provider.PolygonIo] <- job.Job{
		ReceiptId: aws.String("message1"),
		JobId:     "TEST_ID",
		Provider:  provider.PolygonIo,
		Type:      job.LoadHistoricalPrices,
		TickerId:  "AMZN",
		Attempts:  0,
	}

	addInvokeWorkerStub(stubber, `{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":0}`)
	addDeleteQueueStub(stubber, "message1")

	mockHandler.invokeNextJob(provider.PolygonIo)

	err := stubber.VerifyAllStubsCalled()
	assert.Nil(t, err)
}

func invokeNextJobNoJobs(t *testing.T) {
	stubber, _ := test.Enter()
	mockClock := clock.NewMock()

	mockHandler := newHandler(
		handlers.NewMock(*stubber.SdkConfig),
		mockClock,
	)

	mockHandler.invokeNextJob(provider.PolygonIo)

	err := stubber.VerifyAllStubsCalled()
	assert.Nil(t, err)
}

func invokeNextJobErrors(t *testing.T) {
	stubber, _ := test.Enter()
	mockClock := clock.NewMock()

	mockHandler := newHandler(
		handlers.NewMock(*stubber.SdkConfig),
		mockClock,
	)

	mockHandler.queueManager.queues[provider.PolygonIo] <- job.Job{
		ReceiptId: aws.String("message1"),
		JobId:     "TEST_ID",
		Provider:  provider.PolygonIo,
		Type:      job.LoadHistoricalPrices,
		TickerId:  "AMZN",
		Attempts:  0,
	}

	stubber.Add(test.StubLambdaInvoke(
		"LAMBDA_WORKER_NAME",
		nil, //[]byte(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":0}`),
		fmt.Errorf("something went wrong"),
	))

	stubber.Add(test.StubSqsSendMessage(
		"SQS_QUEUE_URL",
		`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":1}`,
		nil,
	))

	stubber.Add(test.StubSqsDeleteMessage(
		"SQS_QUEUE_URL",
		"message1",
		nil,
	))

	mockHandler.invokeNextJob(provider.PolygonIo)

	err := stubber.VerifyAllStubsCalled()
	assert.Nil(t, err)
}
