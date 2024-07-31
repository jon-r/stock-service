package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubDynamoDbAddTicker(tableName string, item map[string]types.AttributeValue, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "PutItem",
		Input:         &dynamodb.PutItemInput{Item: item, TableName: aws.String(tableName)},
		Output:        &dynamodb.PutItemOutput{},
		Error:         StubbedError(raiseErr),
	}
}
