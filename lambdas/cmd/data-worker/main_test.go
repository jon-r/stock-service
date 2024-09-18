package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
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
	t.Run("No Errors", setTickerDescriptionNoErrors)
	t.Run("Invalid action type", invalidActionType)
	// other errors are handled in other tests
}

func setTickerDescriptionNoErrors(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

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

	assert.NoError(t, err)
	testtools.ExitTest(stubber, t)
}

func invalidActionType(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}
