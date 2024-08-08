package scheduler_old

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type EventsRepository struct {
	svc    *eventbridge.Client
	lambda *lambda.Client
}

func CreateEventClients() (*eventbridge.Client, *lambda.Client) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return eventbridge.NewFromConfig(sdkConfig), lambda.NewFromConfig(sdkConfig)
}

func NewEventsService(client *eventbridge.Client, lambdaClient *lambda.Client) *EventsRepository {
	return &EventsRepository{
		svc:    client,
		lambda: lambdaClient,
	}
}

func (events EventsRepository) StartTickerScheduler() error {
	var err error

	ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")

	request := eventbridge.EnableRuleInput{
		Name: aws.String(ruleName),
	}

	_, err = events.svc.EnableRule(context.TODO(), &request)

	if err != nil {
		return err
	}

	lambdaErr := events.InvokeTicker()

	if lambdaErr != nil {
		log.Printf("Failed to manually trigger poller but continuing anyway: %v\n", lambdaErr)
	}

	return err
}

func (events EventsRepository) InvokeTicker() error {
	functionName := os.Getenv("LAMBDA_TICKER_NAME")
	lambdaReq := lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: types.InvocationTypeEvent,
	}

	_, err := events.lambda.Invoke(context.TODO(), &lambdaReq)

	return err
}

func (events EventsRepository) StopTickerScheduler() error {
	var err error

	ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")

	request := eventbridge.DisableRuleInput{
		Name: aws.String(ruleName),
	}

	_, err = events.svc.DisableRule(context.TODO(), &request)

	return err
}
