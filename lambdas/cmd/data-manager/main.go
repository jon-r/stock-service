package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
)

type handler struct{ *handlers.LambdaHandler }

var dataManagerHandler = handler{handlers.NewLambdaHandler()}

func (h *handler) HandleRequest(ctx context.Context) error {
	// todo look at zap docs to see if this can be done better. its not passing context to controllers
	h.Log = h.Log.LoadLambdaContext(ctx)
	defer h.Log.Sync()

	// 1. get all tickers
	tickerList, err := h.Tickers.GetAll()

	if err != nil {
		return err
	}

	// 2. convert the jobs into update actions
	return h.Jobs.LaunchDailyTickerJobs(tickerList)
}

func main() {
	lambda.Start(dataManagerHandler.HandleRequest)
}
