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

func chunkIds(tickers []string, size int) [][]string {
	numberOfChunks := len(tickers) / size

	if len(tickers)%size != 0 {
		numberOfChunks += 1
	}

	result := make([][]string, 0, numberOfChunks)

	for i := 0; i < numberOfChunks; i++ {
		last := (i + 1) * size
		if last > len(tickers) {
			last = len(tickers)
		}
		result = append(result, tickers[i*size:last])
	}

	return result
}
