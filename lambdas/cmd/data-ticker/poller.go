package main

import (
	"context"
	"time"
)

// TODO TEST -> trigger stop
func (h *handler) pollUntilCancelled(ctx context.Context, fn func(), interval time.Duration) {
	ticker := h.Clock.Ticker(interval)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			fn()
		}
	}
}
