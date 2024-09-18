package handlers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/prices"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

func NewMock(cfg aws.Config) *LambdaHandler {
	log := logger.NewMock()
	jobsCtrl := jobs.NewMock(cfg, log)
	tickersCtrl := tickers.NewMock(cfg, log)
	pricesCtrl := prices.NewMock(cfg, log)

	return &LambdaHandler{tickersCtrl, jobsCtrl, pricesCtrl, log}
}
