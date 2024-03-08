package jobs

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	//"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type EventsRepository struct {
	svc *eventbridge.Client
}

func NewEventsService() *EventsRepository {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &EventsRepository{
		svc: eventbridge.NewFromConfig(sdkConfig),
	}
}

func (events EventsRepository) startTickerScheduler() error {
	var err error

	ruleArn := os.Getenv("EVENTBRIDGE_RULE_ARN")
	ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")

	request := eventbridge.EnableRuleInput{
		Name:         aws.String(ruleName),
		EventBusName: aws.String(ruleArn),
	}

	_, err = events.svc.EnableRule(context.TODO(), &request)

	return err
}
func (events EventsRepository) stopTickerScheduler() error {
	var err error

	ruleArn := os.Getenv("EVENTBRIDGE_RULE_ARN")
	ruleName := os.Getenv("EVENTBRIDGE_RULE_NAME")

	request := eventbridge.DisableRuleInput{
		Name:         aws.String(ruleName),
		EventBusName: aws.String(ruleArn),
	}

	_, err = events.svc.DisableRule(context.TODO(), &request)

	return err
}
