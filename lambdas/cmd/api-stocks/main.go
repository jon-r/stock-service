package main

import (
	"context"
	"fmt"
	"net/http"

	awsEvents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/config"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
	"github.com/jon-r/stock-service/lambdas/internal/utils/response"
	"go.uber.org/zap/zapcore"
)

//type ResponseBody struct {
//	Message string `json:"message"`
//	Status  int    `json:"status"`
//}

//type OLD_ApiStockHandler struct {
//	types_old.ServiceHandler
//}

type apiStockHandler interface {
	handleRequest(ctx context.Context, request awsEvents.APIGatewayProxyRequest) (*awsEvents.APIGatewayProxyResponse, error)
}

type handler struct {
	queueBroker     queue.Broker
	eventsScheduler events.Scheduler
	idGen           queue.NewIdFunc
	dbRepository    db.Repository
	log             logger.Logger
}

func newHandler() apiStockHandler {
	cfg := config.GetAwsConfig()
	idGen := uuid.NewString

	return &handler{
		queueBroker:     queue.NewBroker(cfg, idGen),
		eventsScheduler: events.NewScheduler(cfg),
		dbRepository:    db.NewRepository(cfg),
		log:             logger.NewLogger(zapcore.InfoLevel),
		idGen:           idGen,
	}
}

func (h *handler) handleRequest(ctx context.Context, request awsEvents.APIGatewayProxyRequest) (*awsEvents.APIGatewayProxyResponse, error) {
	// todo look at zap docs to see if this can be done better
	h.log = h.log.LoadLambdaContext(ctx)

	switch request.HTTPMethod {
	case http.MethodPost:
		return h.handlePost(request)
	default:
		err := fmt.Errorf("request method %s not supported", request.HTTPMethod)
		h.log.Errorw("Request error",
			"status", http.StatusMethodNotAllowed,
			"message", err,
		)

		return response.StatusMethodNotAllowed(err)
	}
}

//func (handler OLD_ApiStockHandler) handleRequest(ctx context.Context, request awsEvents.APIGatewayProxyRequest) (*awsEvents.APIGatewayProxyResponse, error) {
//	if handler.LogService == nil { // todo this might not work?
//		handler.LogService = logging_old.NewLogger(ctx)
//	}
//	defer handler.LogService.Sync()
//
//	switch request.HTTPMethod {
//	case "POST":
//		return handler.handlePost(request)
//	default:
//		err := fmt.Errorf("request method %s not supported", request.HTTPMethod)
//		handler.LogService.Errorw("Request error",
//			"status", http.StatusMethodNotAllowed,
//			"message", err,
//		)
//		return clientError(http.StatusMethodNotAllowed, err)
//	}
//}

//func clientError(status int, err error) (*awsEvents.APIGatewayProxyResponse, error) {
//	body, _ := json.Marshal(ResponseBody{
//		Message: http.StatusText(status),
//		Status:  status,
//	})
//
//	return &awsEvents.APIGatewayProxyResponse{
//		StatusCode: status,
//		Body:       string(body),
//	}, err
//}

//func clientSuccess(message string) *awsEvents.APIGatewayProxyResponse {
//	if message == "" {
//		message = "Success"
//	}
//
//	body, _ := json.Marshal(ResponseBody{
//		Message: message,
//		Status:  http.StatusOK,
//	})
//
//	return &awsEvents.APIGatewayProxyResponse{
//		StatusCode: http.StatusOK,
//		Body:       string(body),
//	}
//}

//var serviceHandler_old = types_old.ServiceHandler{
//	QueueService:  jobs_old.NewQueueService(jobs_old.CreateSqsClient()),
//	EventsService: scheduler_old.NewEventsService(scheduler_old.CreateEventClients()),
//	DbService:     db_old.NewDatabaseService(db_old.CreateDatabaseClient()),
//	NewUuid:       uuid.NewString,
//}

var serviceHandler = newHandler()

func main() {
	lambda.Start(serviceHandler.handleRequest)
}
