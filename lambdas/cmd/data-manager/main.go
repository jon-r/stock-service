package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/db_old"
	"github.com/jon-r/stock-service/lambdas/internal/jobs_old"
	"github.com/jon-r/stock-service/lambdas/internal/logging_old"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler_old"
	"github.com/jon-r/stock-service/lambdas/internal/types_old"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

type DataManagerHandler struct {
	types_old.ServiceHandler
}

type dataManagerHandler interface {
	HandleRequest(ctx context.Context) error
}

// remove any unused things
type handler struct {
	queueBroker     queue.Broker
	eventsScheduler events.Scheduler
	idGen           queue.NewIdFunc
	dbRepository    db.Repository
	log             logger.Logger
}

func (h *handler) updateAllTickers() error {
	var err error

	// 1. get all tickers

	// 2. convert the jobs_old into update actions
	// 3. add queue jobs_old for ticker prices + dividends
	// 4. enable the jobs_old ticker
}

func (h *handler) HandleRequest(ctx context.Context) error {
	// todo look at zap docs to see if this can be done better
	h.log = h.log.LoadLambdaContext(ctx)

	return h.updateAllTickers()
}

func (handler DataManagerHandler) updateAllTickers(ctx context.Context) error {
	// todo this might not work?
	if handler.LogService == nil {
		handler.LogService = logging_old.NewLogger(ctx)
	}
	defer handler.LogService.Sync()

	var err error

	// 1. get all tickers
	tickers, err := handler.DbService.GetAllTickers()

	if err != nil {
		handler.LogService.Fatalw("Errors in fetching the tickers",
			"error", err,
		)
	}

	if len(tickers) == 0 {
		handler.LogService.Fatal("No tickers found")
	}

	// 2. convert the jobs_old into update actions
	jobActions := jobs_old.MakeUpdateJobs(tickers, handler.NewUuid)

	// 3. add queue jobs_old for ticker prices + dividends
	err = handler.QueueService.AddJobs(*jobActions, handler.NewUuid)

	if err != nil {
		handler.LogService.Fatalw("Failed to add jobs_old",
			"error", err,
		)
	} else {
		handler.LogService.Infow("Added Jobs for tickers",
			"tickers", tickers,
		)
	}

	// 4. enable the jobs_old ticker
	err = handler.EventsService.StartTickerScheduler()

	if err != nil {
		handler.LogService.Fatalw("Failed to start the ticker",
			"error", err,
		)
	}

	return err
}

var serviceHandler = types_old.ServiceHandler{
	QueueService:  jobs_old.NewQueueService(jobs_old.CreateSqsClient()),
	EventsService: scheduler_old.NewEventsService(scheduler_old.CreateEventClients()),
	DbService:     db_old.NewDatabaseService(db_old.CreateDatabaseClient()),
	NewUuid:       uuid.NewString,
}

func main() {
	handler := DataManagerHandler{ServiceHandler: serviceHandler}
	lambda.Start(handler.updateAllTickers)
}
