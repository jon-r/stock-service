package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/testutil"
)

func enterTest() (*testtools.AwsmStubber, *db.DatabaseRepository) {
	stubber := testtools.NewStubber()
	repository := &db.DatabaseRepository{
		Svc: dynamodb.NewFromConfig(*stubber.SdkConfig),
	}
}

func TestHandleRequest(t *testing.T) {
	var s3Event events.APIGatewayProxyRequest

	testutil.ReadTestJson("./testevents/api-stocks_POST.json", &s3Event)

	fmt.Printf("%+v+", s3Event)

}