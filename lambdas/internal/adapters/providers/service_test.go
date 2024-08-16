package providers

import (
	"fmt"
	"testing"
	"time"

	"github.com/jon-r/stock-service/lambdas/internal/models/prices"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
	"github.com/jon-r/stock-service/lambdas/internal/models/ticker"
	"github.com/jon-r/stock-service/lambdas/internal/utils/test"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	apiStubber := test.NewApiStubber()

	service := NewService(apiStubber.NewTestClient())

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
			// todo thisll break right away, need to mock time :/
			URL:    fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/AAPL/range/1/day/1638316800000/%v?adjusted=true&sort=desc", time.Now().UnixMilli()),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getHistoricalPricesRes.json"),
		})

		res, err := service.GetHistoricalPrices(provider.PolygonIo, "AAPL")

		assert.Nil(t, err)
		assert.Equal(t, &ticker.Description{
			FullName:   "Apple Inc.",
			FullTicker: "XNAS:AAPL",
			Currency:   "usd",
			Icon:       "https://api.polygon.io/v1/reference/company-branding/d3d3LmFwcGxlLmNvbQ/images/2022-01-10_icon.png",
		}, res)
	})

	t.Run("GetDailyPrices", func(t *testing.T) {
		apiStubber.AddRequest(test.ReqStub{
			Method: "GET",
			// todo tidy this up
			URL:    fmt.Sprintf("https://api.polygon.io/v2/aggs/grouped/locale/us/market/stocks/%v?adjusted=true", time.Now().Add(time.Hour*-24).Format(time.DateOnly)),
			Input:  "",
			Output: test.ReadJsonToString("./testdata/getDailyPricesRes.json"),
		})

		res, err := service.GetDailyPrices(provider.PolygonIo, []string{"AAPL", "META"})

		assert.Nil(t, err)
		//_ := models.Millis.UnmarshalJSON()
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
