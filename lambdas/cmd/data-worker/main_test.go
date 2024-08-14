package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/providers"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/prices"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"go.uber.org/zap/zapcore"
)

func TestHandleJobAction(t *testing.T) {
	t.Run("SetTickerDescriptionNoErrors", handleSetTickerDescriptionNoErrors)
	t.Run("SetHistoricalPricesNoErrors", handleSetHistoricalPricesNoErrors)
	t.Run("UpdatePricesNoErrors", handleUpdatePricesNoErrors)
}

func mockHandler(cfg aws.Config) dataWorkerHandler {
	//cfg := config.GetAwsConfig()
	//idGen := func() string { return "TEST_ID" }
	log := logger.NewLogger(zapcore.DebugLevel)
	//idGen := uuid.NewString

	// todo once tests split up, some of this can be moved to the controller
	providersService := providers.NewMock()
	queueBroker := queue.NewBroker(cfg, nil)
	//eventsScheduler := events.NewScheduler(cfg)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, nil, nil, log)
	tickersCtrl := tickers.NewController(dbRepository, providersService, log)
	pricesCtrl := prices.NewController(dbRepository, providersService, log)

	return &handler{tickersCtrl, jobsCtrl, pricesCtrl, log}
}

func handleSetTickerDescriptionNoErrors(t *testing.T) {
	//stubber, mockServiceHander := testutil_old.EnterTest(nil)
	//mockHandler := DataWorkerHandler{
	//	*mockServiceHander,
	//	providers_old.NewMock(),
	//}

	stubber, ctx := test.Enter()
	mockServiceHandler := mockHandler(*stubber.SdkConfig)

	expectedUpdate := &dynamodb.UpdateItemInput{
		TableName: aws.String("DB_STOCKS_TABLE_NAME"),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "T#TestTicker"},
			"SK": &types.AttributeValueMemberS{Value: "T#TestTicker"},
		},
		ExpressionAttributeNames: map[string]string{
			"#0": "Description",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":0": &types.AttributeValueMemberM{
				Value: map[string]types.AttributeValue{
					"Currency":   &types.AttributeValueMemberS{Value: "GBP"},
					"FullName":   &types.AttributeValueMemberS{Value: "Full name TestTicker"},
					"FullTicker": &types.AttributeValueMemberS{Value: "Ticker:TestTicker"},
					"Icon":       &types.AttributeValueMemberS{Value: "Icon:POLYGON_IO/TestTicker"},
				},
			},
		},
		UpdateExpression: aws.String("SET #0 = :0\n"),
	}
	stubber.Add(test.StubDynamoDbUpdate(expectedUpdate, nil))

	jobEvent := job.Job{
		JobId:    "TestJob",
		Provider: provider.PolygonIo,
		Type:     job.LoadTickerDescription,
		TickerId: "TestTicker",
		Attempts: 0,
	}
	err := mockServiceHandler.HandleRequest(ctx, jobEvent)
	test.Assert(t, stubber, err, nil)
}

func handleSetHistoricalPricesNoErrors(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := mockHandler(*stubber.SdkConfig)

	//var entity map[string]types.AttributeValue
	//itemTest := json.Unmarshal([]byte(`{"DT"}`))

	expectedInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"DB_STOCKS_TABLE_NAME": {
				{PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"DT": &types.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &types.AttributeValueMemberS{Value: "T#TestTicker"},
						"SK": &types.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Prices": &types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"Open":      &types.AttributeValueMemberN{Value: "10"},
								"Close":     &types.AttributeValueMemberN{Value: "20"},
								"High":      &types.AttributeValueMemberN{Value: "30"},
								"Low":       &types.AttributeValueMemberN{Value: "5"},
								"Id":        &types.AttributeValueMemberS{Value: "TestTicker"},
								"Timestamp": &types.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
				{PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"DT": &types.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &types.AttributeValueMemberS{Value: "T#TestTicker"},
						"SK": &types.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Prices": &types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"Open":      &types.AttributeValueMemberN{Value: "20"},
								"Close":     &types.AttributeValueMemberN{Value: "30"},
								"High":      &types.AttributeValueMemberN{Value: "35"},
								"Low":       &types.AttributeValueMemberN{Value: "15"},
								"Id":        &types.AttributeValueMemberS{Value: "TestTicker"},
								"Timestamp": &types.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
			},
		},
	}
	stubber.Add(test.StubDynamoDbBatchWriteTicker(expectedInput, nil))

	jobEvent := job.Job{
		JobId:    "TestJob",
		Provider: provider.PolygonIo,
		Type:     job.LoadHistoricalPrices,
		TickerId: "TestTicker",
		Attempts: 0,
	}

	err := mockServiceHandler.HandleRequest(ctx, jobEvent)
	test.Assert(t, stubber, err, nil)
}

func handleUpdatePricesNoErrors(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := mockHandler(*stubber.SdkConfig)

	expectedInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"DB_STOCKS_TABLE_NAME": {
				{PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"DT": &types.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &types.AttributeValueMemberS{Value: "T#TestTicker1"},
						"SK": &types.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Prices": &types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"Open":      &types.AttributeValueMemberN{Value: "40"},
								"Close":     &types.AttributeValueMemberN{Value: "50"},
								"High":      &types.AttributeValueMemberN{Value: "55"},
								"Low":       &types.AttributeValueMemberN{Value: "35"},
								"Id":        &types.AttributeValueMemberS{Value: "TestTicker1"},
								"Timestamp": &types.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
				{PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"DT": &types.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &types.AttributeValueMemberS{Value: "T#TestTicker2"},
						"SK": &types.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Prices": &types.AttributeValueMemberM{
							Value: map[string]types.AttributeValue{
								"Open":      &types.AttributeValueMemberN{Value: "40"},
								"Close":     &types.AttributeValueMemberN{Value: "50"},
								"High":      &types.AttributeValueMemberN{Value: "55"},
								"Low":       &types.AttributeValueMemberN{Value: "35"},
								"Id":        &types.AttributeValueMemberS{Value: "TestTicker2"},
								"Timestamp": &types.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
			},
		},
	}
	stubber.Add(test.StubDynamoDbBatchWriteTicker(expectedInput, nil))

	jobEvent := job.Job{
		JobId:    "TestJob",
		Provider: provider.PolygonIo,
		Type:     job.LoadDailyPrices,
		TickerId: "TestTicker1,TestTicker2",
		Attempts: 0,
	}

	err := mockServiceHandler.HandleRequest(ctx, jobEvent)
	test.Assert(t, stubber, err, nil)
}
