package grpc

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func newUnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return grpc_zap.UnaryServerInterceptor(logger)
}

func newServerLoggingInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return grpc_zap.StreamServerInterceptor(logger)
}

func NewServer(
	logger *zap.Logger,
	tp trace.TracerProvider,
) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			newServerLoggingInterceptor(logger),
			otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(tp)),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			newUnaryLoggingInterceptor(logger),
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(tp)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		),
	)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(server)
	return server
}
