package logger

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(ctx context.Context, level zapcore.Level) *zap.SugaredLogger {
	lc, _ := lambdacontext.FromContext(ctx)

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)

	logger := zap.Must(config.Build())

	return logger.With(
		zap.String("requestId", lc.AwsRequestID),
	).Sugar()
}
