package prices

import (
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/providers"
	"github.com/jon-r/stock-service/lambdas/internal/models/prices"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

type Controller interface {
	LoadHistoricalPrices(p provider.Name, id string) error
	LoadDailyPrices(p provider.Name, ids []string) error
}

type pricesController struct {
	providers providers.Service
	db        db.Repository
	log       logger.Logger
}

func (c *pricesController) LoadHistoricalPrices(provider provider.Name, id string) error {
	var err error

	historicalPrices, err := c.providers.GetHistoricalPrices(provider, id)

	if err != nil {
		c.log.Errorw("failed to get historical prices", "provider", provider, "error", err)
		return err
	}

	return c.addPricesToDb(historicalPrices)
}

func (c *pricesController) LoadDailyPrices(provider provider.Name, ids []string) error {
	var err error

	historicalPrices, err := c.providers.GetDailyPrices(provider, ids)

	if err != nil {
		c.log.Errorw("failed to get historical prices", "provider", provider, "error", err)
		return err
	}

	return c.addPricesToDb(historicalPrices)
}

func (c *pricesController) addPricesToDb(pricesList *[]prices.TickerPrices) error {
	pricesEntities := prices.MapPriceEntities(pricesList)

	_, err := c.db.AddMany(prices.TableName(), *pricesEntities)

	if err != nil {
		c.log.Errorw("could not add prices to database", "error", err)
	}

	return err
}

func NewController(db db.Repository, providers providers.Service, log logger.Logger) Controller {
	return &pricesController{providers, db, log}
}
