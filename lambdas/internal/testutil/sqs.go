package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubSqsSendMessageBatch(queue string, items []types.SendMessageBatchRequestEntry, raiseErr *testtools.StubError) testtools.Stub {
	return testtools.Stub{
		OperationName: "SendMessageBatch",
		Input:         &sqs.SendMessageBatchInput{QueueUrl: aws.String(queue), Entries: items},
		Output:        &sqs.SendMessageBatchOutput{},
		Error:         raiseErr,
	}
}