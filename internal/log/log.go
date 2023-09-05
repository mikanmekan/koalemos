package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	Panic(msg string, fields ...zapcore.Field)

	// Debugf(template string, msg string, fields ...zapcore.Field)
	// Infof(template string, msg string, fields ...zapcore.Field)
	// Warnf(template string, msg string, fields ...zapcore.Field)
	// Errorf(template string, msg string, fields ...zapcore.Field)
	// Fatalf(template string, msg string, fields ...zapcore.Field)
	// Panicf(template string, msg string, fields ...zapcore.Field)

	With(fields ...zapcore.Field) *zap.Logger
}

type LogImpl struct {
	*zap.Logger
}

func NewLogger() Logger {
	logger, _ := zap.NewProduction()

	return &LogImpl{logger}
}

var _ Logger = (*LogImpl)(nil)

func (l *LogImpl) Fatal(msg string, fields ...zapcore.Field) {
	l.Logger.Fatal(msg, fields...)
	l.Logger.Sync()
}
