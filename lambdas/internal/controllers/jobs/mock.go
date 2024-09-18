package jobs

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/events"
	"github.com/jon-r/stock-service/lambdas/internal/adapters/queue"
	"github.com/jon-r/stock-service/lambdas/internal/utils/logger"
)

func NewMock(cfg aws.Config, log logger.Logger) Controller {
	idGen := func() string { return "TEST_ID" }

	queueBroker := queue.NewBroker(cfg, idGen)
	eventsScheduler := events.NewScheduler(cfg)

	return NewController(queueBroker, eventsScheduler, idGen, log)
}
