package logger

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

//nolint:sloglint
func NewServerGRPCMiddleware(log *slog.Logger) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(
		logging.LoggerFunc(
			func(ctx context.Context, level logging.Level, msg string, fields ...any) {
				log.Log(ctx, slog.Level(level), msg, fields...)
			},
		),
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	)
}

//nolint:sloglint
func NewClientGRPCMiddleware(log *slog.Logger) grpc.UnaryClientInterceptor {
	return logging.UnaryClientInterceptor(
		logging.LoggerFunc(
			func(ctx context.Context, level logging.Level, msg string, fields ...any) {
				log.Log(ctx, slog.Level(level), msg, fields...)
			},
		),
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	)
}
