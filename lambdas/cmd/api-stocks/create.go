package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/providers"
)

type RequestParams struct {
	Provider providers.ProviderName `json:"provider"`
	TickerId string                 `json:"ticker"`
}

// var dbService = db.NewDatabaseService()
var queueService = jobs.NewQueueService()
var eventsService = jobs.NewEventsService()

func createStockIndex(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error

	// 1. get ticket and provider from post request body
	var params RequestParams
	err = json.Unmarshal([]byte(request.Body), &params)

	if err != nil {
		return clientError(http.StatusBadRequest, err)
	}

	// 2. Create new job queue item
	job := jobs.JobAction{
		Type:     jobs.NewTickerItem,
		Provider: params.Provider,
		TickerId: params.TickerId,
		Attempts: 0,
	}

	err = queueService.AddJobs([]jobs.JobAction{job})

	if err != nil {
		return clientError(http.StatusInternalServerError, err)
	}

	// 3. enable the event timer
	err = eventsService.StartTickerScheduler()

	if err != nil {
		return clientError(http.StatusInternalServerError, err)
	}

	return clientSuccess()
}

func create(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return createStockIndex(request)
}
