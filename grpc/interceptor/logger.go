package interceptor

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/elog"
)

func ServerLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Now().UTC().Sub(start)

	if err := logRequest(ctx, info.FullMethod, duration, req, err); err != nil {
		return resp, err
	}

	return resp, err
}

func ClientLogger(
	ctx context.Context,
	method string,
	req, reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	duration := time.Now().UTC().Sub(start)

	if err := logRequest(ctx, method, duration, req, err); err != nil {
		return err
	}

	return err
}

func logRequest(ctx context.Context, method string, duration time.Duration, req any, err error) error {
	var requestData []byte
	if req != nil {
		var err error
		requestData, err = json.Marshal(req)
		if err != nil {
			return err
		}
	}

	logger := econtext.Logger(ctx)
	if !logger.Enabled(ctx, slog.LevelError) {
		logger = slog.Default()
	}

	logger = logger.With(
		slog.String(string(semconv.RPCSystemNameKey), "grpc"),
		slog.String(string(semconv.RPCMethodKey), method),
		slog.Duration(elog.RPCRequestDuration, duration),
	)

	if len(requestData) > 0 {
		logger = logger.With(slog.String(elog.RPCRequestBodyContent, string(requestData)))
	}

	if errStatus, ok := status.FromError(err); ok {
		logger = logger.With(slog.String(string(semconv.RPCResponseStatusCodeKey), errStatus.Code().String()))
	}

	switch {
	case err != nil:
		logger.With(elog.Err(err)).ErrorContext(ctx, "gRPC request")
	default:
		logger.InfoContext(ctx, "gRPC request")
	}

	return nil
}
