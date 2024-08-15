package ticker

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jon-r/stock-service/lambdas/internal/models/provider"
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

func NewStubsFromDynamoDb(entities []map[string]dbTypes.AttributeValue) (*[]EntityStub, error) {
	tickers := make([]EntityStub, len(entities))
	err := attributevalue.UnmarshalListOfMaps(entities, &tickers)

	return &tickers, err
}

//func NewStubsFromSQS(messages []sqsTypes.Message) (*[]EntityStub, error) {
//	var err error
//
//	tickers := make([]EntityStub, len(messages))
//	for i, message := range messages {
//		err = json.Unmarshal([]byte(*message.Body), &tickers[i])
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &tickers, nil
//}

func TableName() string {
	return os.Getenv("DB_STOCKS_TABLE_NAME")
}

func GroupByProvider(tickers []EntityStub) map[provider.Name][]string {
	list := map[provider.Name][]string{}

	for _, item := range tickers {
		key := item.Provider
		tickerId, _ := strings.CutPrefix(item.TickerSort, "T#")

		list[key] = append(list[key], tickerId)
	}

	return list
}
