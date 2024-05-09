package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/digisata/auth-service/api/controller"
	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/gateway"
	"github.com/digisata/auth-service/pkg/grpcclient"
	"github.com/digisata/auth-service/pkg/grpcserver"
	"github.com/digisata/auth-service/pkg/interceptors"
	"github.com/digisata/auth-service/pkg/jwtio"
	memcachedRepo "github.com/digisata/auth-service/repository/memcached"
	mongoRepo "github.com/digisata/auth-service/repository/mongo"
	authPb "github.com/digisata/auth-service/stubs/auth"
	"github.com/digisata/auth-service/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Setup app
	app, err := bootstrap.App()
	if err != nil {
		panic(err)
	}

	cfg := app.Cfg

	jwt := jwtio.NewJSONWebToken(&cfg.Jwt, app.MemcachedDB)

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	sugar := logger.Sugar()

	// Dependencies injection
	db := app.Mongo.Database(cfg.Mongo.DBName)
	defer app.CloseDBConnection()

	userRepository := mongoRepo.NewUserRepository(db, domain.USER_COLLECTION)
	profileRepository := mongoRepo.NewProfileRepository(db, domain.USER_COLLECTION)
	cacheRepository := memcachedRepo.NewCacheRepository(app.MemcachedDB)
	timeout := time.Duration(cfg.ContextTimeout) * time.Second
	authController := &controller.AuthController{
		UserUsecase:    usecase.NewUserUsecase(jwt, cfg, userRepository, cacheRepository, timeout),
		ProfileUsecase: usecase.NewProfileUsecase(jwt, cfg, profileRepository, cacheRepository, timeout),
	}

	// Setup GRPC server
	im := interceptors.NewInterceptorManager(jwt, sugar)
	grpcServer, err := grpcserver.NewGrpcServer(cfg.GrpcServer, im, sugar)
	if err != nil {
		panic(err)
	}

	authPb.RegisterAuthServiceServer(grpcServer, authController)
	grpc_health_v1.RegisterHealthServer(grpcServer.Server, health.NewServer())

	err = grpcServer.Run()
	if err != nil {
		panic(err)
	}
	defer grpcServer.Stop(ctx)

	// Setup GRPC client
	grpcClientConn, err := grpcclient.NewGrpcClient(ctx, cfg.GrpcServer, im, grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer grpcClientConn.Close()

	// Setup gateway mux
	gatewayServer := gateway.NewGateway(cfg.Port)
	err = authPb.RegisterAuthServiceHandler(ctx, gatewayServer.ServeMux, grpcClientConn)
	if err != nil {
		panic(err)
	}

	err = gatewayServer.Run(ctx, cfg)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}
