package logging_old

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"
)

func NewLogger(ctx context.Context) *zap.SugaredLogger {
	lc, _ := lambdacontext.FromContext(ctx)
	logger := zap.Must(zap.NewProduction())

	return logger.With(
		zap.String("requestId", lc.AwsRequestID),
	).Sugar()
}
