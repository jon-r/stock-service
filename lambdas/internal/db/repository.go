package db

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//type DatabaseService interface {
//	NewDatabaseService(session session.Session) *DatabaseRepository
//
//	InsertJobs(jobInputs []JobInput) error
//}

type DatabaseRepository struct {
	svc *dynamodb.DynamoDB
}

func NewDatabaseService(session session.Session) *DatabaseRepository {
	return &DatabaseRepository{
		svc: dynamodb.New(&session),
	}
}

//type JobInput struct {
//	Name string
//	Group string
//	Speed int
//}
//type JobItem struct {
//	JobId string
//	Name string
//	Group string
//}
//
//var jobTableName = "stock-app_Job"
//
//func (db DatabaseRepository) InsertJobs(jobInputs []JobInput) error {
//	var err error
//
//	writeRequests := make([]*dynamodb.WriteRequest, len(jobInputs))
//
//	for i, jobInput := range jobInputs {
//		job := JobItem{
//			JobId: uuid.NewString(),
//			Name:  jobInput.Name,
//			Group: jobInput.Group,
//		}
//
//		av, err := dynamodbattribute.MarshalMap(job)
//		if err != nil {
//			log.Print("could not convert input to JobItem")
//			break
//		}
//		writeRequests[i] = &dynamodb.WriteRequest{
//			PutRequest: &dynamodb.PutRequest{ Item: av },
//		}
//	}
//
//	if err != nil { return err }
//
//	dbInput := &dynamodb.BatchWriteItemInput{
//		RequestItems: map[string][]*dynamodb.WriteRequest{
//			jobTableName: writeRequests,
//		},
//	}
//	_, err = db.svc.BatchWriteItem(dbInput)
//
//	return err
//}
