package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Warning(msg string, fields ...Field)
}

type logger struct {
	zap *zap.Logger
}

func (l logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

func (l logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

func (l logger) Warning(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

func New(namespace string) Logger {
	return logger{
		zap: newZapLogger(namespace),
	}
}
