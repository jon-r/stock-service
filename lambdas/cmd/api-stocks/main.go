package main

import (
	"context"
	"fmt"
	"net/http"

	awsEvents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/utils/response"
)

type handler struct{ *handlers.LambdaHandler }

var apiStocksHandler = handler{handlers.NewLambdaHandler()}

func (h *handler) HandleRequest(ctx context.Context, request awsEvents.APIGatewayProxyRequest) (*awsEvents.APIGatewayProxyResponse, error) {
	h.Log.LoadContext(ctx)
	defer h.Log.Sync()

	switch request.HTTPMethod {
	case http.MethodPost:
		return h.handlePost(request)
	default:
		err := fmt.Errorf("request method %s not supported", request.HTTPMethod)
		h.Log.Errorw("Request error",
			"status", http.StatusMethodNotAllowed,
			"message", err,
		)

		return response.StatusMethodNotAllowed(err)
	}
}

func main() {
	lambda.Start(apiStocksHandler.HandleRequest)
}
