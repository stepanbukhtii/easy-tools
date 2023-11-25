package grpc

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/stepanbukhtii/easy-tools/grpc/interceptor"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	PerRetryTimeout = 30 * time.Second
	MaxRetry        = uint(3)

	DefaultClientInterceptors = []grpc.UnaryClientInterceptor{
		interceptor.ClientLogger,
		retry.UnaryClientInterceptor(
			retry.WithMax(MaxRetry),
			retry.WithPerRetryTimeout(PerRetryTimeout),
			retry.WithBackoff(retry.BackoffExponential(100*time.Millisecond)),
		),
		interceptor.ClientAuth,
	}
	DefaultClientDialOptions = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(DefaultClientInterceptors...),
	}
)

func NewClientConnection(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address, DefaultClientDialOptions...)
}

type ClientConnectionParams struct {
	TransportCredentials credentials.TransportCredentials
	UnaryInterceptors    []grpc.UnaryClientInterceptor
	DialOptions          []grpc.DialOption
}

func NewClientConnectionParams(address string, params ClientConnectionParams) (*grpc.ClientConn, error) {
	if params.TransportCredentials == nil {
		params.TransportCredentials = insecure.NewCredentials()
	}

	unaryInterceptors := DefaultClientInterceptors
	if len(params.UnaryInterceptors) > 0 {
		unaryInterceptors = append(unaryInterceptors, params.UnaryInterceptors...)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(params.TransportCredentials),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
	}
	if params.DialOptions != nil {
		dialOpts = append(dialOpts, params.DialOptions...)
	}

	return grpc.NewClient(address, dialOpts...)
}
