package config

import (
	"errors"
	"github.com/spf13/viper"
	"websocket-chat-service/init/logger"
	"websocket-chat-service/pkg/constants"
)

var ServerConfig Config

type Config struct {
	ApiDebug bool   `mapstructure:"API_DEBUG"`
	ApiPort  int    `mapstructure:"API_PORT"`
	ApiEntry string `mapstructure:"API_ENTRY"`

	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	ScyllaHosts    []string `mapstructure:"SCYLLA_HOSTS"`
	ScyllaKeyspace string   `mapstructure:"SCYLLA_KEYSPACE"`

	WebsocketURL string `mapstructure:"WEBSOCKET_URL"`
}

func InitConfig() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil {
		logger.Error(err.Error(), constants.ConfigLogger)
		return err
	}

	if err := viper.Unmarshal(&ServerConfig); err != nil {
		logger.Error(err.Error(), constants.ConfigLogger)
		return err
	}

	if err := CheckVars(); err != nil {
		return err
	}

	return nil
}

func CheckVars() error {
	if ServerConfig.ApiPort == 0 || ServerConfig.ApiEntry == "" {
		return errors.New("API port and API entry is required environment variables")
	}

	if ServerConfig.ScyllaKeyspace == "" {
		return errors.New("ScyllaKeyspace is required environment variable")
	}

	return nil
}
