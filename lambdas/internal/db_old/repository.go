package db_old

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DatabaseRepository struct {
	svc             *dynamodb.Client
	StocksTableName *string
	LogsTableName   *string
}

func CreateDatabaseClient() *dynamodb.Client {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(sdkConfig)
}

func NewDatabaseService(client *dynamodb.Client) *DatabaseRepository {
	return &DatabaseRepository{
		svc:             client,
		StocksTableName: aws.String(os.Getenv("DB_STOCKS_TABLE_NAME")),
		LogsTableName:   aws.String(os.Getenv("DB_LOGS_TABLE_NAME")),
	}
}

func (item *StocksTableItem) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(item.Id)
	if err != nil {
		panic(err)
	}
	sort, err := attributevalue.Marshal(item.Sort)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"PK": id, "SK": sort}
}

func (item *StocksTableItem) SetKey(partitionKeyType KeyType, partitionId string, sortKeyType KeyType, sortId string) {
	partitionKey := string(partitionKeyType) + partitionId
	sortKey := string(sortKeyType) + sortId

	item.Id = partitionKey
	item.Sort = sortKey
}
