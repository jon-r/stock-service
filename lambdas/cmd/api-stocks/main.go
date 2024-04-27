package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"go.uber.org/zap"
)

type ResponseBody struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

var queueService = jobs.NewQueueService()
var eventsService = scheduler.NewEventsService()
var dbService = db.NewDatabaseService()

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return create(ctx, request)
	default:
		return clientError(ctx, http.StatusMethodNotAllowed, fmt.Errorf("request method %s not supported", request.HTTPMethod))
	}
}

func clientError(ctx context.Context, status int, err error) (*events.APIGatewayProxyResponse, error) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	log.Errorw("Request error",
		"status", status,
		"message", err,
	)

	body, _ := json.Marshal(ResponseBody{
		Message: http.StatusText(status),
		Status:  status,
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
	}, err
}

func clientSuccess(message string) *events.APIGatewayProxyResponse {
	if message == "" {
		message = "Success"
	}

	body, _ := json.Marshal(ResponseBody{
		Message: message,
		Status:  http.StatusOK,
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}
}

func main() {
	lambda.Start(handleRequest)
}
