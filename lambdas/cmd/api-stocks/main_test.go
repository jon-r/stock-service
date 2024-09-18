package main

import (
	"fmt"
	"testing"

	awsEvents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/response"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicker(t *testing.T) {
	stubber, ctx := test.SetupLambdaEnvironment()
	mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

	var postEvent awsEvents.APIGatewayProxyRequest
	test.ReadTestJson("./testdata/api-stocks_POST.json", &postEvent)

	var putEvent awsEvents.APIGatewayProxyRequest
	test.ReadTestJson("./testdata/api-stocks_PUT.json", &putEvent)

	emptyEvent := awsEvents.APIGatewayProxyRequest{HTTPMethod: "POST", Body: ""}

	t.Run("No Errors", func(t *testing.T) {
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
		stubber.Add(test.StubEventbridgeEnableRule("EVENTBRIDGE_RULE_NAME", nil))
		stubber.Add(test.StubLambdaInvoke("LAMBDA_TICKER_NAME", nil, nil))

		res, err := mockServiceHandler.HandleRequest(ctx, postEvent)
		expectedResponse, _ := response.StatusOK("Success: ticker 'AMZN' queued")

		assert.Equal(t, expectedResponse, res)
		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Invalid request method", func(t *testing.T) {

		res, err := mockServiceHandler.HandleRequest(ctx, putEvent)

		expectedError := fmt.Errorf("request method PUT not supported")
		expectedResponse, _ := response.StatusMethodNotAllowed(expectedError)

		assert.Equal(t, expectedResponse, res)
		assert.Equal(t, expectedError, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		res, err := mockServiceHandler.HandleRequest(ctx, emptyEvent)

		expectedResponse, expectedError := response.StatusBadRequest(err)

		assert.Equal(t, expectedResponse, res)
		assert.Equal(t, expectedError, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("AWS database error", func(t *testing.T) {
		stubber.Add(test.StubDynamoDbPutItem(
			"DB_STOCKS_TABLE_NAME",
			nil,
			fmt.Errorf("test error"),
		))

		res, err := mockServiceHandler.HandleRequest(ctx, postEvent)

		expectedResponse, _ := response.StatusServerError(err)
		expectedError := test.StubbedError(fmt.Errorf("test error"))

		assert.Equal(t, expectedResponse, res)
		testtools.VerifyError(err, expectedError, t)
		testtools.ExitTest(stubber, t)
	})
}
