package main

import (
	"context"
	"testing"

	"github.com/jon-r/stock-service/lambdas/internal/testutil"
)

func TestHandleRequest(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) { handleRequestNoErrors(t) })
}

func handleRequestNoErrors(t *testing.T) {
	stubber, mockServiceHandler := testutil.EnterTest(nil)
	mockHandler := DataManagerHandler{*mockServiceHandler}

	err := mockHandler.updateAllTickers(context.TODO())

	testutil.Assert(stubber, err, nil, t)
}
