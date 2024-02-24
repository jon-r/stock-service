package db

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DatabaseRepository struct {
	svc *dynamodb.DynamoDB
}

func NewDatabaseService(session session.Session) *DatabaseRepository {
	return &DatabaseRepository{
		svc: dynamodb.New(&session),
	}
}
