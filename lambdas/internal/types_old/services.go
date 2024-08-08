package types_old

import (
	"github.com/jon-r/stock-service/lambdas/internal/db_old"
	"github.com/jon-r/stock-service/lambdas/internal/jobs_old"
	"github.com/jon-r/stock-service/lambdas/internal/scheduler_old"
	"go.uber.org/zap"
)

type ServiceHandler struct {
	QueueService  *jobs_old.QueueRepository
	EventsService *scheduler_old.EventsRepository
	DbService     *db_old.DatabaseRepository
	LogService    *zap.SugaredLogger
	NewUuid       jobs_old.UuidGen
}
