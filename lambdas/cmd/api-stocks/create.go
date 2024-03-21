package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/logging"
	"jon-richards.com/stock-app/internal/providers"
)

type RequestParams struct {
	Provider providers.ProviderName `json:"provider"`
	TickerId string                 `json:"ticker"`
}

func newStockTickerJobs(provider providers.ProviderName, tickerId string) *[]jobs.JobAction {
	newItemActions := []jobs.JobTypes{
		jobs.LoadTickerDescription,
		jobs.LoadHistoricalPrices,
		// TODO jobs.LoadHistoricalDividends,
		// TODO jobs.LoadTickerIcon,
	}

	jobActions := make([]jobs.JobAction, len(newItemActions))
	for i, jobType := range newItemActions {
		job := jobs.JobAction{
			JobId:    uuid.NewString(),
			Provider: provider,
			Type:     jobType,
			TickerId: tickerId,
			Attempts: 0,
		}
		jobActions[i] = job
	}

	return &jobActions
}

func createStockIndex(ctx context.Context, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
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
	newItemJobs := newStockTickerJobs(params.Provider, params.TickerId)

	log.Infow("Add jobs to the queue",
		"jobs", *newItemJobs,
	)
	err = queueService.AddJobs(*newItemJobs)

	if err != nil {
		return clientError(ctx, http.StatusInternalServerError, err)
	}

	// 4. enable the event timer
	err = eventsService.StartTickerScheduler(ctx)

	if err != nil {
		return clientError(ctx, http.StatusInternalServerError, err)
	}

	return clientSuccess(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId))
}

func create(ctx context.Context, request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	// todo would be switch if multiple endpoints
	return createStockIndex(ctx, request)
}
