package interceptor

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var resp any
	var err error
	defer func() {
		if r := recover(); r != nil {
			var recoverErr error
			switch x := r.(type) {
			case string:
				recoverErr = errors.New(x)
			case error:
				recoverErr = x
			}

			slog.With(
				slog.String(string(semconv.ExceptionStacktraceKey), string(debug.Stack())),
				slog.Any(string(semconv.ExceptionMessageKey), recoverErr),
			).Error("panic recovered")

			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()

	resp, err = handler(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
