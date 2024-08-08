package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jon-r/stock-service/lambdas/internal/db_old"
	"github.com/jon-r/stock-service/lambdas/internal/providers_old"
	"github.com/jon-r/stock-service/lambdas/internal/testutil_old"
)

func TestUpdateAllTickers(t *testing.T) {
	t.Run("NoErrors", updateAllTickerNoErrors)
}

func updateAllTickerNoErrors(t *testing.T) {
	stubber, mockServiceHandler := testutil_old.EnterTest(nil)
	mockHandler := DataManagerHandler{*mockServiceHandler}

	expectedTickers := []db_old.TickerItem{
		{
			StocksTableItem: db_old.StocksTableItem{Id: "T#AMZN", Sort: "T#AMZN"},
			Provider:        providers_old.PolygonIo,
		},
		{
			StocksTableItem: db_old.StocksTableItem{Id: "T#META", Sort: "T#META"},
			Provider:        providers_old.PolygonIo,
		},
	}
	expectedQuery := &dynamodb.ScanInput{
		TableName: aws.String("DB_STOCKS_TABLE_NAME"),
		ExpressionAttributeNames: map[string]string{
			"#0": "SK",
			"#1": "Provider",
		},
		ExpressionAttributeValues: map[string]dbTypes.AttributeValue{
			":0": &dbTypes.AttributeValueMemberS{Value: "T#"},
		},
		FilterExpression:     aws.String("begins_with (#0, :0)"),
		ProjectionExpression: aws.String("#0, #1"),
	}
	stubber.Add(testutil_old.StubDynamoDbScan(expectedQuery, expectedTickers, nil))

	expectedQueueItems := []sqsTypes.SendMessageBatchRequestEntry{
		{
			Id:          aws.String("TEST_ID"),
			MessageBody: aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"UPDATE_PRICES","TickerId":"AMZN,META","Attempts":0}`),
		},
	}
	stubber.Add(testutil_old.StubSqsSendMessageBatch("SQS_QUEUE_URL", expectedQueueItems, nil))

	expectedRule := "EVENTBRIDGE_RULE_NAME"
	stubber.Add(testutil_old.StubEventbridgeEnableRule(expectedRule, nil))
	expectedLambda := "LAMBDA_TICKER_NAME"
	stubber.Add(testutil_old.StubLambdaInvoke(expectedLambda, nil, nil))

	err := mockHandler.updateAllTickers(context.TODO())

	testutil_old.Assert(stubber, err, nil, t)
}
