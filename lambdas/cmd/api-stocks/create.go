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
		h.Log.Errorw("error unmarshalling ticker", "error", err)
		return response.StatusBadRequest(err)
	}

	// 2. enter basic content to the database
	err = h.Tickers.New(params)

	if err != nil {
		return response.StatusServerError(err)
	}

	// 3. Create new job queue items
	err = h.Jobs.LaunchNewTickerJobs(params)

	if err != nil {
		return response.StatusServerError(err)
	}

	return response.StatusOK(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId))
}

func (h *handler) handlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	// todo add a 'toggle event ticker' endpoint
	return h.createTicker(req)
}
