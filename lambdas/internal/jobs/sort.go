package jobs

import (
	"strings"

	"jon-richards.com/stock-app/internal/providers"
)

func groupByProvider(tickers []providers.TickerItemStub) map[providers.ProviderName][]string {
	list := map[providers.ProviderName][]string{}

	for _, item := range tickers {
		key := item.Provider
		tickerId, _ := strings.CutPrefix(item.TickerSort, "T#")

		list[key] = append(list[key], tickerId)
	}

	return list
}
