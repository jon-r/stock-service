package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jon-r/stock-service/lambdas/internal/jobs_old"
	"github.com/jon-r/stock-service/lambdas/internal/providers_old"
)

func (handler ApiStockHandler) createTicker(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error

	// 1. get ticket and provider from post request body
	var params providers_old.NewTickerParams
	err = json.Unmarshal([]byte(request.Body), &params)

	if err != nil {
		handler.LogService.Errorw("Request error",
			"status", http.StatusBadRequest,
			"message", err,
		)
		return clientError(http.StatusBadRequest, err)
	}

	// 2. enter basic content to the database
	err = handler.DbService.NewTickerItem(handler.LogService, params)

	if err != nil {
		handler.LogService.Errorw("Request error",
			"status", http.StatusInternalServerError,
			"message", err,
		)
		return clientError(http.StatusInternalServerError, err)
	}

	// 3. Create new job queue item
	newItemJobs := jobs_old.MakeCreateJobs(params.Provider, params.TickerId, handler.NewUuid)

	handler.LogService.Infow("Add jobs_old to the queue",
		"jobs_old", *newItemJobs,
	)
	err = handler.QueueService.AddJobs(*newItemJobs, handler.NewUuid)

	if err != nil {
		handler.LogService.Errorw("Request error",
			"status", http.StatusInternalServerError,
			"message", err,
		)
		return clientError(http.StatusInternalServerError, err)
	}

	// 4. enable the jobs_old ticker
	err = handler.EventsService.StartTickerScheduler()

	if err != nil {
		handler.LogService.Errorw("Request error",
			"status", http.StatusInternalServerError,
			"message", err,
		)
		return clientError(http.StatusInternalServerError, err)
	}

	return clientSuccess(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId)), nil
}

func (handler ApiStockHandler) create(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return handler.createTicker(request)
}
