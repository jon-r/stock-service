package jobs

import (
	"os"

	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

type Controller interface {
	LaunchNewTickerJobs(ticker *ticker.NewTickerParams) error
	LaunchDailyTickerJobs(tickers *[]ticker.EntityStub) error
}

type jobsController struct {
	queueBroker     queue.Broker
	eventsScheduler events.Scheduler
	idGen           queue.NewIdFunc
	log             logger.Logger
}

func (c *jobsController) LaunchNewTickerJobs(newTicker *ticker.NewTickerParams) error {
	var err error

	newJobs := []job.Job{
		*job.NewJob(job.LoadTickerDescription, c.idGen(), newTicker.Provider, newTicker.TickerId),
		*job.NewJob(job.LoadHistoricalPrices, c.idGen(), newTicker.Provider, newTicker.TickerId),
	}

	c.log.Debugw("add jobs to the queue", "jobs", newJobs)
	_, err = c.queueBroker.SendMessages(job.QueueUrl(), newJobs)

	if err != nil {
		c.log.Errorw("error sending messages", "error", err)
		return err
	}

	return c.startTicker()
}

func (c *jobsController) startTicker() error {
	var err error

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

func (c *jobsController) LaunchDailyTickerJobs(tickers *[]ticker.EntityStub) error {
	var err error
	var bulkJob *job.Job
	var newJobs []job.Job

	groupedTickerIds := groupByProvider(*tickers)
	for p, tickerGroup := range groupedTickerIds {

		bulkJob = job.NewBulkJob(job.UpdatePrices, c.idGen(), p, tickerGroup)
		newJobs = append(newJobs, *bulkJob)

		// todo STK-90 no need to chunk for prices, just dividends
		//chunkedTickers := lo.Chunk(tickerGroup, tickerLimit)
		// have a look at AddTickerPrices for how to chunk in a way that dynamoDB likes
		//for _, chunk := range chunkedTickers {
		//
		//	for _, jobType := range updateItemActions {
		//		job := MakeBulkJob(provider, chunk, jobType)
		//
		//		jobActions = append(jobActions, job)
		//	}
		//}
	}

	c.log.Debugw("add jobs to the queue", "jobs", newJobs)
	_, err = c.queueBroker.SendMessages(job.QueueUrl(), newJobs)

	if err != nil {
		c.log.Errorw("error sending messages", "error", err)
		return err
	}

	return c.startTicker()
}

func NewController(queueBroker queue.Broker, eventsScheduler events.Scheduler, idGen queue.NewIdFunc, log logger.Logger) Controller {
	return &jobsController{
		queueBroker,
		eventsScheduler,
		idGen,
		log,
	}
}
