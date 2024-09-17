package test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/stretchr/testify/assert"
)

func Enter() (*testtools.AwsmStubber, context.Context) {
	stubber := testtools.NewStubber()
	os.Setenv("DB_STOCKS_TABLE_NAME", "DB_STOCKS_TABLE_NAME")
	os.Setenv("DB_LOGS_TABLE_NAME", "DB_LOGS_TABLE_NAME")
	os.Setenv("EVENTBRIDGE_RULE_NAME", "EVENTBRIDGE_RULE_NAME")
	os.Setenv("LAMBDA_TICKER_NAME", "LAMBDA_TICKER_NAME")
	os.Setenv("LAMBDA_WORKER_NAME", "LAMBDA_WORKER_NAME")
	os.Setenv("SQS_QUEUE_URL", "SQS_QUEUE_URL")
	os.Setenv("SQS_DLQ_URL", "SQS_DLQ_URL")
	os.Setenv("TICKER_TIMEOUT", "2")

	ctx := lambdacontext.NewContext(context.TODO(), &lambdacontext.LambdaContext{
		AwsRequestID: "test_request",
	})

	return stubber, ctx
}

func AssertAll(t *testing.T, stubber *testtools.AwsmStubber, actualError error, expectedError error) {
	assert.Equal(t, expectedError, actualError)

	testtools.ExitTest(stubber, t)
}
