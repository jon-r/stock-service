package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	actual, _ := handleRequest(context.TODO(), events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Headers:    map[string]string{"Accept": "application/json"},
	})

	assert.Equal(t, actual.StatusCode, 200)
	assert.Contains(t, actual.Body, `"greeting":"hello LOGS world!"`)
}
