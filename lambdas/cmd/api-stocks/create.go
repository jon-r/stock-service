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
	params, err := ticker.ParamsFromJsonString(req.Body)

	if err != nil {
		h.log.Errorw("error unmarshalling ticker", "error", err)
		return response.StatusBadRequest(err)
	}

	// 2. enter basic content to the database
	//newTicker := ticker.NewTickerEntity(params)
	err = h.tickers.New(params)

	if err != nil {
		return response.StatusServerError(err)
	}

	// 3. Create new job queue items
	err = h.jobs.LaunchNewTickerJobs(params.Provider, params.TickerId)

	//jobs := []job.Types{job.LoadTickerDescription, job.LoadHistoricalPrices}
	//newTickerJobs := job.NewJobs(jobs, h.idGen(), params.Provider, params.TickerId)
	//
	//h.log.Debugw("add jobs to the queue", "jobs", newTickerJobs)
	//_, err = h.queueBroker.SendMessages(job.QueueUrl(), newTickerJobs)
	//
	if err != nil {
		return response.StatusServerError(err)
	}
	//
	//// 4. enable the jobs ticker
	//ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")
	//_, err = h.eventsScheduler.EnableRule(ruleName)
	//
	//if err != nil {
	//	h.log.Errorw("error enabling rule", "error", err)
	//	return response.StatusServerError(err)
	//}
	//
	//// 5. manually trigger the lambda
	//functionName := os.Getenv("LAMBDA_TICKER_NAME")
	//_, err = h.eventsScheduler.InvokeFunction(functionName, nil)
	//
	//if err != nil {
	//	h.log.Errorw("error invoking function but continuing anyway", "error", err)
	//}

	return response.StatusOK(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId))
}

func (h *handler) handlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return h.createTicker(req)
}
