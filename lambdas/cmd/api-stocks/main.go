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
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
	"github.com/jon-r/stock-service/lambdas/internal/utils/response"
	"go.uber.org/zap/zapcore"
)

type apiStockHandler interface {
	HandleRequest(ctx context.Context, request awsEvents.APIGatewayProxyRequest) (*awsEvents.APIGatewayProxyResponse, error)
}

type handler struct {
	tickers tickers.Controller
	jobs    jobs.Controller
	log     logger.Logger
}

func newHandler() apiStockHandler {
	cfg := config.GetAwsConfig()
	log := logger.NewLogger(zapcore.InfoLevel)
	idGen := uuid.NewString

	// todo once tests split up, some of this can be moved to the controller
	queueBroker := queue.NewBroker(cfg, idGen)
	eventsScheduler := events.NewScheduler(cfg)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, eventsScheduler, idGen, log)
	tickersCtrl := tickers.NewController(dbRepository, nil, log)

	return &handler{tickersCtrl, jobsCtrl, log}
}

func (h *handler) HandleRequest(ctx context.Context, request awsEvents.APIGatewayProxyRequest) (*awsEvents.APIGatewayProxyResponse, error) {
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

var serviceHandler = newHandler()

func main() {
	lambda.Start(serviceHandler.HandleRequest)
}
