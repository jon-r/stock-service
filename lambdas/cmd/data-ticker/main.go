package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest() {
	log.Println("Hello world")

	// 1. get all items in queue

	// 2. if queue is empty, disable the event rule and end the function

	// 3. group queue jobs by provider

	// 4 for each provider have a ticker function that invokes event provider/ticker/type to the worker fn
}

func main() {
	lambda.Start(handleRequest)
}
