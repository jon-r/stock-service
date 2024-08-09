package job

import (
	"os"
	"strings"

	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

func NewJob(id string, provider provider.Name, tickerId string, jobType Types) *Job {
	return &Job{
		JobId:    id,
		Provider: provider,
		Type:     jobType,
		TickerId: tickerId,
		Attempts: 0,
	}
}

func NewBulkJob(id string, provider provider.Name, tickerIds []string, jobType Types) *Job {
	return &Job{
		JobId:    id,
		Provider: provider,
		Type:     jobType,
		TickerId: strings.Join(tickerIds, ","),
		Attempts: 0,
	}
}

func (j *Job) GetQueueUrl() string {
	return os.Getenv("SQS_QUEUE_URL")
}
func (j *Job) GetDLQUrl() string {
	return os.Getenv("SQS_DLQ_URL")
}
