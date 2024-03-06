package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest() {
	log.Println("Hello world")
}

func main() {
	lambda.Start(handleRequest)
}
