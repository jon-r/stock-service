package tickers

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/providers"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

type Controller interface {
	New(params *ticker.NewTickerParams) error
	LoadDescription(provider provider.Name, tickerId string) error
	GetOne(tickerId string) (*ticker.Entity, error)
	GetAll() (*[]ticker.EntityStub, error)
}

type tickersController struct {
	db        db.Repository
	providers providers.Service
	Log       logger.Logger
}

func (c *tickersController) New(params *ticker.NewTickerParams) error {
	entity := ticker.NewTickerEntity(params)

	c.Log.Debugw("new ticker", "entity", entity, "item", params)

	_, err := c.db.AddOne(ticker.TableName(), entity)

	if err != nil {
		c.Log.Errorw("error adding ticker", "error", err)
	}

	return err
}

func (c *tickersController) LoadDescription(provider provider.Name, tickerId string) error {
	var err error

	description, err := c.providers.GetDescription(provider, tickerId)

	if err != nil {
		c.Log.Errorw("error loading description", "error", err)
		return err
	}

	updateEx := expression.Set(expression.Name("Description"), expression.Value(*description))
	update, err := expression.NewBuilder().WithUpdate(updateEx).Build()

	if err != nil {
		c.Log.Errorw("error building update expression", "error", err)
		return err
	}

	item := ticker.NewTickerEntity(&ticker.NewTickerParams{Provider: provider, TickerId: tickerId})

	c.Log.Debugw("Update item", "item", item, "key", item.GetKey())

	_, err = c.db.Update(ticker.TableName(), item.GetKey(), update)

	if err != nil {
		c.Log.Errorw("error updating ticker", "error", err)
	}

	return err
}

func (c *tickersController) GetOne(tickerId string) (*ticker.Entity, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (c *tickersController) GetAll() (*[]ticker.EntityStub, error) {
	var err error

	filterEx := expression.Name("SK").BeginsWith(string(ticker.KeyTickerId))
	projEx := expression.NamesList(
		expression.Name("SK"), expression.Name("Provider"),
	)
	query, err := expression.NewBuilder().WithFilter(filterEx).WithProjection(projEx).Build()

	if err != nil {
		c.Log.Errorw("error building query", "error", err)
		return nil, err
	}

	entities, err := c.db.GetMany(ticker.TableName(), query)

	if err != nil {
		c.Log.Errorw("error getting tickers", "error", err)
		return nil, err
	}

	return ticker.NewStubsFromDynamoDb(entities)
}

func NewController(db db.Repository, providers providers.Service, log logger.Logger) Controller {
	return &tickersController{db, providers, log}
}
