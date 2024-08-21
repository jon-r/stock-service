package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jon-r/stock-service/lambdas/internal/handlers"
	"github.com/jon-r/stock-service/lambdas/internal/models/job"
)

type handler struct{ *handlers.LambdaHandler }

var dataWorkerHandler = handler{handlers.NewLambdaHandler()}

func (h *handler) HandleRequest(ctx context.Context, j job.Job) error {
	h.Log.LoadContext(ctx)
	defer h.Log.Sync()

	// 1. handle action
	err := h.doJob(j)

	if err == nil {
		h.Log.Infoln("job completed", "jobId", j.JobId)
		return nil // job done
	}

	// 2. if action failed or new queue actions after last, try again
	h.Log.Warnw("failed to process event, re-adding it to queue",
		"jobId", j.JobId,
		"error", err,
	)

	queueErr := h.Jobs.RequeueJob(j, err.Error())

	if queueErr != nil {
		h.Log.Errorw("Failed to add item to DLQ",
			"jobId", j.JobId,
			"error", queueErr,
		)
		return queueErr
	}

	return err
}

func main() {
	lambda.Start(dataWorkerHandler.HandleRequest)
}
