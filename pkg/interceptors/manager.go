package interceptors

import (
	"context"

	"github.com/digisata/auth-service/pkg/constants"
	"github.com/digisata/auth-service/pkg/jwtio"
	"github.com/digisata/auth-service/pkg/tracing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type InterceptorManager interface {
	Logger(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error)
	ClientRequestLoggerInterceptor() func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error
	AuthInterceptor(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error)
}

// InterceptorManager struct
type interceptorManager struct {
	logger            *zap.SugaredLogger
	jwtManager        *jwtio.JSONWebToken
	restrictedMethods map[string][]string
}

// NewInterceptorManager InterceptorManager constructor
func NewInterceptorManager(jwtManager *jwtio.JSONWebToken, logger *zap.SugaredLogger) *interceptorManager {
	return &interceptorManager{
		logger:            logger,
		jwtManager:        jwtManager,
		restrictedMethods: restrictedMethods(),
	}
}

// Logger Interceptor
func (im interceptorManager) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	reply, err := handler(ctx, req)
	if err != nil {
		im.logger.Errorw(constants.ERROR,
			"method", info.FullMethod,
			"request", req,
			"error", err.Error(),
		)

		return reply, err
	}

	im.logger.Infow(constants.INFO,
		"method", info.FullMethod,
		"request", req,
		"error", nil,
	)

	return reply, err
}

// ClientRequestLoggerInterceptor gRPC client interceptor
func (im interceptorManager) ClientRequestLoggerInterceptor() func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		im.logger.Infow(constants.INFO,
			"method", method,
			"request", req,
			"error", nil,
		)

		return err
	}
}

func (im interceptorManager) AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "Interceptors.AuthInterceptor")
	defer span.End()

	if _, ok := im.isRestricted(ctx, info.FullMethod); !ok {
		return handler(ctx, req)
	}

	claims, err := im.jwtManager.Verify(ctx)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, "claims", claims)

	return handler(newCtx, req)
}

func (im interceptorManager) isRestricted(ctx context.Context, method string) ([]string, bool) {
	_, span := tracing.StartGrpcServerTracerSpan(ctx, "Interceptors.isRestricted")
	defer span.End()

	value, ok := im.restrictedMethods[method]

	return value, ok
}

func restrictedMethods() map[string][]string {
	const path string = "/auth_service.user.AuthService/"

	return map[string][]string{
		path + "CreateUser":  {constants.ACCESS_TOKEN},
		path + "GetUserByID": {constants.ACCESS_TOKEN},
	}
}
