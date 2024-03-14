package jobs

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func (events EventsRepository) InvokeWorker(job JobAction) error {
	var err error

	functionName := os.Getenv("LAMBDA_WORKER_NAME")

	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	lambdaReq := lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        payload,
	}

	_, err = events.lambda.Invoke(context.TODO(), &lambdaReq)

	return err
}
