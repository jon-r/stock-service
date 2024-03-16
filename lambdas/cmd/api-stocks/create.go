package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"jon-richards.com/stock-app/internal/jobs"
	"jon-richards.com/stock-app/internal/providers"
)

type RequestParams struct {
	Provider providers.ProviderName `json:"provider"`
	TickerId string                 `json:"ticker"`
}

func createStockIndex(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error

	// 1. get ticket and provider from post request body
	var params RequestParams
	err = json.Unmarshal([]byte(request.Body), &params)

	if err != nil {
		return clientError(http.StatusBadRequest, err)
	}

	// 4. enter basic content to the database
	err = dbService.NewTickerItem(params.Provider, params.TickerId)

	if err != nil {
		return clientError(http.StatusInternalServerError, err)
	}

	// 3. Create new job queue item
	jobId := uuid.NewString()
	job := jobs.JobAction{
		JobId:    jobId,
		Type:     jobs.NewTickerItem,
		Provider: params.Provider,
		TickerId: params.TickerId,
		Attempts: 0,
	}

	err = queueService.AddJobs([]jobs.JobAction{job})

	if err != nil {
		return clientError(http.StatusInternalServerError, err)
	} else {
		log.Printf("Added Job '%s' to queue", jobId)
	}

	// 4. enable the event timer
	err = eventsService.StartTickerScheduler()

	if err != nil {
		return clientError(http.StatusInternalServerError, err)
	}

	return clientSuccess(fmt.Sprintf("Success: Job '%s' queued", jobId))
}

func create(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return createStockIndex(request)
}
