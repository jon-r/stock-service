package testutil

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/jon-r/stock-service/lambdas/internal/jobs"
)

type mockSQSClient struct {
	sqs.Client
	Res sqs.SendMessageOutput
}

type mockedQueueRepository struct {
	jobs.QueueRepository
	svc mockSQSClient
}

func (m mockSQSClient) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return &m.Res, nil
}

// todo make service clients that dont call the real SDK, but just check the response
//
//	 https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/testtools/awsm_stubber.go
//		https://github.com/aws/aws-sdk-go/blob/main/example/service/sqs/mockingClientsForTests/ifaceExample_test.go
func mockNewQueueService() mockedQueueRepository {
	repo := jobs.NewQueueService()
}
