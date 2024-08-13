package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"go.uber.org/zap/zapcore"
)

func TestUpdateAllTickers(t *testing.T) {
	t.Run("NoErrors", updateAllTickerNoErrors)
}

func mockHandler(cfg aws.Config) dataManagerHandler {
	idGen := func() string { return "TEST_ID" }
	log := logger.NewLogger(zapcore.DPanicLevel)

	// todo once tests split up, some of this can be moved to the controller
	queueBroker := queue.NewBroker(cfg, idGen)
	eventsScheduler := events.NewScheduler(cfg)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, eventsScheduler, idGen, log)
	tickersCtrl := tickers.NewController(dbRepository, log)

	return &handler{tickersCtrl, jobsCtrl, log}
}

func updateAllTickerNoErrors(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := mockHandler(*stubber.SdkConfig)

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
	stubber.Add(test.StubDynamoDbScan(expectedQuery, expectedTickers, nil))

	expectedQueueItems := []sqsTypes.SendMessageBatchRequestEntry{
		{
			Id:          aws.String("TEST_ID"),
			MessageBody: aws.String(`{"JobId":"TEST_ID","Provider":"POLYGON_IO","Type":"UPDATE_PRICES","TickerId":"AMZN,META","Attempts":0}`),
		},
	}
	stubber.Add(test.StubSqsSendMessageBatch("SQS_QUEUE_URL", expectedQueueItems, nil))

	expectedRule := "EVENTBRIDGE_RULE_NAME"
	stubber.Add(test.StubEventbridgeEnableRule(expectedRule, nil))
	expectedLambda := "LAMBDA_TICKER_NAME"
	stubber.Add(test.StubLambdaInvoke(expectedLambda, nil, nil))

	err := mockServiceHandler.HandleRequest(ctx)

	test.Assert(t, stubber, err, nil)
}
