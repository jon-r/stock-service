package test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubLambdaInvoke(functionName string, payload []byte, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "Invoke",
		Input: &lambda.InvokeInput{
			FunctionName:   aws.String(functionName),
			InvocationType: types.InvocationTypeEvent,
			Payload:        payload,
		},
		Output: &lambda.InvokeOutput{},
		Error:  StubbedError(raiseErr),
	}
}
