package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/stretchr/testify/assert"
)

func TestSetTickerDescription(t *testing.T) {
	t.Run("API Error", setTickerDescriptionApiError)
	t.Run("AWS Error", setTickerDescriptionAWSError)
}

func setTickerDescriptionApiError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}

func setTickerDescriptionAWSError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}

func TestSetHistoricalPrices(t *testing.T) {
	t.Run("No Errors", setHistoricalPricesNoErrors)
	t.Run("API Error", setHistoricalPricesApiError)
	t.Run("AWS Error", setHistoricalPricesAWSError)
}

func setHistoricalPricesNoErrors(t *testing.T) {
	stubber, _ := test.Enter()
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

	err := mockServiceHandler.doJob(jobEvent)

	assert.NoError(t, err)
	testtools.ExitTest(stubber, t)
}

func setHistoricalPricesApiError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}

func setHistoricalPricesAWSError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}

func TestUpdatePrices(t *testing.T) {
	t.Run("NoErrors", updatePricesNoErrors)
	t.Run("API Error", updatePricesApiError)
	t.Run("AWS Error", updatePricesAWSError)
}

func updatePricesNoErrors(t *testing.T) {
	stubber, _ := test.Enter()
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

	err := mockHandler.doJob(jobEvent)

	assert.NoError(t, err)
	testtools.ExitTest(stubber, t)
}

func updatePricesApiError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}

func updatePricesAWSError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}
