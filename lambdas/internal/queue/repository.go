package queue

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type QueueRepository struct {
	svc *sqs.Client
}

func NewQueueService() *QueueRepository {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &QueueRepository{
		svc: sqs.NewFromConfig(sdkConfig),
	}
}
