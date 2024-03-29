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

type RequestParams struct {
	Provider providers.ProviderName `json:"provider"`
	TickerId string                 `json:"ticker"`
}

func createTicker(ctx context.Context, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	log := logging.NewLogger(ctx)
	defer log.Sync()

	var err error

	// 1. get ticket and provider from post request body
	var params RequestParams
	err = json.Unmarshal([]byte(request.Body), &params)

	if err != nil {
		return clientError(ctx, http.StatusBadRequest, err)
	}

	// 2. enter basic content to the database
	err = dbService.NewTickerItem(params.Provider, params.TickerId)

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

	return clientSuccess(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId))
}

func create(ctx context.Context, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	// todo would be switch if multiple endpoints
	return createTicker(ctx, request)
}
