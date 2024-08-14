package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
)

//type dataManagerHandler interface {
//	HandleRequest(ctx context.Context) error
//}
//
//type handler struct {
//	tickers tickers.Controller
//	jobs    jobs.Controller
//	log     logger.Logger
//}
//
//func newHandler() dataManagerHandler {
//	cfg := config.GetAwsConfig()
//	log := logger.NewLogger(zapcore.InfoLevel)
//	idGen := uuid.NewString
//
//	// todo once tests split up, some of this can be moved to the controller
//	queueBroker := queue.NewBroker(cfg, idGen)
//	eventsScheduler := events.NewScheduler(cfg)
//	dbRepository := db.NewRepository(cfg)
//
//	jobsCtrl := jobs.NewController(queueBroker, eventsScheduler, idGen, log)
//	tickersCtrl := tickers.NewController(dbRepository, nil, log)
//
//	return &handler{tickersCtrl, jobsCtrl, log}
//}

type handler struct{ *handlers.LambdaHandler }

var dataManagerHandler = handler{handlers.NewLambdaHandler()}

func (h *handler) HandleRequest(ctx context.Context) error {
	// todo look at zap docs to see if this can be done better. its not passing context to controllers
	h.Log = h.Log.LoadLambdaContext(ctx)
	defer h.Log.Sync()

	// 1. get all tickers
	tickerList, err := h.Tickers.GetAll()

	if err != nil {
		return err
	}

	// 2. convert the jobs into update actions
	return h.Jobs.LaunchDailyTickerJobs(tickerList)
}

func main() {
	lambda.Start(dataManagerHandler.HandleRequest)
}
