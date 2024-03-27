package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
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
	tickers, queueErr := dbService.GetAllTickers()

	if queueErr != nil {
		log.Errorw("Errors in fetching the tickers",
			"error", err,
		)
	}

	// 2. convert the jobs into update actions
	jobActions := jobs.MakeUpdateJobs(tickers)

	// 3. add queue jobs for ticker prices + dividends
	// todo compare jobs here to the jobs made by api controller.
	//   also check the bulk tool if the jobs are valid
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

	// 4. enable the ticker
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
