package tickers

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/db"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

type Controller interface {
	New(params *ticker.NewTickerParams) error
	SetDescription(tickerId string, description *ticker.Description) (*dynamodb.UpdateItemOutput, error)
	GetOne(tickerId string) (*ticker.Entity, error)
	GetAll() (*[]ticker.EntityStub, error)
}

type tickersController struct {
	db  db.Repository
	log logger.Logger
}

func (c *tickersController) New(params *ticker.NewTickerParams) error {
	entity := ticker.NewTickerEntity(params)

	c.log.Debugw("new ticker", "entity", entity, "item", params)

	_, err := c.db.AddOne(ticker.TableName(), entity)

	if err != nil {
		c.log.Errorw("error adding ticker", "error", err)
	}

	return err
}

func (c *tickersController) SetDescription(tickerId string, description *ticker.Description) (*dynamodb.UpdateItemOutput, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
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
		c.log.Errorw("error building query", "error", err)
		return nil, err
	}

	entities, err := c.db.GetMany(ticker.TableName(), query)

	if err != nil {
		c.log.Errorw("error getting tickers", "error", err)
		return nil, err
	}

	return ticker.NewStubsFromDynamoDb(entities)
}

func NewController(db db.Repository, log logger.Logger) Controller {
	return &tickersController{db, log}
}
