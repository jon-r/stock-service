package main

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func TestPollUntilCanceled(t *testing.T) {
	t.Run("triggers function until canceled", handleRequestsUntilStop)
}

func handleRequestsUntilStop(t *testing.T) {
	mockClock := clock.NewMock()
	mockHandler := handler{Clock: mockClock}

	functionSpyCount := 0
	functionSpy := func() { functionSpyCount++ }

	ctx, cancel := context.WithCancel(context.TODO())

	go mockHandler.pollUntilCancelled(ctx, functionSpy, 2*time.Second)

	// start at 0
	assert.Equal(t, 0, functionSpyCount)

	for range [9]int{} {
		mockClock.Add(time.Second)
	}

	// 10 invocations
	assert.Equal(t, 4, functionSpyCount)

	// stop the poll
	cancel()

	for range [9]int{} {
		mockClock.Add(time.Second)
	}

	// no more invocations
	assert.Equal(t, 4, functionSpyCount)
}
