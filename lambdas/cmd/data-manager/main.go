package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"jon-richards.com/stock-app/internal/db"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
	"jon-richards.com/stock-app/internal/scheduler"
)

var dbService = db.NewDatabaseService()
var queueService = jobs.NewQueueService()
var eventsService = scheduler.NewEventsService()

func updateAllTickers(ctx context.Context) {
	var err error

	log := logging.NewLogger(ctx)
	defer log.Sync()

	// 1. get all tickers
	tickers, err := dbService.GetAllTickers()

	if err != nil {
		log.Fatalw("Errors in fetching the tickers",
			"error", err,
		)
	}

	if len(tickers) == 0 {
		log.Fatal("No tickers found")
	}

	// 2. convert the jobs into update actions
	jobActions := jobs.MakeUpdateJobs(tickers, uuid.NewString)

	// 3. add queue jobs for ticker prices + dividends
	err = queueService.AddJobs(*jobActions)

	if err != nil {
		log.Fatalw("Failed to add jobs",
			"error", err,
		)
	} else {
		log.Infow("Added Jobs for tickers",
			"tickers", tickers,
		)
	}

	// 4. enable the jobs ticker
	err = eventsService.StartTickerScheduler()

	if err != nil {
		log.Fatalw("Failed to start the ticker",
			"error", err,
		)
	}
}

func main() {
	lambda.Start(updateAllTickers)
}
