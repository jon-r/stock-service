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

func MakeCreateJobs(provider providers.ProviderName, tickerId string) *[]JobAction {
	newItemActions := []JobTypes{
		LoadTickerDescription,
		LoadHistoricalPrices,
		// TODO jobs.LoadHistoricalDividends,
		// TODO jobs.LoadTickerIcon,
	}

	jobActions := make([]JobAction, len(newItemActions))
	for i, jobType := range newItemActions {
		job := MakeJob(provider, tickerId, jobType)

		jobActions[i] = job
	}

	return &jobActions
}

func MakeUpdateJobs(tickers []providers.TickerItemStub) *[]JobAction {
	updateItemActions := []JobTypes{
		UpdatePrices,
		// TODO jobs.UpdateDividends,
	}

	groupedTickerIds := groupByProvider(tickers)

	jobCount := len(updateItemActions) * ((len(tickers) / 20) + 1)
	jobActions := make([]JobAction, jobCount)
	for provider, tickerGroup := range groupedTickerIds {
		chunkedTickers := chunkIds(tickerGroup, 20)

		for _, chunk := range chunkedTickers {

			for _, jobType := range updateItemActions {
				job := MakeBulkJob(provider, chunk, jobType)

				jobActions = append(jobActions, job)
			}
		}
	}

	return &jobActions
}
