package main

import (
	"testing"

	awsEvents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
)

func TestHandleRequest(t *testing.T) {
	t.Run("CreateTicker - No Errors", handleCreateTicker)
}

func handleCreateTicker(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

	expectedTicker := ticker.Entity{
		EntityBase: db.EntityBase{Id: "T#AMZN", Sort: "T#AMZN"},
		Provider:   provider.PolygonIo,
	}
	stubber.Add(test.StubDynamoDbPutItem("DB_STOCKS_TABLE_NAME", expectedTicker, nil))

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
	stubber.Add(test.StubSqsSendMessageBatch("SQS_QUEUE_URL", expectedQueueItems, nil))

	expectedRule := "EVENTBRIDGE_RULE_NAME"
	stubber.Add(test.StubEventbridgeEnableRule(expectedRule, nil))
	expectedLambda := "LAMBDA_TICKER_NAME"
	stubber.Add(test.StubLambdaInvoke(expectedLambda, nil, nil))

	var postEvent awsEvents.APIGatewayProxyRequest
	test.ReadTestJson("./testdata/api-stocks_POST.json", &postEvent)

	_, err := mockServiceHandler.HandleRequest(ctx, postEvent)

	test.Assert(t, stubber, err, nil)
}
