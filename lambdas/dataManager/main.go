package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type JobEvent struct {
	JobId string `json:"jobId"`
	Name  string `json:"name"`
	Delay int    `json:"delay"`
}

func loadData(ctx context.Context, event JobEvent) {
	//names := []string{
	//	"Phoebe",
	//	"Harley",
	//	"Bandit",
	//	"Delilah",
	//	"Tiger",
	//	"Panda",
	//	"Whiskey",
	//	"Jasper",
	//	"Belle",
	//	"Shelby",
	//	"Zara",
	//	"Bruno",
	//}

	dbSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dbService := dynamodb.New(dbSession)

	av, err := dynamodbattribute.MarshalMap(event)
	if err != nil {
		log.Fatalf("Error marshalling new job item: %s", err)
	}

	tableName := "JobTable"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dbService.PutItem(input)
	if err != nil {
		log.Fatalf("Error calling PutItem: %s", err)
	} else {
		log.Println("Succesfully added '" + event.Name + "' to table " + tableName)
	}
}

func main() {
	lambda.Start(loadData)
}
