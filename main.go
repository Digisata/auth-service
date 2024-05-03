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
	"github.com/digisata/auth-service/repository"
	userPb "github.com/digisata/auth-service/stubs/user"
	"github.com/digisata/auth-service/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	app, err := bootstrap.App()
	if err != nil {
		panic(err)
	}

	cfg := app.Cfg

	db := app.Mongo.Database(cfg.DBName)
	defer app.CloseDBConnection()

	jwt := jwtio.NewJSONWebToken(cfg)

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	sugar := logger.Sugar()

	im := interceptors.NewInterceptorManager(jwt, sugar)

	grpcServer, err := grpcserver.NewGrpcServer(cfg, im, sugar)
	if err != nil {
		panic(err)
	}

	ur := repository.NewUserRepository(db, domain.CollectionUser)
	timeout := time.Duration(cfg.ContextTimeout) * time.Second
	uc := &controller.UserController{
		LoginUsecase:        usecase.NewLoginUsecase(ur, timeout),
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(ur, timeout),
		UserUsecase:         usecase.NewUserUsecase(ur, timeout),
	}

	userPb.RegisterAuthServiceServer(grpcServer, uc)
	grpc_health_v1.RegisterHealthServer(grpcServer.Server, health.NewServer())

	err = grpcServer.Run()
	if err != nil {
		panic(err)
	}
	defer grpcServer.Stop(ctx)

	grpcClientConn, err := grpcclient.NewGrpcClient(ctx, cfg, im, grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer grpcClientConn.Close()

	gatewayServer := gateway.NewGateway(cfg.ServerAddress)
	err = userPb.RegisterAuthServiceHandler(ctx, gatewayServer.ServeMux, grpcClientConn)
	if err != nil {
		panic(err)
	}

	err = gatewayServer.Run(ctx, cfg)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}
