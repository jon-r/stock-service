package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/logging"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"github.com/jon-r/stock-service/lambdas/internal/types"
)

type DataManagerHandler struct {
	types.ServiceHandler
}

func (handler DataManagerHandler) updateAllTickers(ctx context.Context) {
	// todo this might not work?
	if handler.LogService == nil {
		handler.LogService = logging.NewLogger(ctx)
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

	// 2. convert the jobs into update actions
	jobActions := jobs.MakeUpdateJobs(tickers, handler.NewUuid)

	// 3. add queue jobs for ticker prices + dividends
	err = handler.QueueService.AddJobs(*jobActions, handler.NewUuid)

	if err != nil {
		handler.LogService.Fatalw("Failed to add jobs",
			"error", err,
		)
	} else {
		handler.LogService.Infow("Added Jobs for tickers",
			"tickers", tickers,
		)
	}

	// 4. enable the jobs ticker
	err = handler.EventsService.StartTickerScheduler()

	if err != nil {
		handler.LogService.Fatalw("Failed to start the ticker",
			"error", err,
		)
	}
}

var serviceHandler = types.ServiceHandler{
	QueueService:  jobs.NewQueueService(jobs.CreateSqsClient()),
	EventsService: scheduler.NewEventsService(scheduler.CreateEventClients()),
	DbService:     db.NewDatabaseService(db.CreateDatabaseClient()),
	NewUuid:       uuid.NewString,
}

func main() {
	handler := DataManagerHandler{ServiceHandler: serviceHandler}
	lambda.Start(handler.updateAllTickers)
}
