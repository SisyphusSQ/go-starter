package log

import (
	"context"

	"go.uber.org/zap"
)

type LarkZapLogger struct {
	logger *zap.SugaredLogger
}

func NewLarkZapLogger(logger *zap.SugaredLogger) *LarkZapLogger {
	return &LarkZapLogger{logger: logger}
}

func (l *LarkZapLogger) Debug(ctx context.Context, args ...any) {
	l.logger.Debug(args...)
}

func (l *LarkZapLogger) Info(ctx context.Context, args ...any) {
	l.logger.Info(args...)
}

func (l *LarkZapLogger) Warn(ctx context.Context, args ...any) {
	l.logger.Warn(args...)
}

func (l *LarkZapLogger) Error(ctx context.Context, args ...any) {
	l.logger.Error(args...)
}
