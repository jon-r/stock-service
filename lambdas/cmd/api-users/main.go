package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	Message string `json:"greeting"`
	Event   events.APIGatewayProxyRequest
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := &response{
		Message: "hello USERS world!",
		Event:   request,
	}
	body, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Error parsing payload", StatusCode: 400}, err
	}
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
