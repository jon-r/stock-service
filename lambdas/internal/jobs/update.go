package jobs

import "jon-richards.com/stock-app/internal/providers"

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
