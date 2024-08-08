package jobs_old

import (
	"testing"

	"github.com/jon-r/stock-service/lambdas/internal/providers_old"
	"github.com/stretchr/testify/assert"
)

// https://dev.to/salesforceeng/mocks-in-go-tests-with-testify-mock-6pd
// todo look at coverage library?
//   or pipeline
//   https://medium.com/synechron/how-to-set-up-a-test-coverage-threshold-in-go-and-github-167f69b940dc

func mockUuid() string {
	return "test"
}

func TestMakeCreateJobs(t *testing.T) {
	tickerId := "EXAMPLE"
	actual := MakeCreateJobs(providers_old.PolygonIo, tickerId, mockUuid)

	expected := []JobAction{
		{JobId: "test", Provider: providers_old.PolygonIo, Type: LoadTickerDescription, TickerId: tickerId},
		{JobId: "test", Provider: providers_old.PolygonIo, Type: LoadHistoricalPrices, TickerId: tickerId},
	}

	assert.Equal(t, actual, &expected)
}
