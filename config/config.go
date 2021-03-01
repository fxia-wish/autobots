package config

import (
	"path"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Root    string
		Service *ServiceConfig
		Clients *ClientsConfig
	}

	ServiceConfig struct {
		ServiceName     string
		ShutDownTimeOut time.Duration
		ShutDownDelay   time.Duration

		HTTP struct {
			Port int
		}
		GRPC struct {
			Port int
		}
	}

	ClientsConfig struct {
		Logger   *LoggerConfig
		Temporal *TemporalConfig
	}

	TemporalConfig struct {
		TaskQueue string
		HostPort  string
	}

	LoggerConfig struct {
		Level string
	}
)

func Init() (*Config, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	_, filename, _, _ := runtime.Caller(0)
	config := &Config{Root: path.Join(path.Dir(filename), "../..")}

	viper.AddConfigPath(path.Join(config.Root, "autobots/config/yaml"))
	viper.SetConfigName("base")
	logger.Infof("merge base config")
	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	logger.WithField("config", config).Info("config info")

	return config, nil
}
