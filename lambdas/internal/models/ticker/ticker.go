package ticker

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func NewTickerEntity(params NewTickerParams) *Entity {
	entity := &Entity{
		Provider: params.Provider,
	}
	entity.SetKey(KeyTicker, params.TickerId, KeyTickerId, params.TickerId)

	return entity
}

func ParamsFromJsonString(jsonString string) (NewTickerParams, error) {
	var params NewTickerParams
	err := json.Unmarshal([]byte(jsonString), &params)

	return params, err
}

func NewFromJsonString(jsonString string) (*Entity, error) {
	var params NewTickerParams
	err := json.Unmarshal([]byte(jsonString), &params)

	if err != nil {
		return nil, err
	}

	return NewTickerEntity(params), nil
}

func TableName() string {
	return os.Getenv("DB_STOCKS_TABLE_NAME")
}

func NewSelectAllQuery() (expression.Expression, error) {
	filterEx := expression.Name("SK").BeginsWith(string(KeyTickerId))
	projEx := expression.NamesList(
		expression.Name("SK"), expression.Name("Provider"),
	)

	return expression.NewBuilder().WithFilter(filterEx).WithProjection(projEx).Build()
}
