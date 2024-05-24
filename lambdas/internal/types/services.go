package types

import (
	"github.com/jon-r/stock-service/lambdas/internal/db"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler"
	"go.uber.org/zap"
)

type ServiceHandler struct {
	QueueService  *jobs.QueueRepository
	EventsService *scheduler.EventsRepository
	DbService     *db.DatabaseRepository
	LogService    *zap.SugaredLogger
	NewUuid       jobs.UuidGen
}
