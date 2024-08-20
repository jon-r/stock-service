package providers

import (
	"fmt"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/jon-r/stock-service/lambdas/internal/models/prices"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	apiStubber := test.NewApiStubber()
	mockClock := clock.NewMock()

	service := NewService(apiStubber.NewTestClient(), mockClock)

	date, _ := time.Parse(time.DateOnly, "2022-10-26")
	mockClock.Set(date)

	mockToday := mockClock.Now()
	mockYesterday := mockToday.Add(24 * -time.Hour)

	t.Run("GetDescription", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL:    "https://api.polygon.io/v3/reference/tickers/AAPL",
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getDescriptionRes.json"),
		})

		res, err := service.GetDescription(provider.PolygonIo, "AAPL")

		assert.Nil(t, err)
		assert.Equal(t, &ticker.Description{
			FullName:   "Apple Inc.",
			FullTicker: "XNAS:AAPL",
			Currency:   "usd",
			Icon:       "https://api.polygon.io/v1/reference/company-branding/d3d3LmFwcGxlLmNvbQ/images/2022-01-10_icon.png",
		}, res)
	})

	t.Run("GetHistoricalPrices", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/%v/%v?adjusted=true&sort=desc",
				startDate.UnixMilli(),
				mockToday.UnixMilli(),
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getHistoricalPricesRes.json"),
		})
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/%v/%v?adjusted=true&sort=desc?cursor=%v",
				startDate.UnixMilli(),
				mockToday.UnixMilli(),
				"bGltaXQ9MiZzb3J0PWFzYw",
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getHistoricalPricesRes2.json"),
		})

		res, err := service.GetHistoricalPrices(provider.PolygonIo, "AAPL")

		assert.Nil(t, err)
		assert.Equal(t, &[]prices.TickerPrices{
			{
				Id:        "AAPL",
				Open:      74.06,
				Close:     75.0875,
				High:      75.15,
				Low:       73.7975,
				Timestamp: models.Millis(time.UnixMilli(1577941200000)),
			},
			{
				Id:        "AAPL",
				Open:      74.2875,
				Close:     74.3575,
				High:      75.145,
				Low:       74.125,
				Timestamp: models.Millis(time.UnixMilli(1578027600000)),
			},
		}, res)
	})

	t.Run("GetDailyPrices", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			URL: fmt.Sprintf(
				"https://api.polygon.io/v2/aggs/grouped/locale/us/market/stocks/%v?adjusted=true",
				mockYesterday.Format(time.DateOnly),
			),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getDailyPricesRes.json"),
		})

		res, err := service.GetDailyPrices(provider.PolygonIo, []string{"AAPL", "META"})

		assert.Nil(t, err)
		assert.Equal(t, &[]prices.TickerPrices{{
			Id:        "AAPL",
			Open:      26.07,
			Close:     25.9102,
			High:      26.25,
			Low:       25.91,
			Timestamp: models.Millis(time.UnixMilli(1602705600000)),
		}, {
			Id:        "META",
			Open:      24.5,
			Close:     23.4,
			High:      24.763,
			Low:       22.65,
			Timestamp: models.Millis(time.UnixMilli(1602705600000)),
		}}, res)
	})

	err := apiStubber.VerifyAllStubsCalled()
	if err != nil {
		t.Error(err)
	}
}
