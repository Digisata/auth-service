package interceptors

import (
	"context"

	"github.com/digisata/auth-service/pkg/constants"
	"github.com/digisata/auth-service/pkg/jwtio"
	"github.com/digisata/auth-service/pkg/tracing"
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
	jwtManager        *jwtio.JSONWebToken
	restrictedMethods map[string][]string
}

// NewInterceptorManager InterceptorManager constructor
func NewInterceptorManager(jwtManager *jwtio.JSONWebToken) *interceptorManager {
	return &interceptorManager{jwtManager: jwtManager, restrictedMethods: restrictedMethods()}
}

// Logger Interceptor
func (im *interceptorManager) Logger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	reply, err := handler(ctx, req)
	if err != nil {
		return reply, err
	}

	return reply, err
}

// ClientRequestLoggerInterceptor gRPC client interceptor
func (im *interceptorManager) ClientRequestLoggerInterceptor() func(
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
		return err
	}
}

func restrictedMethods() map[string][]string {
	const path string = "/auth_service.user.AuthService/"

	return map[string][]string{
		path + "Verify": {constants.RefreshToken},

		path + "CreateCNSToken":            {constants.AccessToken},
		path + "GetAllCNSToken":            {constants.AccessToken},
		path + "GetAllCNSTokenByAccountId": {constants.AccessToken},
		path + "GetCNSToken":               {constants.AccessToken},
		path + "DeleteCNSToken":            {constants.AccessToken},

		path + "GetNotifToken":    {constants.AccessToken},
		path + "ResetNotifToken":  {constants.AccessToken},
		path + "GetAllCNSAccount": {constants.AccessToken},
		path + "GetCNSAccount":    {constants.AccessToken},
		path + "UpdateCNSAccount": {constants.AccessToken},
		path + "DeleteCNSAccount": {constants.AccessToken},

		path + "CreateCNSChannelGroup": {constants.AccessToken},
		path + "GetAllCNSChannelGroup": {constants.AccessToken},
		path + "GetCNSChannelGroup":    {constants.AccessToken},
		path + "UpdateCNSChannelGroup": {constants.AccessToken},
		path + "DeleteCNSChannelGroup": {constants.AccessToken},

		path + "CreateCNSChannelAccountGroup": {constants.AccessToken},
		path + "GetAllCNSChannelAccountGroup": {constants.AccessToken},
		path + "GetCNSChannelAccountGroup":    {constants.AccessToken},
		path + "DeleteCNSChannelAccountGroup": {constants.AccessToken},

		path + "CreateCNSNotificationCategory":            {constants.AccessToken},
		path + "GetAllCNSNotificationCategory":            {constants.AccessToken},
		path + "GetAllCNSNotificationCategoryByAccountId": {constants.AccessToken},
		path + "GetCNSNotificationCategory":               {constants.AccessToken, constants.RefreshToken},
		path + "DeleteCNSNotificationCategory":            {constants.AccessToken},

		path + "CreateCNSNotificationCategoryDetail": {constants.AccessToken},
		path + "GetAllCNSNotificationCategoryDetail": {constants.AccessToken},
		path + "GetCNSNotificationCategoryDetail":    {constants.AccessToken},
		path + "UpdateCNSNotificationCategoryDetail": {constants.AccessToken},
		path + "DeleteCNSNotificationCategoryDetail": {constants.AccessToken},

		path + "CreateCNSNotificationType":                         {constants.AccessToken},
		path + "GetAllCNSNotificationType":                         {constants.AccessToken},
		path + "GetAllCNSNotificationTypeByNotificationCategoryId": {constants.AccessToken},
		path + "GetCNSNotificationType":                            {constants.AccessToken, constants.RefreshToken},
		path + "DeleteCNSNotificationType":                         {constants.AccessToken},
		path + "GetAllCNSNotificationTypeByCategoryID":             {constants.AccessToken},

		path + "CreateCNSNotificationTypeDetail": {constants.AccessToken},
		path + "GetAllCNSNotificationTypeDetail": {constants.AccessToken},
		path + "GetCNSNotificationTypeDetail":    {constants.AccessToken},
		path + "UpdateCNSNotificationTypeDetail": {constants.AccessToken},
		path + "DeleteCNSNotificationTypeDetail": {constants.AccessToken},
	}
}

func (im *interceptorManager) AuthInterceptor(
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

	claims, err := im.jwtManager.Validate(ctx)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, "claims", claims)

	return handler(newCtx, req)
}

func (im *interceptorManager) isRestricted(ctx context.Context, method string) ([]string, bool) {
	_, span := tracing.StartGrpcServerTracerSpan(ctx, "Interceptors.isRestricted")
	defer span.End()

	value, restricted := im.restrictedMethods[method]
	return value, restricted
}
