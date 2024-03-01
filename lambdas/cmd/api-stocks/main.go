package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// todo return different error statuses and look at the error
	switch request.HTTPMethod {
	case "POST":
		return create(request)
	default:
		return clientError(http.StatusMethodNotAllowed, nil)
	}
}

func clientError(status int, err error) (*events.APIGatewayProxyResponse, error) {
	// todo more detailed error handling?
	if err != nil {
		fmt.Printf("request error: %v", err)
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func clientSuccess() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
