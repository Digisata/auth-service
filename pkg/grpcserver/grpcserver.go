package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/pkg/interceptors"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type GrpcServer struct {
	*grpc.Server
	Listener net.Listener
	Port     string
	Network  string
}

const (
	maxConnectionIdle = 300
	gRPCTimeout       = 15
	maxConnectionAge  = 300
	gRPCTime          = 600
)

func NewGrpcServer(cfg *bootstrap.Config, im interceptors.InterceptorManager, opts ...grpc.ServerOption) (*GrpcServer, error) {
	if cfg.GrpcTls {
		certFile := "ssl/certificates/server.crt" // => your certFile file path
		keyFile := "ssl/server.pem"               // => your keyFile file patn

		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(
		opts,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Second,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Second,
			Time:              gRPCTime * time.Second,
		}),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcCtxtags.UnaryServerInterceptor(),
			grpcPrometheus.UnaryServerInterceptor,
			grpcRecovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
			im.AuthInterceptor,
			im.Logger,
		)),
	)

	server := grpc.NewServer(opts...)
	grpcPrometheus.Register(server)

	return &GrpcServer{
		Server:  server,
		Network: cfg.GrpcNetwork,
		Port:    cfg.GrpcPort,
	}, nil
}

func (grpcServer *GrpcServer) Run() error {
	listener, err := net.Listen(grpcServer.Network, fmt.Sprintf(":%v", grpcServer.Port))
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	grpcServer.Listener = listener

	go func() {
		if err := grpcServer.Server.Serve(grpcServer.Listener); err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func (grpcServer *GrpcServer) Stop(ctx context.Context) {
	if err := grpcServer.Listener.Close(); err != nil {
		panic(err)
	}

	go func() {
		defer grpcServer.Server.GracefulStop()
		<-ctx.Done()
	}()
}
