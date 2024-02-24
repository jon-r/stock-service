package queue

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type QueueRepository struct {
	svc *sqs.SQS
}

func NewQueueService(session session.Session) *QueueRepository {
	return &QueueRepository{
		svc: sqs.New(&session),
	}
}
