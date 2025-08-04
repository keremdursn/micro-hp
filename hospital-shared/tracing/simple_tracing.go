package tracing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Simple tracing without complex dependencies
type SimpleSpan struct {
	TraceID       string            `json:"traceId"`
	SpanID        string            `json:"id"`
	Name          string            `json:"name"`
	Timestamp     int64             `json:"timestamp"`
	Duration      int64             `json:"duration"`
	LocalEndpoint LocalEndpoint     `json:"localEndpoint"`
	Tags          map[string]string `json:"tags"`
	started       time.Time
}

type LocalEndpoint struct {
	ServiceName string `json:"serviceName"`
}

type SimpleTracer struct {
	serviceName string
	zipkinURL   string
	client      *http.Client
}

var simpleTracer *SimpleTracer

func InitSimpleTracing(serviceName string) error {
	simpleTracer = &SimpleTracer{
		serviceName: serviceName,
		zipkinURL:   "http://zipkin:9411",
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
	return nil
}

func StartSimpleSpan(ctx context.Context, operationName string) *SimpleSpan {
	// Check if tracer is initialized
	if simpleTracer == nil {
		return nil
	}

	now := time.Now()
	span := &SimpleSpan{
		TraceID:   generateID(),
		SpanID:    generateID(),
		Name:      operationName,
		Timestamp: now.UnixMicro(),
		LocalEndpoint: LocalEndpoint{
			ServiceName: simpleTracer.serviceName,
		},
		Tags:    make(map[string]string),
		started: now,
	}

	// Store in context if needed
	return span
}

func (s *SimpleSpan) SetTag(key, value string) {
	if s.Tags == nil {
		s.Tags = make(map[string]string)
	}
	s.Tags[key] = value
}

func (s *SimpleSpan) Finish() {
	if simpleTracer == nil {
		return
	}

	s.Duration = time.Since(s.started).Microseconds()

	// Send to Zipkin
	go func() {
		spans := []SimpleSpan{*s}
		jsonData, err := json.Marshal(spans)
		if err != nil {
			log.Printf("❌ Failed to marshal trace: %v", err)
			return
		}

		log.Printf("🚀 Sending trace to Zipkin: %s", string(jsonData))

		resp, err := simpleTracer.client.Post(
			simpleTracer.zipkinURL+"/api/v2/spans",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			log.Printf("❌ Failed to send trace to Zipkin: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("✅ Trace sent to Zipkin, status: %d", resp.StatusCode)
	}()
}

func generateID() string {
	return fmt.Sprintf("%016x", time.Now().UnixNano())
}

func FinishSimpleSpanWithError(span *SimpleSpan, err error) {
	if err != nil {
		span.SetTag("error", "true")
		span.SetTag("error.message", err.Error())
	}
	span.Finish()
}
