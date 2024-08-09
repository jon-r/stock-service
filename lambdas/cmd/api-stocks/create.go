package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/response"
)

func (h *handler) createTicker(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error

	// 1. get ticker and provider from post request body
	var params ticker.NewTickerParams
	err = json.Unmarshal([]byte(req.Body), &params)

	if err != nil {
		h.log.Errorw("error unmarshalling ticker", "error", err)
		return response.StatusBadRequest(err)
	}

	// 2. enter basic content to the database
	newTicker := ticker.NewTickerEntity(params)
	h.log.Debugw("new ticker", "params", params, "item", newTicker)
	_, err = h.dbRepository.AddOne(ticker.TableName(), newTicker)

	if err != nil {
		h.log.Errorw("error adding ticker", "error", err)
		return response.StatusServerError(err)
	}

	// 3. Create new job queue items
	newTickerJobs := []job.Job{
		*job.NewJob(job.LoadTickerDescription, h.idGen(), params.Provider, params.TickerId),
		*job.NewJob(job.LoadHistoricalPrices, h.idGen(), params.Provider, params.TickerId),
	}

	h.log.Debugw("add jobs to the queue",
		"jobs", newTickerJobs,
	)
	_, err = h.queueBroker.SendMessages(job.QueueUrl(), newTickerJobs)

	if err != nil {
		h.log.Errorw("error sending messages", "error", err)
		return response.StatusServerError(err)
	}

	// 4. enable the jobs ticker
	ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")
	_, err = h.eventsScheduler.EnableRule(ruleName)

	if err != nil {
		h.log.Errorw("error enabling rule", "error", err)
		return response.StatusServerError(err)
	}

	// 5. manually trigger the lambda
	functionName := os.Getenv("LAMBDA_TICKER_NAME")
	_, err = h.eventsScheduler.InvokeFunction(functionName, nil)

	if err != nil {
		h.log.Errorw("error invoking function but continuing anyway", "error", err)
	}

	return response.StatusOK(fmt.Sprintf("Success: ticker '%s' queued", params.TickerId))
}

func (h *handler) handlePost(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo would be switch if multiple endpoints
	return h.createTicker(req)
}
