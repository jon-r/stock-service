package response

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Body struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func response(status int, message string) *events.APIGatewayProxyResponse {
	body, _ := json.Marshal(Body{
		Message: message,
		Status:  http.StatusText(status),
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(body),
	}
}

func StatusOK(message string) (*events.APIGatewayProxyResponse, error) {
	if message == "" {
		message = http.StatusText(http.StatusOK)
	}

	return response(http.StatusOK, message), nil
}

func StatusMethodNotAllowed(err error) (*events.APIGatewayProxyResponse, error) {
	return response(http.StatusMethodNotAllowed, "Method not allowed"), err
}

func StatusBadRequest(err error) (*events.APIGatewayProxyResponse, error) {
	return response(http.StatusBadRequest, "Bad request"), err
}

func StatusServerError(err error) (*events.APIGatewayProxyResponse, error) {
	return response(http.StatusInternalServerError, "Internal server error"), err
}
