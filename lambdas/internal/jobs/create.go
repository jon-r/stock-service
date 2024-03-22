package jobs

import "jon-richards.com/stock-app/internal/providers"

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
