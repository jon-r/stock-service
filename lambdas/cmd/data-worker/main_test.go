package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
)

func TestHandleJobAction(t *testing.T) {
	t.Run("SetTickerDescriptionNoErrors", handleSetTickerDescriptionNoErrors)
	t.Run("SetHistoricalPricesNoErrors", handleSetHistoricalPricesNoErrors)
	t.Run("UpdatePricesNoErrors", handleUpdatePricesNoErrors)
}

func handleSetTickerDescriptionNoErrors(t *testing.T) {
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
	test.Assert(t, stubber, err, nil)
}

func handleSetHistoricalPricesNoErrors(t *testing.T) {
	stubber, ctx := test.Enter()
	mockServiceHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

	var jsonData interface{}
	test.ReadTestJson("./testdata/testTicker1Price.json", &jsonData)
	item1, _ := attributevalue.MarshalMap(jsonData)
	test.ReadTestJson("./testdata/testTicker2Price.json", &jsonData)
	item2, _ := attributevalue.MarshalMap(jsonData)

	expectedInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"DB_STOCKS_TABLE_NAME": {
				{PutRequest: &types.PutRequest{Item: item1}},
				{PutRequest: &types.PutRequest{Item: item2}},
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
	mockHandler := handler{handlers.NewMock(*stubber.SdkConfig)}

	var jsonData interface{}
	test.ReadTestJson("./testdata/testTicker3Price.json", &jsonData)
	item3, _ := attributevalue.MarshalMap(jsonData)
	test.ReadTestJson("./testdata/testTicker4Price.json", &jsonData)
	item4, _ := attributevalue.MarshalMap(jsonData)

	expectedInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"DB_STOCKS_TABLE_NAME": {
				{PutRequest: &types.PutRequest{Item: item3}},
				{PutRequest: &types.PutRequest{Item: item4}},
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

	err := mockHandler.HandleRequest(ctx, jobEvent)
	test.Assert(t, stubber, err, nil)
}
