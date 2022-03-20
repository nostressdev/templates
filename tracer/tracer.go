package tracer

import (
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"os"
)

const (
	jaegerEndpoint = "JAEGER"
)

func NewJaegerWithEnv() (trace.TracerProvider, error) {
	endpoint := os.Getenv(jaegerEndpoint)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}
	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
	), nil
}

func NewStdoutTracer() (trace.TracerProvider, error) {
	exp, err := stdouttrace.New()
	if err != nil {
		return nil, err
	}
	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
	), nil
}

func NewDefault() (trace.TracerProvider, error) {
	if os.Getenv(jaegerEndpoint) == "" {
		return NewStdoutTracer()
	}
	return NewJaegerWithEnv()
}
