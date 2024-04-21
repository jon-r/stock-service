package jobs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"jon-richards.com/stock-app/internal/providers"
)

// https://dev.to/salesforceeng/mocks-in-go-tests-with-testify-mock-6pd

func mockUuid() string {
	return "test"
}

func TestMakeCreateJobs(t *testing.T) {
	tickerId := "EXAMPLE"
	actual := MakeCreateJobs(providers.PolygonIo, tickerId, mockUuid)

	expected := []JobAction{
		{JobId: "test", Provider: providers.PolygonIo, Type: LoadTickerDescription, TickerId: tickerId},
		{JobId: "test", Provider: providers.PolygonIo, Type: LoadHistoricalPrices, TickerId: tickerId},
	}

	assert.Equal(t, actual, &expected)
}
