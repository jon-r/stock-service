package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/stretchr/testify/assert"
)

func TestSetTickerDescription(t *testing.T) {
	// todo redo other tests like this to dry them out
	stubber, _ := test.Enter()
	apiStubber := test.NewApiStubber()
	mockClock := clock.NewMock()

	mockServiceHandler := handler{handlers.NewMockWithHttpClient(*stubber.SdkConfig, apiStubber.NewTestClient(), mockClock)}

	jobEvent := job.Job{
		JobId:    "TestJob",
		Provider: provider.PolygonIo,
		Type:     job.LoadTickerDescription,
		TickerId: "TestTicker",
		Attempts: 0,
	}

	t.Run("API Error", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL:    "https://api.polygon.io/v3/reference/tickers/TestTicker",
			Input:  "",
			Status: http.StatusUnauthorized,
			Output: test.ReadJsonToString("./testdata/errorRes.json"),
		})

		err := mockServiceHandler.doJob(jobEvent)
		expectedErr := &models.ErrorResponse{
			BaseResponse: models.BaseResponse{
				Status:       "ERROR",
				RequestID:    "85cc09a1ce73359badd942fa78412fba",
				ErrorMessage: "Unknown API Key",
			},
			StatusCode: http.StatusUnauthorized,
		}

		assert.Equal(t, expectedErr, err)
		apiStubber.VerifyAllStubsCalled(t)
		testtools.ExitTest(stubber, t)
	})

	t.Run("AWS Error", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL:    "https://api.polygon.io/v3/reference/tickers/TestTicker",
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getDescriptionRes.json"),
		})

		stubber.Add(test.StubDynamoDbUpdate(
			nil,
			fmt.Errorf("test error"),
		))

		err := mockServiceHandler.doJob(jobEvent)
		expectedError := test.StubbedError(fmt.Errorf("test error"))

		testtools.VerifyError(err, expectedError, t)
		apiStubber.VerifyAllStubsCalled(t)
		testtools.ExitTest(stubber, t)
	})
}

func TestSetHistoricalPrices(t *testing.T) {
	stubber, _ := test.Enter()
	apiStubber := test.NewApiStubber()
	mockClock := clock.NewMock()

	mockServiceHandler := handler{handlers.NewMockWithHttpClient(*stubber.SdkConfig, apiStubber.NewTestClient(), mockClock)}
	mockToday := mockClock.Now()

	startDate, _ := time.Parse(time.DateOnly, "2021-12-01")

	t.Run("No Errors", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/ticker/TestTicker/range/1/day/%v/%v?adjusted=true&sort=desc",
				startDate.UnixMilli(),
				mockToday.UnixMilli(),
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getHistoricalPricesRes.json"),
		})
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/ticker/TestTicker/range/1/day/%v/%v?adjusted=true&sort=desc?cursor=%v",
				startDate.UnixMilli(),
				"1666742400000",
				"bGltaXQ9MiZzb3J0PWFzYw",
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getHistoricalPricesRes2.json"),
		})

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
		apiStubber.VerifyAllStubsCalled(t)
		testtools.ExitTest(stubber, t)
	})

	t.Run("API Error", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/ticker/TestTicker/range/1/day/%v/%v?adjusted=true&sort=desc",
				startDate.UnixMilli(),
				mockToday.UnixMilli(),
			),
			Input:  "",
			Status: http.StatusUnauthorized,
			Output: test.ReadJsonToString("./testdata/errorRes.json"),
		})

		jobEvent := job.Job{
			JobId:    "TestJob",
			Provider: provider.PolygonIo,
			Type:     job.LoadHistoricalPrices,
			TickerId: "TestTicker",
			Attempts: 0,
		}

		err := mockServiceHandler.doJob(jobEvent)
		expectedErr := &models.ErrorResponse{
			BaseResponse: models.BaseResponse{
				Status:       "ERROR",
				RequestID:    "85cc09a1ce73359badd942fa78412fba",
				ErrorMessage: "Unknown API Key",
			},
			StatusCode: http.StatusUnauthorized,
		}

		assert.Equal(t, expectedErr, err)
		apiStubber.VerifyAllStubsCalled(t)
		testtools.ExitTest(stubber, t)
	})

	t.Run("AWS Error", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/ticker/TestTicker/range/1/day/%v/%v?adjusted=true&sort=desc",
				startDate.UnixMilli(),
				mockToday.UnixMilli(),
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getHistoricalPricesRes.json"),
		})
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/ticker/TestTicker/range/1/day/%v/%v?adjusted=true&sort=desc?cursor=%v",
				startDate.UnixMilli(),
				"1666742400000",
				"bGltaXQ9MiZzb3J0PWFzYw",
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getHistoricalPricesRes2.json"),
		})

		stubber.Add(test.StubDynamoDbBatchWriteTicker(
			nil,
			fmt.Errorf("test error"),
		))

		jobEvent := job.Job{
			JobId:    "TestJob",
			Provider: provider.PolygonIo,
			Type:     job.LoadHistoricalPrices,
			TickerId: "TestTicker",
			Attempts: 0,
		}
		err := mockServiceHandler.doJob(jobEvent)
		expectedError := test.StubbedError(fmt.Errorf("test error"))

		testtools.VerifyError(err, expectedError, t)
		apiStubber.VerifyAllStubsCalled(t)
		testtools.ExitTest(stubber, t)
	})
}

