package db

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"log"
)

var JobsTableName = "stock-app_Job"

type JobInput struct {
	Name  string
	Group string
	Speed int
}

type JobItem struct {
	JobId string
	Name  string
	Group string
}

func (db DatabaseRepository) InsertJobs(jobInputs []JobInput) error {
	var err error

	writeRequests := make([]*dynamodb.WriteRequest, len(jobInputs))

	for i, jobInput := range jobInputs {
		job := JobItem{
			JobId: uuid.NewString(),
			Name:  jobInput.Name,
			Group: jobInput.Group,
		}

		av, err := dynamodbattribute.MarshalMap(job)
		if err != nil {
			log.Print("could not convert input to JobItem")
			break
		}
		writeRequests[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: av},
		}
	}

	if err != nil {
		return err
	}

	input := dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			JobsTableName: writeRequests,
		},
	}
	_, err = db.svc.BatchWriteItem(&input)

	return err
}
