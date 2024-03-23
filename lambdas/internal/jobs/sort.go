package jobs

import "jon-richards.com/stock-app/internal/providers"

func groupByProvider(tickers []providers.TickerItemStub) map[providers.ProviderName][]string {
	list := map[providers.ProviderName][]string{}

	for _, item := range tickers {
		key := item.Provider

		list[key] = append(list[key], item.TickerId)
	}

	return list
}
