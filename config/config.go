package config

import (
	"errors"
	"os"
	"strings"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("example")

type Config struct {
	Servers []string
}

// func LoadConfig(path string) (config *Config, err error) {
// 	if path == "" {
// 		path = "./app.env"
// 	}

// 	viper.SetConfigFile(path)
// 	if err := viper.ReadInConfig(); err != nil {
// 		logger.Error("Error to reading config file, %s", err)
// 		return nil, err
// 	}

// 	err = viper.Unmarshal(&config)
// 	if err != nil {
// 		logger.Error("error to decode, %v", err)
// 		return nil, err
// 	}

// 	return config, nil
// }

func LoadConfigDockerfile() (config *Config, err error) {
	prefix := os.Getenv("PREFIX_SERVICE")

	if prefix == "" {
		logger.Debug("PREFIX_SERVICE is not set")
		return nil, errors.New("PREFIX_SERVICE is not set")
	}

	logger.Debug("Prefix server is :", prefix)
	envs := os.Environ()

	config = &Config{}

	for _, env := range envs {
		// Split biến môi trường thành key và value
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		logger.Debug("value", value)
		logger.Debug("key", key)

		// Kiểm tra nếu key bắt đầu với tiền tố
		if strings.HasPrefix(key, prefix) {
			logger.Debug("key", key)
			config.Servers = append(config.Servers, value)
		}
	}

	for index, service := range config.Servers {
		logger.Debug("Service ", index, service)
	}

	return config, err
}
