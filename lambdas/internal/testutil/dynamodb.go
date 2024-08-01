package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

// todo pass in any and marshal here (dry out the test?)
func StubDynamoDbAddTicker(tableName string, item map[string]types.AttributeValue, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "PutItem",
		Input:         &dynamodb.PutItemInput{Item: item, TableName: aws.String(tableName)},
		Output:        &dynamodb.PutItemOutput{},
		Error:         StubbedError(raiseErr),
	}
}

func StubDynamoDbScan(request *dynamodb.ScanInput, response interface{}, raiseErr error) testtools.Stub {
	// query, _ := attributevalue.MarshalMap(request)
	list := unpackArray(response)

	var items []map[string]types.AttributeValue
	for _, item := range list {
		data, _ := attributevalue.MarshalMap(item)
		items = append(items, data)
	}

	return testtools.Stub{
		OperationName: "Scan",
		Input:         request,
		Output:        &dynamodb.ScanOutput{Items: items},
		Error:         StubbedError(raiseErr),
	}
}
