package job

import (
	"os"
	"strings"

	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
)

func NewJob(jobType Types, id string, provider provider.Name, tickerId string) *Job {
	return &Job{
		JobId:    id,
		Provider: provider,
		Type:     jobType,
		TickerId: tickerId,
		Attempts: 0,
	}
}

func NewJobs(jobTypes []Types, id string, provider provider.Name, tickerId string) *[]Job {
	jobs := make([]Job, len(jobTypes))
	for i, jobType := range jobTypes {
		jobs[i] = *NewJob(jobType, id, provider, tickerId)
	}
	return &jobs
}

func NewBulkJob(jobType Types, id string, provider provider.Name, tickerIds []string) *Job {
	return &Job{
		JobId:    id,
		Provider: provider,
		Type:     jobType,
		TickerId: strings.Join(tickerIds, ","),
		Attempts: 0,
	}
}

func QueueUrl() string {
	return os.Getenv("SQS_QUEUE_URL")
}
func DLQUrl() string {
	return os.Getenv("SQS_DLQ_URL")
}
