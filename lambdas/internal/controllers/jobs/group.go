package jobs

import (
	"strings"

	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
)

func groupByProvider(tickers []ticker.EntityStub) map[provider.Name][]string {
	list := map[provider.Name][]string{}

	for _, item := range tickers {
		key := item.Provider
		tickerId, _ := strings.CutPrefix(item.TickerSort, "T#")

		list[key] = append(list[key], tickerId)
	}

	return list
}
