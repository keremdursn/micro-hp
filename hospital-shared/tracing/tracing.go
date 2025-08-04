package tracing

import (
	"context"
	"io"
	"log"
)

// MockCloser implements io.Closer
type MockCloser struct{}

func (m *MockCloser) Close() error {
	return nil
}

// InitTracing initializes simple tracing (fallback to simple implementation)
func InitTracing(serviceName string) (io.Closer, error) {
	err := InitSimpleTracing(serviceName)
	if err != nil {
		log.Printf("Could not initialize simple tracer: %s", err.Error())
		return &MockCloser{}, err
	}

	log.Printf("Simple tracing initialized for service: %s", serviceName)
	return &MockCloser{}, nil
}

// Wrapper functions for simple spans that match the old interface
func StartHTTPSpan(ctx context.Context, operationName, method, url string) (*SimpleSpan, context.Context) {
	span := StartSimpleSpan(ctx, operationName)
	span.SetTag("http.method", method)
	span.SetTag("http.url", url)
	span.SetTag("component", "http-client")
	return span, ctx
}

func StartDatabaseSpan(ctx context.Context, operation, table string) (*SimpleSpan, context.Context) {
	span := StartSimpleSpan(ctx, "db."+operation)
	span.SetTag("db.type", "postgresql")
	span.SetTag("db.statement", operation)
	if table != "" {
		span.SetTag("db.table", table)
	}
	span.SetTag("component", "database")
	return span, ctx
}

func StartServiceSpan(ctx context.Context, serviceName, operation string) (*SimpleSpan, context.Context) {
	span := StartSimpleSpan(ctx, serviceName+"."+operation)
	span.SetTag("service.name", serviceName)
	span.SetTag("service.operation", operation)
	span.SetTag("component", "microservice")
	return span, ctx
}

func FinishSpanWithError(span *SimpleSpan, err error) {
	if span == nil {
		return
	}
	FinishSimpleSpanWithError(span, err)
}
