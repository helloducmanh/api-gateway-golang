package config

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var logger = logging.MustGetLogger("example")

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config *Config, err error) {
	if path == "" {
		path = "./app.env"
	}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error to reading config file, %s", err)
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Error("error to decode, %v", err)
		return nil, err
	}

	return config, nil
}
