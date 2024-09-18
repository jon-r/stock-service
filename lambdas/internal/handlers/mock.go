package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/prices"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

func NewMock(cfg aws.Config) *LambdaHandler {
	log := logger.NewMock()
	jobsCtrl := jobs.NewMock(cfg, log)
	// clock and httpClient only used for api calls, so can be nil here
	tickersCtrl := tickers.NewMock(cfg, log, nil, nil)
	pricesCtrl := prices.NewMock(cfg, log, nil, nil)

	return &LambdaHandler{tickersCtrl, jobsCtrl, pricesCtrl, log}
}

func NewMockWithHttpClient(cfg aws.Config, c clock.Clock, httpClient *http.Client) *LambdaHandler {
	log := logger.NewMock()
	jobsCtrl := jobs.NewMock(cfg, log)
	tickersCtrl := tickers.NewMock(cfg, log, c, httpClient)
	pricesCtrl := prices.NewMock(cfg, log, c, httpClient)

	return &LambdaHandler{tickersCtrl, jobsCtrl, pricesCtrl, log}
}
