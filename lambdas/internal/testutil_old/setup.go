package testutil_old

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/db_old"
	"github.com/jon-r/stock-service/lambdas/internal/jobs_old"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler_old"
	"github.com/jon-r/stock-service/lambdas/internal/types_old"
	"go.uber.org/zap"
)

type TestSettings struct {
	MuteErrors bool
}

func EnterTest(settings *TestSettings) (*testtools.AwsmStubber, *types_old.ServiceHandler) {
	if settings == nil {
		settings = &TestSettings{MuteErrors: false}
	}

	stubber := testtools.NewStubber()

	os.Setenv("DB_STOCKS_TABLE_NAME", "DB_STOCKS_TABLE_NAME")
	os.Setenv("DB_LOGS_TABLE_NAME", "DB_LOGS_TABLE_NAME")
	os.Setenv("EVENTBRIDGE_RULE_NAME", "EVENTBRIDGE_RULE_NAME")
	os.Setenv("LAMBDA_TICKER_NAME", "LAMBDA_TICKER_NAME")
	os.Setenv("LAMBDA_WORKER_NAME", "LAMBDA_WORKER_NAME")
	os.Setenv("SQS_QUEUE_URL", "SQS_QUEUE_URL")
	os.Setenv("SQS_DLQ_URL", "SQS_DLQ_URL")
	os.Setenv("TICKER_TIMEOUT", "2")

	mockSqsClient := sqs.NewFromConfig(*stubber.SdkConfig)
	mockDbClient := dynamodb.NewFromConfig(*stubber.SdkConfig)
	mockEventsClient := eventbridge.NewFromConfig(*stubber.SdkConfig)
	mockLambdaClient := lambda.NewFromConfig(*stubber.SdkConfig)

	var mockLogger *zap.Logger
	if settings.MuteErrors {
		mockLogger = zap.NewNop()
	} else {
		mockLogger = zap.Must(zap.NewDevelopment())
	}

	mockHandler := &types_old.ServiceHandler{
		QueueService:  jobs_old.NewQueueService(mockSqsClient),
		EventsService: scheduler_old.NewEventsService(mockEventsClient, mockLambdaClient),
		DbService:     db_old.NewDatabaseService(mockDbClient),
		LogService:    mockLogger.Sugar(),
		NewUuid:       func() string { return "TEST_ID" },
	}
	return stubber, mockHandler
}

func Assert(stubber *testtools.AwsmStubber, actualError error, expectedError error, t *testing.T) {
	testtools.VerifyError(actualError, StubbedError(expectedError), t)

	testtools.ExitTest(stubber, t)
}
