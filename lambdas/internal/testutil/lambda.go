package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubLambdaInvoke(functionName string, raiseErr *testtools.StubError) testtools.Stub {
	return testtools.Stub{
		OperationName: "Invoke",
		Input: &lambda.InvokeInput{
			FunctionName:   aws.String(functionName),
			InvocationType: types.InvocationTypeEvent,
		},
		Output: &lambda.InvokeOutput{},
		Error:  raiseErr,
	}
}
