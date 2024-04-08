package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
	"jon-richards.com/stock-app/internal/providers"
)

func createTicker(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	var err error

	// 1. get ticket and provider from post request body
	var params providers.NewTickerParams
	err = json.Unmarshal([]byte(request.Body), &params)

	if err != nil {
		return clientError(ctx, http.StatusBadRequest, err)
	}

	// 2. enter basic content to the database
	err = dbService.NewTickerItem(log, params)

	if err != nil {
		return clientError(ctx, http.StatusInternalServerError, err)
	}

	// 3. Create new job queue item
	newItemJobs := jobs.MakeCreateJobs(params.Provider, params.TickerId)

	log.Infow("Add jobs to the queue",
		"jobs", *newItemJobs,
	)
	err = queueService.AddJobs(*newItemJobs)

	if err != nil {
		return clientError(ctx, http.StatusInternalServerError, err)
	}

	// 4. enable the event timer
	err = eventsService.StartTickerScheduler()

	if err != nil {
		return clientError(ctx, http.StatusInternalServerError, err)
	}

	return clientSuccess(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId)), nil
}

func create(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return createTicker(ctx, request)
}
