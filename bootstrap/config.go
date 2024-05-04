package bootstrap

import (
	"fmt"
	"log"

	"github.com/digisata/auth-service/pkg/grpcserver"
	"github.com/digisata/auth-service/pkg/jwtio"
	"github.com/digisata/auth-service/pkg/mongo"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv         string            `mapstructure:"APP_ENV"`
	ServerAddress  string            `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout int               `mapstructure:"CONTEXT_TIMEOUT"`
	Jwt            jwtio.Config      `mapstructure:"JWT"`
	Mongo          mongo.Config      `mapstructure:"MONGO"`
	GrpcServer     grpcserver.Config `mapstructure:"GRPC_SERVER"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("can't find the file .env : %v", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("environment can't be loaded: %v", err)
	}

	if cfg.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &cfg, nil
}
