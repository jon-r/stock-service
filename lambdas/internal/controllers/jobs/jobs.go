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
	InvokeWorker(j job.Job) error
	RequeueJob(j job.Job, failReason string) error
	ReceiveJobs() (*[]job.Job, error)
	StopScheduledRule()
}

type jobsController struct {
	queueBroker     queue.Broker
	eventsScheduler events.Scheduler
	idGen           queue.NewIdFunc
	log             logger.Logger

	//attempts int
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

	return c.startScheduledRule()
}

func (c *jobsController) LaunchDailyTickerJobs(tickers *[]ticker.EntityStub) error {
	var err error
	var bulkJob *job.Job
	var newJobs []job.Job

	groupedTickerIds := ticker.GroupByProvider(*tickers)
	for p, tickerGroup := range groupedTickerIds {

		bulkJob = job.NewBulkJob(job.LoadDailyPrices, c.idGen(), p, tickerGroup)
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

	return c.startScheduledRule()
}

func (c *jobsController) startScheduledRule() error {
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

func (c *jobsController) InvokeWorker(j job.Job) error {
	var err error

	functionName := os.Getenv("LAMBDA_WORKER_NAME")
	_, err = c.eventsScheduler.InvokeFunction(functionName, j)

	if err != nil {
		c.log.Errorw("could not invoke function", "error", err)
		err = c.RequeueJob(j, err.Error())
	}

	_, err = c.queueBroker.DeleteMessage(job.QueueUrl(), j.JobId)
	if err != nil {
		c.log.Errorw("could not delete message", "error", err)
	}

	return err
}

func (c *jobsController) StopScheduledRule() {
	var err error

	ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")
	_, err = c.eventsScheduler.EnableRule(ruleName)

	if err != nil {
		c.log.Errorw("error disabling rule", "error", err)
	}
}

func (c *jobsController) RequeueJob(j job.Job, failReason string) error {
	var err error

	if j.Attempts > 2 {
		failedJob := job.NewFailedJob(j, failReason)
		_, err = c.queueBroker.SendMessage(job.DLQUrl(), failedJob)
	} else {
		// put the failed item back into the queue
		j.Attempts += 1
		_, err = c.queueBroker.SendMessage(job.QueueUrl(), j)
	}

	return err
}

func (c *jobsController) ReceiveJobs() (*[]job.Job, error) {
	var err error

	c.log.Debug("attempting to receive jobs")
	messages, err := c.queueBroker.ReceiveMessages(job.QueueUrl())

	if err != nil {
		c.log.Errorw("error receiving messages", "error", err)
		return nil, err
	}

	return job.NewJobsFromSqs(messages)
}

func NewController(queueBroker queue.Broker, eventsScheduler events.Scheduler, idGen queue.NewIdFunc, log logger.Logger) Controller {
	return &jobsController{
		queueBroker,
		eventsScheduler,
		idGen,
		log,
		//0,
	}
}
