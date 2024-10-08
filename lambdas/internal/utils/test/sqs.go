package test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubSqsSendMessage(queue string, message string, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "SendMessage",
		Input:         &sqs.SendMessageInput{QueueUrl: aws.String(queue), MessageBody: aws.String(message)},
		Output:        &sqs.SendMessageOutput{},
		Error:         StubbedError(raiseErr),
	}
}

func StubSqsSendMessageBatch(queue string, items []types.SendMessageBatchRequestEntry, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "SendMessageBatch",
		Input:         &sqs.SendMessageBatchInput{QueueUrl: aws.String(queue), Entries: items},
		Output:        &sqs.SendMessageBatchOutput{},
		Error:         StubbedError(raiseErr),
	}
}

func StubSqsReceiveMessages(request *sqs.ReceiveMessageInput, response *sqs.ReceiveMessageOutput, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "ReceiveMessage",
		Input:         request,
		Output:        response,
		Error:         StubbedError(raiseErr),
	}
}

func StubSqsDeleteMessage(queue string, messageId string, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "DeleteMessage",
		Input:         &sqs.DeleteMessageInput{QueueUrl: aws.String(queue), ReceiptHandle: aws.String(messageId)},
		Output:        &sqs.DeleteMessageOutput{},
		Error:         StubbedError(raiseErr),
	}
}
