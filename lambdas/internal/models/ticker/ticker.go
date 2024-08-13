package ticker

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func NewTickerEntity(params *NewTickerParams) *Entity {
	entity := &Entity{
		Provider: params.Provider,
	}
	entity.SetKey(KeyTicker, params.TickerId, KeyTickerId, params.TickerId)

	return entity
}

func NewParamsFromJsonString(jsonString string) (*NewTickerParams, error) {
	var params NewTickerParams
	err := json.Unmarshal([]byte(jsonString), &params)

	return &params, err
}

func NewStubsFromDynamoDb(entities []map[string]types.AttributeValue) (*[]EntityStub, error) {
	var tickers []EntityStub
	err := attributevalue.UnmarshalListOfMaps(entities, &tickers)

	return &tickers, err
}

func TableName() string {
	return os.Getenv("DB_STOCKS_TABLE_NAME")
}
