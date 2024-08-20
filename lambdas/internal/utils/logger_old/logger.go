package logger_old

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//func NewLogger(ctx context.Context, level zapcore.Level) *zap.SugaredLogger {
//	lc, _ := lambdacontext.FromContext(ctx)
//
//	config := zap.NewProductionConfig()
//	config.Level = zap.NewAtomicLevelAt(level)
//
//	logger_old := zap.Must(config.Build())
//
//	return logger_old.With(
//		zap.String("requestId", lc.AwsRequestID),
//	).Sugar()
//}
//
//func LoadLambdaContext(ctx context.Context) {
//	lc, _ := lambdacontext.FromContext(ctx)
//
//
//}

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(level zapcore.Level) Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)

	sugaredLogger := zap.Must(config.Build()).Sugar()

	return Logger{sugaredLogger}
}

func (l *Logger) LoadLambdaContext(ctx context.Context) Logger {
	lc, _ := lambdacontext.FromContext(ctx)

	return Logger{l.With(zap.String("requestId", lc.AwsRequestID))}
}
