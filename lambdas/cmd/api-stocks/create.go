package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/response"
)

func (h *handler) createTicker(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error

	// 1. get ticker and provider from post request body
	params, err := ticker.NewParamsFromJsonString(req.Body)

	if err != nil {
		h.log.Errorw("error unmarshalling ticker", "error", err)
		return response.StatusBadRequest(err)
	}

	// 2. enter basic content to the database
	err = h.tickers.New(params)

	if err != nil {
		return response.StatusServerError(err)
	}

	// 3. Create new job queue items
	err = h.jobs.LaunchNewTickerJobs(params)

	if err != nil {
		return response.StatusServerError(err)
	}

	return response.StatusOK(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId))
}

func (h *handler) handlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return h.createTicker(req)
}
