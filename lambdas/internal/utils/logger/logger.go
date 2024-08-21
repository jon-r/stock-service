package logger

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	LoadContext(ctx context.Context)

	Sync() error

	Errorw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Errorf(format string, args ...interface{})

	Infoln(args ...interface{})
	Debugln(args ...interface{})
}

type l struct {
	instance *zap.SugaredLogger
}

func (l *l) LoadContext(ctx context.Context) {
	lc, _ := lambdacontext.FromContext(ctx)

	l.instance = l.instance.With(zap.String("requestId", lc.AwsRequestID))
}

func (l *l) Errorf(format string, args ...interface{}) {
	l.instance.
		With(zap.String("source", getFunctionSource())).
		Errorf(format, args...)
}
func (l *l) Errorw(msg string, keysAndValues ...interface{}) {
	l.instance.
		With(zap.String("source", getFunctionSource())).
		Errorw(msg, keysAndValues...)
}
func (l *l) Warnw(msg string, keysAndValues ...interface{}) {
	l.instance.
		With(zap.String("source", getFunctionSource())).
		Warnw(msg, keysAndValues...)
}
func (l *l) Infoln(args ...interface{}) {
	l.instance.
		With(zap.String("source", getFunctionSource())).
		Infoln(args...)
}
func (l *l) Debugw(msg string, keysAndValues ...interface{}) {
	l.instance.
		With(zap.String("source", getFunctionSource())).
		Debugw(msg, keysAndValues...)
}
func (l *l) Debugln(args ...interface{}) {
	l.instance.
		With(zap.String("source", getFunctionSource())).
		Debugln(args...)
}

func (l *l) Sync() error {
	return l.instance.Sync()
}

func getFunctionSource() string {
	pc, _, line, ok := runtime.Caller(2)
	if ok {
		fullPath := runtime.FuncForPC(pc).Name()
		fnName := fullPath[strings.LastIndex(fullPath, "/")+1:]

		return fmt.Sprintf("%s@%d", fnName, line)
	} else {
		return ""
	}
}

func New(level zapcore.Level) Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.DisableCaller = true // replaced by getFunctionSource because the logger is called indirectly

	sugaredLogger := zap.Must(config.Build()).Sugar()

	return &l{sugaredLogger}
}
