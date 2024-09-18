package prices

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/providers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

func NewMock(cfg aws.Config, log logger.Logger, c clock.Clock, httpClient *http.Client) Controller {
	// todo redo this passing along the mock service
	providersService := providers.NewService(httpClient, c)
	dbRepository := db.NewRepository(cfg)

	return NewController(dbRepository, providersService, log)
}
