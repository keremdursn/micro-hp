package logging

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"

var GlobalLogger *Logger

func InitLogger(serviceName string) *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"

	// JSON format for structured logging
	config.Encoding = "json"

	// Log level configuration
	if os.Getenv("LOG_LEVEL") == "debug" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// Add service name to all logs
	logger = logger.With(zap.String("service", serviceName))

	GlobalLogger = &Logger{logger}
	return GlobalLogger
}

func (l *Logger) WithCorrelationID(ctx context.Context) *zap.Logger {
	if corrID := ctx.Value(CorrelationIDKey); corrID != nil {
		return l.With(zap.String("correlation_id", corrID.(string)))
	}
	return l.Logger
}

func (l *Logger) LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, userID *uint) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
		zap.String("type", "http_request"),
	}

	if userID != nil {
		fields = append(fields, zap.Uint("user_id", *userID))
	}

	logger := l.WithCorrelationID(ctx)

	if statusCode >= 500 {
		logger.Error("HTTP Request", fields...)
	} else if statusCode >= 400 {
		logger.Warn("HTTP Request", fields...)
	} else {
		logger.Info("HTTP Request", fields...)
	}
}

func (l *Logger) LogError(ctx context.Context, err error, message string, fields ...zap.Field) {
	logger := l.WithCorrelationID(ctx)
	allFields := append(fields, zap.Error(err), zap.String("type", "error"))
	logger.Error(message, allFields...)
}

func (l *Logger) LogInfo(ctx context.Context, message string, fields ...zap.Field) {
	logger := l.WithCorrelationID(ctx)
	allFields := append(fields, zap.String("type", "info"))
	logger.Info(message, allFields...)
}

func (l *Logger) LogDebug(ctx context.Context, message string, fields ...zap.Field) {
	logger := l.WithCorrelationID(ctx)
	allFields := append(fields, zap.String("type", "debug"))
	logger.Debug(message, allFields...)
}

func (l *Logger) LogDatabaseOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Duration("duration", duration),
		zap.String("type", "database"),
	}

	logger := l.WithCorrelationID(ctx)

	if err != nil {
		fields = append(fields, zap.Error(err))
		logger.Error("Database Operation Failed", fields...)
	} else {
		logger.Info("Database Operation", fields...)
	}
}

func (l *Logger) LogServiceCall(ctx context.Context, targetService, endpoint string, statusCode int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("target_service", targetService),
		zap.String("endpoint", endpoint),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
		zap.String("type", "service_call"),
	}

	logger := l.WithCorrelationID(ctx)

	if err != nil {
		fields = append(fields, zap.Error(err))
		logger.Error("Service Call Failed", fields...)
	} else if statusCode >= 400 {
		logger.Warn("Service Call Warning", fields...)
	} else {
		logger.Info("Service Call", fields...)
	}
}
