package db

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	stubber := testtools.NewStubber()
	client := NewRepository(*stubber.SdkConfig)

	t.Run("HealthCheck", func(t *testing.T) {
		stubber.Add(testtools.Stub{
			OperationName: "ListTables",
			Input:         &dynamodb.ListTablesInput{},
			Output:        &dynamodb.ListTablesOutput{TableNames: []string{"Table1", "Table2"}},
		})

		assert.Equal(t, true, client.HealthCheck())
	})

	t.Run("AddOne", func(t *testing.T) {
		item := EntityBase{Id: "123", Sort: "ABC"}

		stubber.Add(testtools.Stub{
			OperationName: "PutItem",
			Input: &dynamodb.PutItemInput{
				TableName: aws.String("Table1"),
				Item: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: "123"},
					"SK": &types.AttributeValueMemberS{Value: "ABC"},
				},
			},
			Output: &dynamodb.PutItemOutput{},
		})

		_, err := client.AddOne("Table1", item)

		assert.NoError(t, err)
	})

	t.Run("AddMany", func(t *testing.T) {
		items := []EntityBase{{Id: "123", Sort: "ABC"}, {Id: "456", Sort: "DEF"}}

		stubber.Add(testtools.Stub{
			OperationName: "BatchWriteItem",
			Input: &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]types.WriteRequest{
					"Table1": {
						{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
							"PK": &types.AttributeValueMemberS{Value: "123"},
							"SK": &types.AttributeValueMemberS{Value: "ABC"},
						}}},
						{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
							"PK": &types.AttributeValueMemberS{Value: "456"},
							"SK": &types.AttributeValueMemberS{Value: "DEF"},
						}}},
					},
				},
			},
			Output: &dynamodb.BatchWriteItemOutput{},
		})

		res, err := client.AddMany("Table1", items)

		assert.NoError(t, err)
		assert.Equal(t, 2, res)
	})
	t.Run("Update", func(t *testing.T) {
		item := EntityBase{Id: "123", Sort: "ABC"}

		updateEx := expression.Set(
			expression.Name("UpdatedValue"),
			expression.Value("New Value"),
		)
		update, _ := expression.NewBuilder().WithUpdate(updateEx).Build()

		stubber.Add(testtools.Stub{
			OperationName: "UpdateItem",
			Input: &dynamodb.UpdateItemInput{
				TableName: aws.String("Table1"),
				Key: EntityKey{
					"PK": &types.AttributeValueMemberS{Value: "123"},
					"SK": &types.AttributeValueMemberS{Value: "ABC"},
				},
				ExpressionAttributeNames: map[string]string{"#0": "UpdatedValue"},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":0": &types.AttributeValueMemberS{Value: "New Value"},
				},
				UpdateExpression: aws.String("SET #0 = :0\n"),
			},
			Output: &dynamodb.UpdateItemOutput{},
		})

		_, err := client.Update("Table1", item.GetKey(), update)

		assert.NoError(t, err)
	})
	t.Run("GetMany", func(t *testing.T) {
		item := EntityBase{Id: "123", Sort: "ABC"}

		entity, _ := attributevalue.MarshalMap(item)

		filterEx := expression.Name("PK").Equal(expression.Value("123"))
		query, _ := expression.NewBuilder().WithFilter(filterEx).Build()

		stubber.Add(testtools.Stub{
			OperationName: "Scan",
			Input: &dynamodb.ScanInput{
				TableName:                aws.String("Table1"),
				ExpressionAttributeNames: map[string]string{"#0": "PK"},
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":0": &types.AttributeValueMemberS{Value: "123"},
				},
				FilterExpression: aws.String("#0 = :0"),
			},
			Output: &dynamodb.ScanOutput{
				Count: 1,
				Items: []map[string]types.AttributeValue{{
					"PK": &types.AttributeValueMemberS{Value: "123"},
					"SK": &types.AttributeValueMemberS{Value: "ABC"},
				}},
			},
		})

		res, err := client.GetMany("Table1", query)
		log.Println(err)

		assert.NoError(t, err)
		assert.Equal(t, []map[string]types.AttributeValue{entity}, res)
	})
}
