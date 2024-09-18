package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAllTickers(t *testing.T) {
	stubber, ctx := test.SetupLambdaEnvironment()
	mockHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

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

	t.Run("No Errors", func(t *testing.T) {
		expectedTickers := []ticker.Entity{
			{
				EntityBase: db.EntityBase{Id: "T#AMZN", Sort: "T#AMZN"},
				Provider:   provider.PolygonIo,
			},
			{
				EntityBase: db.EntityBase{Id: "T#META", Sort: "T#META"},
				Provider:   provider.PolygonIo,
			},
		}

		stubber.Add(test.StubDynamoDbScan(expectedQuery, expectedTickers, nil))

		expectedQueueItems := []sqsTypes.SendMessageBatchRequestEntry{
			{
				Id:          aws.String("TEST_ID"),
				MessageBody: aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"UPDATE_PRICES","TickerId":"AMZN,META","Attempts":0}`),
			},
		}
		stubber.Add(test.StubSqsSendMessageBatch("SQS_QUEUE_URL", expectedQueueItems, nil))
		stubber.Add(test.StubEventbridgeEnableRule("EVENTBRIDGE_RULE_NAME", nil))
		stubber.Add(test.StubLambdaInvoke("LAMBDA_TICKER_NAME", nil, nil))

		err := mockHandler.HandleRequest(ctx)

		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("No Items", func(t *testing.T) {
		stubber.Add(test.StubDynamoDbScan(expectedQuery, []ticker.Entity{}, nil))

		err := mockHandler.HandleRequest(ctx)
		expectedErr := fmt.Errorf("no items found")

		assert.Equal(t, expectedErr, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("AWS Error", func(t *testing.T) {
		stubber.Add(test.StubDynamoDbScan(expectedQuery, nil, fmt.Errorf("test error")))

		err := mockHandler.HandleRequest(ctx)
		expectedErr := test.StubbedError(fmt.Errorf("test error"))

		testtools.VerifyError(err, expectedErr, t)
		testtools.ExitTest(stubber, t)
	})
}