func TestUpdatePrices(t *testing.T) {
	stubber, _ := test.Enter()
	apiStubber := test.NewApiStubber()
	mockClock := clock.NewMock()

	mockServiceHandler := handler{handlers.NewMockWithHttpClient(*stubber.SdkConfig, apiStubber.NewTestClient(), mockClock)}
	mockToday := mockClock.Now()

	mockYesterday := mockToday.Add(24 * -time.Hour)

	jobEvent := job.Job{
		JobId:    "TestJob",
		Provider: provider.PolygonIo,
		Type:     job.LoadDailyPrices,
		TickerId: "TestTicker1,TestTicker2",
		Attempts: 0,
	}

	t.Run("No Errors", func(t *testing.T) {
		var jsonData interface{}
		test.ReadTestJson("./testdata/testTicker3Price.json", &jsonData)
		item3, _ := attributevalue.MarshalMap(jsonData)
		test.ReadTestJson("./testdata/testTicker4Price.json", &jsonData)
		item4, _ := attributevalue.MarshalMap(jsonData)

		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/grouped/locale/us/market/stocks/%v?adjusted=true",
				mockYesterday.Format(time.DateOnly),
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getDailyPricesRes.json"),
		})

		expectedInput := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				"DB_STOCKS_TABLE_NAME": {
					{PutRequest: &types.PutRequest{Item: item3}},
					{PutRequest: &types.PutRequest{Item: item4}},
				},
			},
		}
		stubber.Add(test.StubDynamoDbBatchWriteTicker(expectedInput, nil))

		err := mockServiceHandler.doJob(jobEvent)

		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("API Error", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/grouped/locale/us/market/stocks/%v?adjusted=true",
				mockYesterday.Format(time.DateOnly),
			),
			Input:  "",
			Status: http.StatusUnauthorized,
			Output: test.ReadJsonToString("./testdata/errorRes.json"),
		})

		err := mockServiceHandler.doJob(jobEvent)
		expectedErr := &models.ErrorResponse{
			BaseResponse: models.BaseResponse{
				Status:       "ERROR",
				RequestID:    "85cc09a1ce73359badd942fa78412fba",
				ErrorMessage: "Unknown API Key",
			},
			StatusCode: http.StatusUnauthorized,
		}

		assert.Equal(t, expectedErr, err)
		apiStubber.VerifyAllStubsCalled(t)
		testtools.ExitTest(stubber, t)
	})

	t.Run("AWS Error", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/grouped/locale/us/market/stocks/%v?adjusted=true",
				mockYesterday.Format(time.DateOnly),
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getDailyPricesRes.json"),
		})

		stubber.Add(test.StubDynamoDbBatchWriteTicker(
			nil,
			fmt.Errorf("test error"),
		))

		err := mockServiceHandler.doJob(jobEvent)
		expectedError := test.StubbedError(fmt.Errorf("test error"))

		testtools.VerifyError(err, expectedError, t)
		apiStubber.VerifyAllStubsCalled(t)
		testtools.ExitTest(stubber, t)
	})
}

func updatePricesApiError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}

func updatePricesAWSError(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}
