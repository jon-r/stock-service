package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/stretchr/testify/assert"
)

// TODO TEST -> handle errors!!
// also redo the provider stubs to they are setup for each test?
// see if this pushes coverage above 80 consistently

func TestHandleRequest(t *testing.T) {
	t.Run("No Errors", func(t *testing.T) {
		stubber, ctx := test.Enter()
		apiStubber := test.NewApiStubber()
		mockClock := clock.NewMock()

		mockServiceHandler := handler{handlers.NewMockWithHttpClient(*stubber.SdkConfig, apiStubber.NewTestClient(), mockClock)}

		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL:    "https://api.polygon.io/v3/reference/tickers/TestTicker",
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getDescriptionRes.json"),
		})

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
						"Currency":   &types.AttributeValueMemberS{Value: "usd"},
						"FullName":   &types.AttributeValueMemberS{Value: "Apple Inc."},
						"FullTicker": &types.AttributeValueMemberS{Value: "XNAS:AAPL"},
						"Icon":       &types.AttributeValueMemberS{Value: "https://api.polygon.io/v1/reference/company-branding/d3d3LmFwcGxlLmNvbQ/images/2022-01-10_icon.png"},
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

		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Invalid action type", func(t *testing.T) {
		stubber, ctx := test.Enter()

		mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

		stubber.Add(test.StubSqsSendMessage(
			"SQS_QUEUE_URL",
			`{"JobId":"","Provider":"","Type":"","TickerId":"","Attempts":1}`,
			nil,
		))

		jobEvent := job.Job{}

		err := mockServiceHandler.HandleRequest(ctx, jobEvent)
		expectedErr := fmt.Errorf("invalid action type = ")

		assert.Equal(t, expectedErr, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("AWS error", func(t *testing.T) {
		stubber, ctx := test.Enter()

		mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

		stubber.Add(test.StubSqsSendMessage(
			"SQS_QUEUE_URL",
			`{"JobId":"","Provider":"","Type":"","TickerId":"","Attempts":1}`,
			fmt.Errorf("test error"),
		))

		jobEvent := job.Job{}

		err := mockServiceHandler.HandleRequest(ctx, jobEvent)
		expectedError := test.StubbedError(fmt.Errorf("test error"))

		testtools.VerifyError(err, expectedError, t)
		testtools.ExitTest(stubber, t)
	})
}
