package db

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type EntityBase struct {
	Id   string `dynamodbav:"PK"`
	Sort string `dynamodbav:"SK"`
}

//type TickerItem struct {
//	EntityBase
//	Provider    providers.ProviderName
//	Description providers.TickerDescription
//}

//type PriceItem struct {
//	EntityBase
//	Price providers.TickerPrices
//	Date  string `dynamodbav:"DT"`
//}

type KeyType string

// todo move these to separate entity sub-packages
const (
	KeyTicker         KeyType = "T#"
	KeyTickerPrice    KeyType = "P#"
	KeyTickerId       KeyType = "T#"
	KeyTickerDividend KeyType = "D#"

	KeyUser        KeyType = "U#"
	KeyUserTicker  KeyType = "T#"
	KeyUserTxEvent KeyType = "E#"
)

func (item *EntityBase) GetKey() map[string]types.AttributeValue {
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

func (item *EntityBase) SetKey(partitionKeyType KeyType, partitionId string, sortKeyType KeyType, sortId string) {
	partitionKey := string(partitionKeyType) + partitionId
	sortKey := string(sortKeyType) + sortId

	item.Id = partitionKey
	item.Sort = sortKey
}
