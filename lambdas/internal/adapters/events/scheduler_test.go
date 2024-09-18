package events

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	stubber := testtools.NewStubber()
	client := NewScheduler(*stubber.SdkConfig)

	t.Run("EnableRule", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "EnableRule",
			Input:         &eventbridge.EnableRuleInput{Name: aws.String("rule1")},
			Output:        &eventbridge.EnableRuleOutput{},
		})

		_, err := client.EnableRule("rule1")

		assert.NoError(t, err)
	})

	t.Run("DisableRule", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "DisableRule",
			Input:         &eventbridge.DisableRuleInput{Name: aws.String("rule1")},
			Output:        &eventbridge.DisableRuleOutput{},
		})

		_, err := client.DisableRule("rule1")

		assert.NoError(t, err)
	})

	t.Run("InvokeFunction", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "Invoke",
			Input: lambda.InvokeInput{
				FunctionName:   aws.String("function1"),
				Payload:        []byte(`{"Payload":"do something"}`),
				InvocationType: "Event",
			},
			Output: &lambda.InvokeOutput{},
		})

		type action struct{ Payload string }
		_, err := client.InvokeFunction("function1", action{Payload: "do something"})

		assert.NoError(t, err)
	})

}
