package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"jon-richards.com/stock-app/internal/db"
	"jon-richards.com/stock-app/internal/jobs"
)

type ResponseBody struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

var queueService = jobs.NewQueueService()
var eventsService = jobs.NewEventsService()
var dbService = db.NewDatabaseService()

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

	body, _ := json.Marshal(ResponseBody{
		Message: http.StatusText(status),
		Status:  status,
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
	}, nil
}

func clientSuccess(message string) (*events.APIGatewayProxyResponse, error) {
	if message == "" {
		message = "Success"
	}

	body, _ := json.Marshal(ResponseBody{
		Message: message,
		Status:  http.StatusOK,
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
