package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"jon-richards.com/stock-app/internal/db"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
)

type ResponseBody struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

var queueService = jobs.NewQueueService()
var eventsService = jobs.NewEventsService()
var dbService = db.NewDatabaseService()

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	switch request.HTTPMethod {
	case "POST":
		return create(ctx, request)
	default:
		return clientError(ctx, http.StatusMethodNotAllowed, fmt.Errorf("request method %s not supported", request.HTTPMethod))
	}
}

func clientError(ctx context.Context, status int, err error) *events.APIGatewayProxyResponse {
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
	}
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
