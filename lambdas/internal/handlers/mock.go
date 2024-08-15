package handlers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/providers"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/prices"
	"github.com/jon-r/stock-service/lambdas/internal/controllers/tickers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
	"go.uber.org/zap/zapcore"
)

func NewMock(cfg aws.Config) *LambdaHandler {
	idGen := func() string { return "TEST_ID" }
	log := logger.NewLogger(zapcore.DPanicLevel) // todo raise once finished

	// todo once tests split up, some of this can be moved to the controller
	providersService := providers.NewMock()
	queueBroker := queue.NewBroker(cfg, idGen)
	eventsScheduler := events.NewScheduler(cfg)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, eventsScheduler, idGen, log)
	tickersCtrl := tickers.NewController(dbRepository, providersService, log)
	pricesCtrl := prices.NewController(dbRepository, providersService, log)

	return &LambdaHandler{tickersCtrl, jobsCtrl, pricesCtrl, log}
}
