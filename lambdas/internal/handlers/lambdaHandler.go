package handlers

import (
	"github.com/benbjohnson/clock"
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
	l := logger.New(zapcore.InfoLevel)
	idGen := uuid.NewString

	eventsScheduler := events.NewScheduler(cfg)
	providersService := providers.NewService(nil, clock.New())
	queueBroker := queue.NewBroker(cfg, idGen)
	dbRepository := db.NewRepository(cfg)

	jobsCtrl := jobs.NewController(queueBroker, eventsScheduler, idGen, l)
	tickersCtrl := tickers.NewController(dbRepository, providersService, l)
	pricesCtrl := prices.NewController(dbRepository, providersService, l)

	return &LambdaHandler{tickersCtrl, jobsCtrl, pricesCtrl, l}
}
