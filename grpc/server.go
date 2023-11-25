package grpc

import (
	"fmt"
	"net"

	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/grpc/interceptor"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var (
	DefaultServerInterceptors = []grpc.UnaryServerInterceptor{
		interceptor.Recovery,
		interceptor.ServerLogger,
	}
)

type Server struct {
	*grpc.Server
}

func NewServer(interceptors ...grpc.UnaryServerInterceptor) Server {
	return Server{
		Server: grpc.NewServer(
			grpc.ChainUnaryInterceptor(append(DefaultServerInterceptors, interceptors...)...),
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		),
	}
}

func (s *Server) Serve(cfg config.GRPC) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		return err
	}

	return s.Server.Serve(listener)
}

func (s *Server) GracefulStop() {
	s.Server.GracefulStop()
}
