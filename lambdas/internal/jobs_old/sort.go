package jobs_old

import (
	"strings"

	"github.com/jon-r/stock-service/lambdas/internal/providers_old"
)

func groupByProvider(tickers []providers_old.TickerItemStub) map[providers_old.ProviderName][]string {
	list := map[providers_old.ProviderName][]string{}

	for _, item := range tickers {
		key := item.Provider
		tickerId, _ := strings.CutPrefix(item.TickerSort, "T#")

		list[key] = append(list[key], tickerId)
	}

	return list
}
