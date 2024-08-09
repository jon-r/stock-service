package events

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type events struct {
	eventBridgeClient *eventbridge.Client
	lambdaClient      *lambda.Client
}

type Scheduler interface {
	EnableRule(ruleName string) (*eventbridge.EnableRuleOutput, error)
	DisableRule(ruleName string) (*eventbridge.DisableRuleOutput, error)
	InvokeFunction(functionName string, payload interface{}) (*lambda.InvokeOutput, error)
}

func (e *events) EnableRule(ruleName string) (*eventbridge.EnableRuleOutput, error) {
	request := eventbridge.EnableRuleInput{
		Name: aws.String(ruleName),
	}

	return e.eventBridgeClient.EnableRule(context.TODO(), &request)
}

func (e *events) DisableRule(ruleName string) (*eventbridge.DisableRuleOutput, error) {
	request := eventbridge.DisableRuleInput{
		Name: aws.String(ruleName),
	}

	return e.eventBridgeClient.DisableRule(context.TODO(), &request)
}

func (e *events) InvokeFunction(functionName string, body interface{}) (*lambda.InvokeOutput, error) {
	var err error
	var payload []byte

	if body != nil {
		payload, err = json.Marshal(body)

		if err != nil {
			return nil, err
		}
	}

	lambdaReq := lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	}

	return e.lambdaClient.Invoke(context.TODO(), &lambdaReq)
}

func NewScheduler(config aws.Config) Scheduler {
	return &events{
		eventBridgeClient: eventbridge.NewFromConfig(config),
		lambdaClient:      lambda.NewFromConfig(config),
	}
}
