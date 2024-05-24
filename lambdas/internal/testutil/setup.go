package testutil

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"github.com/jon-r/stock-service/lambdas/internal/types"
	"go.uber.org/zap"
)

func EnterTest() (*testtools.AwsmStubber, *types.ServiceHandler) {
	stubber := testtools.NewStubber()

	os.Setenv("LAMBDA_TICKER_NAME", "LAMBDA_TICKER_NAME")
	os.Setenv("EVENTBRIDGE_RULE_NAME", "EVENTBRIDGE_RULE_NAME")
	os.Setenv("DB_STOCKS_TABLE_NAME", "DB_STOCKS_TABLE_NAME")

	mockSqsClient := sqs.NewFromConfig(*stubber.SdkConfig)
	mockDbClient := dynamodb.NewFromConfig(*stubber.SdkConfig)
	mockEventsClient := eventbridge.NewFromConfig(*stubber.SdkConfig)
	mockLambdaClient := lambda.NewFromConfig(*stubber.SdkConfig)

	mockHandler := &types.ServiceHandler{
		QueueService:  jobs.NewQueueService(mockSqsClient),
		EventsService: scheduler.NewEventsService(mockEventsClient, mockLambdaClient),
		DbService:     db.NewDatabaseService(mockDbClient),
		LogService:    zap.NewNop().Sugar(),
		NewUuid:       func() string { return "TEST_ID" },
	}
	return stubber, mockHandler
}
