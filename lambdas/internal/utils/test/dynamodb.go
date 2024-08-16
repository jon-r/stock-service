package test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
)

func StubDynamoDbPutItem(tableName string, item interface{}, raiseErr error) testtools.Stub {
	data, _ := attributevalue.MarshalMap(item)

	return testtools.Stub{
		OperationName: "PutItem",
		Input:         &dynamodb.PutItemInput{Item: data, TableName: aws.String(tableName)},
		Output:        &dynamodb.PutItemOutput{},
		Error:         StubbedError(raiseErr),
	}
}

func StubDynamoDbUpdate(request *dynamodb.UpdateItemInput, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "UpdateItem",
		Input:         request,
		Output:        &dynamodb.UpdateItemOutput{},
		Error:         StubbedError(raiseErr),
	}
}

func StubDynamoDbBatchWriteTicker(request *dynamodb.BatchWriteItemInput, raiseErr error) testtools.Stub {
	return testtools.Stub{
		OperationName: "BatchWriteItem",
		Input:         request,
		Output:        &dynamodb.BatchWriteItemOutput{},
		Error:         StubbedError(raiseErr),
	}
}

func StubDynamoDbScan(request *dynamodb.ScanInput, response interface{}, raiseErr error) testtools.Stub {
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
