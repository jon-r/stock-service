package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
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

type ApiStockHandler struct {
	queueService  *jobs.QueueRepository
	eventsService *scheduler.EventsRepository
	dbService     *db.DatabaseRepository
	logService    *zap.SugaredLogger
	newUuid       jobs.UuidGen
}

func (handler ApiStockHandler) handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if handler.logService == nil { // todo this might not work?
		handler.logService = logging.NewLogger(ctx)
	}
	defer handler.logService.Sync()

	switch request.HTTPMethod {
	case "POST":
		return handler.create(request)
	default:
		err := fmt.Errorf("request method %s not supported", request.HTTPMethod)
		handler.logService.Errorw("Request error",
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

var handler = ApiStockHandler{
	queueService:  jobs.NewQueueService(jobs.CreateSqsClient()),
	eventsService: scheduler.NewEventsService(scheduler.CreateEventClients()),
	dbService:     db.NewDatabaseService(db.CreateDatabaseClient()),
	newUuid:       uuid.NewString,
}

func main() {
	lambda.Start(handler.handleRequest)
}
