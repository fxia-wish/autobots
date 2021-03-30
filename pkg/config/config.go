package config

import (
	"errors"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	Dev       Env = "dev"
	Ec2       Env = "ec2"
	Local     Env = "local"
	ConfigEnv     = "CONFIG_ENVIRONMENT"
)

type (
	Env    string
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
		Logger       *LoggerConfig
		Temporal     *TemporalConfig
		WishFrontend *WishFrontendConfig
	}

	TemporalConfig struct {
		TaskQueue       string
		TaskQueuePrefix string
		HostPort        string
		Clients         map[string]*TemporalClientConfig
	}

	TemporalClientConfig struct {
		Activities *ActivitiesConfig
		Retention  int
		Worker     *WorkerConfig
	}

	ActivitiesConfig struct {
		StartToCloseTimeout int
		RetryPolicy         *RetryPolicyConfig
	}

	RetryPolicyConfig struct {
		InitialInterval    int
		BackoffCoefficient float64
		MaximumInterval    int
		MaximumAttempts    int32
	}

	WorkerConfig struct {
		MaxConcurrentActivityTaskPollers int
	}

	LoggerConfig struct {
		Level string
	}

	WishFrontendConfig struct {
		Host    string
		Timeout int
	}
)

func Init(env ...Env) (*Config, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	_, filename, _, _ := runtime.Caller(0)
	config := &Config{Root: path.Join(path.Dir(filename), "../..")}

	viper.AddConfigPath(path.Join(config.Root, "pkg/config/yaml"))
	viper.SetConfigName("base")
	logger.Infof("merge base config")
	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	logger.WithField("env", env).Info("enviroment info")
	for _, e := range env {
		switch e {
		case Dev:
		case Ec2:
		case Local:
		default:
			logger.Errorf("unsupported config env type")
			return nil, errors.New("unsupported config env type")
		}
		logger.Infof("merge config: %v", e)
		viper.SetConfigName(string(e))
		err := viper.MergeInConfig()
		if err != nil {
			return nil, err
		}
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	logger.WithField("config", config).Info("config info")

	return config, nil
}

func GetEnvironment() Env {
	if os.Getenv(ConfigEnv) == "" {
		os.Setenv(ConfigEnv, "local")
	}

	switch os.Getenv(ConfigEnv) {
	case "dev":
		return Dev
	case "ec2":
		return Ec2
	case "local":
		return Local
	default:
		panic("invalid env")
	}
}
