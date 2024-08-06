package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/providers"
	"github.com/jon-r/stock-service/lambdas/internal/testutil"
)

func TestHandleJobAction(t *testing.T) {
	t.Run("SetTickerDescriptionNoErrors", handleSetTickerDescriptionNoErrors)
	t.Run("SetHistoricalPricesNoErrors", handleSetHistoricalPricesNoErrors)
	t.Run("UpdatePricesNoErrors", handleUpdatePricesNoErrors)
}

func handleSetTickerDescriptionNoErrors(t *testing.T) {
	stubber, mockServiceHander := testutil.EnterTest(nil)
	mockHandler := DataWorkerHandler{
		*mockServiceHander,
		providers.NewMockProviderService(),
	}

	expectedUpdate := &dynamodb.UpdateItemInput{
		TableName: aws.String("DB_STOCKS_TABLE_NAME"),
		Key: map[string]dbTypes.AttributeValue{
			"PK": &dbTypes.AttributeValueMemberS{Value: "T#TestTicker"},
			"SK": &dbTypes.AttributeValueMemberS{Value: "T#TestTicker"},
		},
		ExpressionAttributeNames: map[string]string{
			"#0": "Description",
		},
		ExpressionAttributeValues: map[string]dbTypes.AttributeValue{
			":0": &dbTypes.AttributeValueMemberM{
				Value: map[string]dbTypes.AttributeValue{
					"Currency":   &dbTypes.AttributeValueMemberS{Value: "GBP"},
					"FullName":   &dbTypes.AttributeValueMemberS{Value: "Full name TestTicker"},
					"FullTicker": &dbTypes.AttributeValueMemberS{Value: "Ticker:TestTicker"},
					"Icon":       &dbTypes.AttributeValueMemberS{Value: "Icon:POLYGON_IO/TestTicker"},
				},
			},
		},
		UpdateExpression: aws.String("SET #0 = :0\n"),
	}
	stubber.Add(testutil.StubDynamoDbUpdate(expectedUpdate, nil))

	job := jobs.JobAction{
		JobId:    "TestJob",
		Provider: providers.PolygonIo,
		Type:     jobs.LoadTickerDescription,
		TickerId: "TestTicker",
		Attempts: 0,
	}
	err := mockHandler.handleJobAction(job)

	testutil.Assert(stubber, err, nil, t)
}
func handleSetHistoricalPricesNoErrors(t *testing.T) {
	stubber, mockServiceHander := testutil.EnterTest(nil)
	mockHandler := DataWorkerHandler{
		*mockServiceHander,
		providers.NewMockProviderService(),
	}

	expectedInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dbTypes.WriteRequest{
			"DB_STOCKS_TABLE_NAME": {
				{PutRequest: &dbTypes.PutRequest{
					Item: map[string]dbTypes.AttributeValue{
						"DT": &dbTypes.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &dbTypes.AttributeValueMemberS{Value: "T#TestTicker:POLYGON_IO"},
						"SK": &dbTypes.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Price": &dbTypes.AttributeValueMemberM{
							Value: map[string]dbTypes.AttributeValue{
								"Open":      &dbTypes.AttributeValueMemberN{Value: "10"},
								"Close":     &dbTypes.AttributeValueMemberN{Value: "20"},
								"High":      &dbTypes.AttributeValueMemberN{Value: "30"},
								"Low":       &dbTypes.AttributeValueMemberN{Value: "5"},
								"Id":        &dbTypes.AttributeValueMemberS{Value: "TestTicker:POLYGON_IO"},
								"Timestamp": &dbTypes.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
				{PutRequest: &dbTypes.PutRequest{
					Item: map[string]dbTypes.AttributeValue{
						"DT": &dbTypes.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &dbTypes.AttributeValueMemberS{Value: "T#TestTicker:POLYGON_IO"},
						"SK": &dbTypes.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Price": &dbTypes.AttributeValueMemberM{
							Value: map[string]dbTypes.AttributeValue{
								"Open":      &dbTypes.AttributeValueMemberN{Value: "20"},
								"Close":     &dbTypes.AttributeValueMemberN{Value: "30"},
								"High":      &dbTypes.AttributeValueMemberN{Value: "35"},
								"Low":       &dbTypes.AttributeValueMemberN{Value: "15"},
								"Id":        &dbTypes.AttributeValueMemberS{Value: "TestTicker:POLYGON_IO"},
								"Timestamp": &dbTypes.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
			},
		},
	}
	stubber.Add(testutil.StubDynamoDbBatchWriteTicker(expectedInput, nil))

	job := jobs.JobAction{
		JobId:    "TestJob",
		Provider: providers.PolygonIo,
		Type:     jobs.LoadHistoricalPrices,
		TickerId: "TestTicker",
		Attempts: 0,
	}
	err := mockHandler.handleJobAction(job)

	testutil.Assert(stubber, err, nil, t)
}
func handleUpdatePricesNoErrors(t *testing.T) {
	stubber, mockServiceHander := testutil.EnterTest(nil)
	mockHandler := DataWorkerHandler{
		*mockServiceHander,
		providers.NewMockProviderService(),
	}

	expectedInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dbTypes.WriteRequest{
			"DB_STOCKS_TABLE_NAME": {
				{PutRequest: &dbTypes.PutRequest{
					Item: map[string]dbTypes.AttributeValue{
						"DT": &dbTypes.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &dbTypes.AttributeValueMemberS{Value: "T#TestTicker1:POLYGON_IO"},
						"SK": &dbTypes.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Price": &dbTypes.AttributeValueMemberM{
							Value: map[string]dbTypes.AttributeValue{
								"Open":      &dbTypes.AttributeValueMemberN{Value: "40"},
								"Close":     &dbTypes.AttributeValueMemberN{Value: "50"},
								"High":      &dbTypes.AttributeValueMemberN{Value: "55"},
								"Low":       &dbTypes.AttributeValueMemberN{Value: "35"},
								"Id":        &dbTypes.AttributeValueMemberS{Value: "TestTicker1:POLYGON_IO"},
								"Timestamp": &dbTypes.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
				{PutRequest: &dbTypes.PutRequest{
					Item: map[string]dbTypes.AttributeValue{
						"DT": &dbTypes.AttributeValueMemberS{Value: "-62135596800000"},
						"PK": &dbTypes.AttributeValueMemberS{Value: "T#TestTicker2:POLYGON_IO"},
						"SK": &dbTypes.AttributeValueMemberS{Value: "P#-62135596800000"},
						"Price": &dbTypes.AttributeValueMemberM{
							Value: map[string]dbTypes.AttributeValue{
								"Open":      &dbTypes.AttributeValueMemberN{Value: "40"},
								"Close":     &dbTypes.AttributeValueMemberN{Value: "50"},
								"High":      &dbTypes.AttributeValueMemberN{Value: "55"},
								"Low":       &dbTypes.AttributeValueMemberN{Value: "35"},
								"Id":        &dbTypes.AttributeValueMemberS{Value: "TestTicker2:POLYGON_IO"},
								"Timestamp": &dbTypes.AttributeValueMemberS{Value: "0001-01-01T00:00:00Z"},
							}},
					},
				}},
			},
		},
	}
	stubber.Add(testutil.StubDynamoDbBatchWriteTicker(expectedInput, nil))

	job := jobs.JobAction{
		JobId:    "TestJob",
		Provider: providers.PolygonIo,
		Type:     jobs.UpdatePrices,
		TickerId: "TestTicker1,TestTicker2",
		Attempts: 0,
	}
	err := mockHandler.handleJobAction(job)

	testutil.Assert(stubber, err, nil, t)
}
