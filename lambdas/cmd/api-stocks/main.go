package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/db_old"
	"github.com/jon-r/stock-service/lambdas/internal/jobs_old"
	"github.com/jon-r/stock-service/lambdas/internal/logging_old"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler_old"
	"github.com/jon-r/stock-service/lambdas/internal/types_old"
)

type ResponseBody struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ApiStockHandler struct {
	types_old.ServiceHandler
}

func (handler ApiStockHandler) handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if handler.LogService == nil { // todo this might not work?
		handler.LogService = logging_old.NewLogger(ctx)
	}
	defer handler.LogService.Sync()

	switch request.HTTPMethod {
	case "POST":
		return handler.create(request)
	default:
		err := fmt.Errorf("request method %s not supported", request.HTTPMethod)
		handler.LogService.Errorw("Request error",
			"status", http.StatusMethodNotAllowed,
			"message", err,
		)
		return clientError(http.StatusMethodNotAllowed, err)
	}
}

func clientError(status int, err error) (*events.APIGatewayProxyResponse, error) {
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

var serviceHandler = types_old.ServiceHandler{
	QueueService:  jobs_old.NewQueueService(jobs_old.CreateSqsClient()),
	EventsService: scheduler_old.NewEventsService(scheduler_old.CreateEventClients()),
	DbService:     db_old.NewDatabaseService(db_old.CreateDatabaseClient()),
	NewUuid:       uuid.NewString,
}

func main() {
	handler := ApiStockHandler{ServiceHandler: serviceHandler}

	lambda.Start(handler.handleRequest)
}
