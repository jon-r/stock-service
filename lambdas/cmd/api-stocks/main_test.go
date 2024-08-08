package main

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jon-r/stock-service/lambdas/internal/db_old"
	"github.com/jon-r/stock-service/lambdas/internal/providers_old"
	"github.com/jon-r/stock-service/lambdas/internal/testutil_old"
)

func TestHandleRequest(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) { handleRequest(nil, t) })
	t.Run("TestError", func(t *testing.T) { handleRequest(errors.New("TestError"), t) })
}

// todo break this up to have tests that hit every error
func handleRequest(raiseErr error, t *testing.T) {
	stubber, mockServiceHandler := testutil_old.EnterTest(&testutil_old.TestSettings{
		MuteErrors: raiseErr != nil,
	})
	mockHandler := ApiStockHandler{*mockServiceHandler}

	expectedTicker := db_old.TickerItem{
		StocksTableItem: db_old.StocksTableItem{Id: "T#AMZN", Sort: "T#AMZN"},
		Provider:        providers_old.PolygonIo,
	}
	item, _ := attributevalue.MarshalMap(expectedTicker)
	stubber.Add(testutil_old.StubDynamoDbAddTicker("DB_STOCKS_TABLE_NAME", item, raiseErr))

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
	stubber.Add(testutil_old.StubSqsSendMessageBatch("SQS_QUEUE_URL", expectedQueueItems, raiseErr))

	expectedRule := "EVENTBRIDGE_RULE_NAME"
	stubber.Add(testutil_old.StubEventbridgeEnableRule(expectedRule, raiseErr))
	expectedLambda := "LAMBDA_TICKER_NAME"
	stubber.Add(testutil_old.StubLambdaInvoke(expectedLambda, nil, raiseErr))

	var postEvent events.APIGatewayProxyRequest
	testutil_old.ReadTestJson("./testevents/api-stocks_POST.json", &postEvent)

	_, err := mockHandler.handleRequest(context.TODO(), postEvent)

	testutil_old.Assert(stubber, err, raiseErr, t)
}
