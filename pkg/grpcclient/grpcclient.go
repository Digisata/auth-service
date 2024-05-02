package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/pkg/interceptors"
	grpcRetry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	backoffLinear  = 100 * time.Millisecond
	backoffRetries = 3
)

func NewGrpcClient(ctx context.Context, cfg *bootstrap.Config, im interceptors.InterceptorManager, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if cfg.GrpcTls {
		certFile := "ssl/certificates/ca.crt" // => file path location your certFile
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		creds := grpc.WithTransportCredentials(insecure.NewCredentials())
		opts = append(opts, creds)
	}

	opts = append(
		opts,
		grpc.WithUnaryInterceptor(im.ClientRequestLoggerInterceptor()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithUnaryInterceptor(grpcRetry.UnaryClientInterceptor([]grpcRetry.CallOption{
			grpcRetry.WithBackoff(grpcRetry.BackoffLinear(backoffLinear)),
			grpcRetry.WithCodes(codes.NotFound, codes.Aborted),
			grpcRetry.WithMax(backoffRetries),
		}...)),
	)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf(":%v", cfg.GrpcPort), opts...)
	if err != nil {
		return nil, errors.Wrap(err, "grpc.DialContext")
	}
	return conn, nil
}
