package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/providers"
	"github.com/jon-r/stock-service/lambdas/internal/testutil"
)

func TestUpdateAllTickers(t *testing.T) {
	t.Run("NoErrors", updateAllTickerNoErrors)
}

func updateAllTickerNoErrors(t *testing.T) {
	stubber, mockServiceHandler := testutil.EnterTest(nil)
	mockHandler := DataManagerHandler{*mockServiceHandler}

	expectedTickers := []db.TickerItem{
		{
			StocksTableItem: db.StocksTableItem{Id: "T#AMZN", Sort: "T#AMZN"},
			Provider:        providers.PolygonIo,
		},
		{
			StocksTableItem: db.StocksTableItem{Id: "T#META", Sort: "T#META"},
			Provider:        providers.PolygonIo,
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
	stubber.Add(testutil.StubDynamoDbScan(expectedQuery, expectedTickers, nil))

	expectedQueueItems := []sqsTypes.SendMessageBatchRequestEntry{
		{
			Id:          aws.String("TEST_ID"),
			MessageBody: aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"UPDATE_PRICES","TickerId":"AMZN,META","Attempts":0}`),
		},
	}
	stubber.Add(testutil.StubSqsSendMessageBatch("", expectedQueueItems, nil))

	expectedRule := "EVENTBRIDGE_RULE_NAME"
	stubber.Add(testutil.StubEventbridgeEnableRule(expectedRule, nil))
	expectedLambda := "LAMBDA_TICKER_NAME"
	stubber.Add(testutil.StubLambdaInvoke(expectedLambda, nil, nil))

	err := mockHandler.updateAllTickers(context.TODO())

	testutil.Assert(stubber, err, nil, t)
}
