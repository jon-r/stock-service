package tickers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/providers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

func NewMock(cfg aws.Config, log logger.Logger) Controller {
	// todo redo this passing along the mock service
	providersService := providers.NewMock()
	dbRepository := db.NewRepository(cfg)

	return NewController(dbRepository, providersService, log)
}
