package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubEventbridgeEnableRule(ruleName string, raiseErr *testtools.StubError) testtools.Stub {
	return testtools.Stub{
		OperationName: "EnableRule",
		Input:         &eventbridge.EnableRuleInput{Name: aws.String(ruleName)},
		Output:        &eventbridge.EnableRuleOutput{},
		Error:         raiseErr,
	}
}
