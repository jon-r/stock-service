package jobs

import (
	"strings"

	"github.com/jon-r/stock-service/lambdas/internal/providers"
)

type UuidGen func() string

func MakeJob(provider providers.ProviderName, tickerId string, jobType JobTypes, newUuid UuidGen) JobAction {
	return JobAction{
		JobId:    newUuid(),
		Provider: provider,
		Type:     jobType,
		TickerId: tickerId,
		Attempts: 0,
	}
}

func MakeBulkJob(provider providers.ProviderName, tickerIds []string, jobType JobTypes, newUuid UuidGen) JobAction {
	tickerId := strings.Join(tickerIds, ",")

	return JobAction{
		JobId:    newUuid(),
		Provider: provider,
		Type:     jobType,
		TickerId: tickerId,
		Attempts: 0,
	}
}

func MakeCreateJobs(provider providers.ProviderName, tickerId string, newUuid UuidGen) *[]JobAction {
	newItemActions := []JobTypes{
		LoadTickerDescription,
		LoadHistoricalPrices,
		// TODO jobs.LoadHistoricalDividends,
		// TODO jobs.LoadTickerIcon,
	}

	jobActions := make([]JobAction, len(newItemActions))
	for i, jobType := range newItemActions {
		job := MakeJob(provider, tickerId, jobType, newUuid)

		jobActions[i] = job
	}

	return &jobActions
}

func MakeUpdateJobs(tickers []providers.TickerItemStub, newUuid UuidGen) *[]JobAction {

	//tickerLimit := 10
	groupedTickerIds := groupByProvider(tickers)

	var job JobAction
	var jobActions []JobAction
	for provider, tickerGroup := range groupedTickerIds {

		job = MakeBulkJob(provider, tickerGroup, UpdatePrices, newUuid)
		jobActions = append(jobActions, job)

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

	return &jobActions
}
