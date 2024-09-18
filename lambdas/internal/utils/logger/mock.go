package logger

import "go.uber.org/zap/zapcore"

func NewMock() Logger {
	return New(zapcore.DPanicLevel) // todo raise once finished
}
