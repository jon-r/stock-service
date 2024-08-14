package handlers

import (
	"github.com/google/uuid"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/config"
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

type LambdaHandler struct {
	Tickers tickers.Controller
	Jobs    jobs.Controller
	Prices  prices.Controller
	Log     logger.Logger
}

func NewLambdaHandler() *LambdaHandler {
	cfg := config.GetAwsConfig()
	log := logger.NewLogger(zapcore.InfoLevel)
	idGen := uuid.NewString

	// todo once tests split up, some of this can be moved to the controllers
	eventsScheduler := events.NewScheduler(cfg)
	providersService := providers.NewService(nil)
	queueBroker := queue.NewBroker(cfg, idGen)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, eventsScheduler, idGen, log)
	tickersCtrl := tickers.NewController(dbRepository, providersService, log)
	pricesCtrl := prices.NewController(dbRepository, providersService, log)

	return &LambdaHandler{tickersCtrl, jobsCtrl, pricesCtrl, log}
}
