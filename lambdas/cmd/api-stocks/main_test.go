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
	t.Run("No Errors", createTickerNoErrors)
	t.Run("Invalid request method", createTickerInvalidRequestMethod)
	t.Run("Invalid request body", createTickerInvalidRequestBody)
	t.Run("AWS database error", createTickerAwsError)
}

func createTickerNoErrors(t *testing.T) {
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

	res, err := mockServiceHandler.HandleRequest(ctx, postEvent)
	expectedResponse, _ := response.StatusOK("Success: ticker 'AMZN' queued")

	assert.Equal(t, expectedResponse, res)
	assert.NoError(t, err)
	testtools.ExitTest(stubber, t)
}

func createTickerInvalidRequestMethod(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

	var putEvent awsEvents.APIGatewayProxyRequest
	test.ReadTestJson("./testdata/api-stocks_PUT.json", &putEvent)

	res, err := mockServiceHandler.HandleRequest(ctx, putEvent)

	expectedError := fmt.Errorf("request method PUT not supported")
	expectedResponse, _ := response.StatusMethodNotAllowed(expectedError)

	assert.Equal(t, expectedResponse, res)
	assert.Error(t, expectedError, err)
	testtools.ExitTest(stubber, t)
}

func createTickerInvalidRequestBody(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

	emptyEvent := awsEvents.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       "",
	}

	res, err := mockServiceHandler.HandleRequest(ctx, emptyEvent)

	expectedResponse, expectedError := response.StatusBadRequest(err)

	assert.Equal(t, expectedResponse, res)
	assert.Error(t, expectedError, err)
	testtools.ExitTest(stubber, t)
}

func createTickerAwsError(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

	stubber.Add(test.StubDynamoDbPutItem(
		"DB_STOCKS_TABLE_NAME",
		nil,
		fmt.Errorf("test error"),
	))

	var postEvent awsEvents.APIGatewayProxyRequest
	test.ReadTestJson("./testdata/api-stocks_POST.json", &postEvent)

	res, err := mockServiceHandler.HandleRequest(ctx, postEvent)

	expectedResponse, _ := response.StatusServerError(err)
	expectedError := fmt.Errorf("test error")

	assert.Equal(t, expectedResponse, res)
	testtools.VerifyError(err, test.StubbedError(expectedError), t)
	testtools.ExitTest(stubber, t)
}
