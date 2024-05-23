package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"go.uber.org/zap"
)

type DataManagerHandler struct {
	queueService  *jobs.QueueRepository
	eventsService *scheduler.EventsRepository
	dbService     *db.DatabaseRepository
	log           *zap.SugaredLogger
}

func (handler DataManagerHandler) updateAllTickers(ctx context.Context) {
	// todo this might not work?
	handler.log = logging.NewLogger(ctx)
	defer handler.log.Sync()

	var err error

	// 1. get all tickers
	tickers, err := handler.dbService.GetAllTickers()

	if err != nil {
		handler.log.Fatalw("Errors in fetching the tickers",
			"error", err,
		)
	}

	if len(tickers) == 0 {
		handler.log.Fatal("No tickers found")
	}

	// 2. convert the jobs into update actions
	jobActions := jobs.MakeUpdateJobs(tickers, uuid.NewString)

	// 3. add queue jobs for ticker prices + dividends
	err = handler.queueService.AddJobs(*jobActions)

	if err != nil {
		handler.log.Fatalw("Failed to add jobs",
			"error", err,
		)
	} else {
		handler.log.Infow("Added Jobs for tickers",
			"tickers", tickers,
		)
	}

	// 4. enable the jobs ticker
	err = handler.eventsService.StartTickerScheduler()

	if err != nil {
		handler.log.Fatalw("Failed to start the ticker",
			"error", err,
		)
	}
}

var handler = DataManagerHandler{
	queueService:  jobs.NewQueueService(jobs.CreateSqsClient()),
	eventsService: scheduler.NewEventsService(scheduler.CreateEventClients()),
	dbService:     db.NewDatabaseService(db.CreateDatabaseClient()),
}

func main() {
	lambda.Start(handler.updateAllTickers)
}
