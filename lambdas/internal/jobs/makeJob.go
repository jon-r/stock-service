package jobs

import (
	"strings"

	"github.com/google/uuid"
	"jon-richards.com/stock-app/internal/providers"
)

func MakeJob(provider providers.ProviderName, tickerId string, jobType JobTypes) JobAction {
	return JobAction{
		JobId:    uuid.NewString(),
		Provider: provider,
		Type:     jobType,
		TickerId: tickerId,
		Attempts: 0,
	}
}

func MakeBulkJob(provider providers.ProviderName, tickerIds []string, jobType JobTypes) JobAction {
	tickerId := strings.Join(tickerIds, ",")

	return JobAction{
		JobId:    uuid.NewString(),
		Provider: provider,
		Type:     jobType,
		TickerId: tickerId,
		Attempts: 0,
	}
}
