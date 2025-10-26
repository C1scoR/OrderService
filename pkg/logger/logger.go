package logger

import (
	"context"
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

const (
	loggerRequestIDKey = "x-request_id"
	loggerTraceIDKey   = "x-trace_id"
)

type L struct {
	z *zap.Logger
}

func NewLogger(env string) Logger {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if env == "dev" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	lo := L{z: logger}
	return &lo
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, loggerTraceIDKey, traceID)
}

func (l *L) Info(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Info(msg, fields...)
}

func (l *L) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Debug(msg, fields...)
}
func (l *L) Error(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Error(msg, fields...)
}
