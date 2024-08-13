package jobs

import (
	"fmt"
	"os"

	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

type Controller interface {
	LaunchNewTickerJobs(provider provider.Name, tickerId string) error
	LaunchDailyTickerJobs(provider provider.Name, tickerIds []string) error
}

type jobsController struct {
	queueBroker     queue.Broker
	eventsScheduler events.Scheduler
	idGen           queue.NewIdFunc
	log             logger.Logger
}

func (c *jobsController) LaunchNewTickerJobs(provider provider.Name, tickerId string) error {
	var err error

	newTickerJobs := []job.Job{
		*job.NewJob(job.LoadTickerDescription, c.idGen(), provider, tickerId),
		*job.NewJob(job.LoadHistoricalPrices, c.idGen(), provider, tickerId),
	}

	c.log.Debugw("add jobs to the queue", "jobs", newTickerJobs)
	_, err = c.queueBroker.SendMessages(job.QueueUrl(), newTickerJobs)

	if err != nil {
		c.log.Errorw("error sending messages", "error", err)
		return err
	}

	ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")
	_, err = c.eventsScheduler.EnableRule(ruleName)

	if err != nil {
		c.log.Errorw("error enabling rule", "error", err)
		return err
	}

	// 5. manually trigger the lambda
	functionName := os.Getenv("LAMBDA_TICKER_NAME")
	_, err = c.eventsScheduler.InvokeFunction(functionName, nil)

	if err != nil {
		c.log.Warnw("could not invoke function but continuing anyway (it may already be running)", "error", err)
	}

	return nil
}

func (c *jobsController) LaunchDailyTickerJobs(provider provider.Name, tickerIds []string) error {
	return fmt.Errorf("NOT IMPLEMENTED")
}

func NewController(queueBroker queue.Broker, eventsScheduler events.Scheduler, idGen queue.NewIdFunc, log logger.Logger) Controller {
	return &jobsController{
		queueBroker,
		eventsScheduler,
		idGen,
		log,
	}
}
