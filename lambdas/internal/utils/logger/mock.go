package logger

import "go.uber.org/zap/zapcore"

func NewMock() Logger {
	return New(zapcore.DebugLevel) // todo raise once finished
}
