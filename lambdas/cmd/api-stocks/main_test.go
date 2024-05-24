package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/providers"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"github.com/jon-r/stock-service/lambdas/internal/testutil"
	"go.uber.org/zap"
)

func enterTest() (*testtools.AwsmStubber, *ApiStockHandler) {
	stubber := testtools.NewStubber()

	os.Setenv("LAMBDA_TICKER_NAME", "LAMBDA_TICKER_NAME")
	os.Setenv("EVENTBRIDGE_RULE_NAME", "EVENTBRIDGE_RULE_NAME")
	os.Setenv("DB_STOCKS_TABLE_NAME", "DB_STOCKS_TABLE_NAME")

	mockSqsClient := sqs.NewFromConfig(*stubber.SdkConfig)
	mockDbClient := dynamodb.NewFromConfig(*stubber.SdkConfig)
	mockEventsClient := eventbridge.NewFromConfig(*stubber.SdkConfig)
	mockLambdaClient := lambda.NewFromConfig(*stubber.SdkConfig)

	mockHandler := &ApiStockHandler{
		queueService:  jobs.NewQueueService(mockSqsClient),
		eventsService: scheduler.NewEventsService(mockEventsClient, mockLambdaClient),
		dbService:     db.NewDatabaseService(mockDbClient),
		logService:    zap.NewNop().Sugar(),
		newUuid:       func() string { return "TEST_ID" },
	}
	return stubber, mockHandler
}

func TestHandleRequest(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) { handleRequest(nil, t) })
	t.Run("TestError", func(t *testing.T) { handleRequest(&testtools.StubError{Err: errors.New("TestError")}, t) })
}

func handleRequest(raiseErr *testtools.StubError, t *testing.T) {
	stubber, mockHandler := enterTest()

	expectedTicker := db.TickerItem{
		StocksTableItem: db.StocksTableItem{Id: "T#AMZN", Sort: "T#AMZN"},
		Provider:        providers.PolygonIo,
	}
	item, _ := attributevalue.MarshalMap(expectedTicker)
	stubber.Add(testutil.StubDynamoDbAddTicker("DB_STOCKS_TABLE_NAME", item, raiseErr))

	expectedQueueItems := []types.SendMessageBatchRequestEntry{
		{
			Id:          aws.String("TEST_ID"),
			MessageBody: aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_TICKER_DESCRIPTION","TickerId":"AMZN","Attempts":0}`),
		},
		{
			Id:          aws.String("TEST_ID"),
			MessageBody: aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"LOAD_HISTORICAL_PRICES","TickerId":"AMZN","Attempts":0}`),
		},
	}
	stubber.Add(testutil.StubSqsSendMessageBatch("", expectedQueueItems, raiseErr))

	expectedRule := "EVENTBRIDGE_RULE_NAME"
	stubber.Add(testutil.StubEventbridgeEnableRule(expectedRule, raiseErr))
	expectedLambda := "LAMBDA_TICKER_NAME"
	stubber.Add(testutil.StubLambdaInvoke(expectedLambda, raiseErr))

	var postEvent events.APIGatewayProxyRequest
	testutil.ReadTestJson("./testevents/api-stocks_POST.json", &postEvent)

	_, err := mockHandler.handleRequest(context.TODO(), postEvent)

	testtools.VerifyError(err, raiseErr, t)

	testtools.ExitTest(stubber, t)
}
