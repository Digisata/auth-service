package bootstrap

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	GrpcNetwork            string `mapstructure:"GRPC_NETWORK"`
	GrpcPort               string `mapstructure:"GRPC_PORT"`
	GrpcTls                bool   `mapstructure:"GRPC_TLS"`
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
