package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/providers"
)

func (handler ApiStockHandler) createTicker(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error

	// 1. get ticket and provider from post request body
	var params providers.NewTickerParams
	err = json.Unmarshal([]byte(request.Body), &params)

	if err != nil {
		handler.logService.Errorw("Request error",
			"status", http.StatusBadRequest,
			"message", err,
		)
		return clientError(http.StatusBadRequest, err)
	}

	// 2. enter basic content to the database
	err = handler.dbService.NewTickerItem(handler.logService, params)

	if err != nil {
		handler.logService.Errorw("Request error",
			"status", http.StatusInternalServerError,
			"message", err,
		)
		return clientError(http.StatusInternalServerError, err)
	}

	// 3. Create new job queue item
	newItemJobs := jobs.MakeCreateJobs(params.Provider, params.TickerId, handler.newUuid)

	handler.logService.Infow("Add jobs to the queue",
		"jobs", *newItemJobs,
	)
	err = handler.queueService.AddJobs(*newItemJobs, handler.newUuid)

	if err != nil {
		handler.logService.Errorw("Request error",
			"status", http.StatusInternalServerError,
			"message", err,
		)
		return clientError(http.StatusInternalServerError, err)
	}

	// 4. enable the jobs ticker
	err = handler.eventsService.StartTickerScheduler()

	if err != nil {
		handler.logService.Errorw("Request error",
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
