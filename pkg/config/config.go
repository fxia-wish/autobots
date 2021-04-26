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
	// Dev env
	Dev Env = "dev"
	// Ec2 env
	Ec2 Env = "ec2"
	// Local env
	Local Env = "local"
	// ConfigEnv env
	ConfigEnv = "CONFIG_ENVIRONMENT"
)

type (
	// Env string
	Env string
	// Config collection
	Config struct {
		Root    string
		Service *ServiceConfig
		Clients *ClientsConfig
	}
	// ServiceConfig contains service related configuration
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

	// ClientsConfig contains autobot related configuration
	ClientsConfig struct {
		Logger       *LoggerConfig
		Temporal     *TemporalConfig
		WishFrontend *WishFrontendConfig
	}
	//TemporalConfig  contains workflow related configuration
	TemporalConfig struct {
		TaskQueue       string
		TaskQueuePrefix string
		HostPort        string
		Clients         map[string]*TemporalClientConfig
	}
	//TemporalClientConfig contains temporal client related configuration
	TemporalClientConfig struct {
		Activities *ActivitiesConfig
		Retention  int
		Worker     *WorkerConfig
	}
	//ActivitiesConfig contains activity related configuration
	ActivitiesConfig struct {
		StartToCloseTimeout int
		RetryPolicy         *RetryPolicyConfig
	}
	//RetryPolicyConfig contains retry related configuration
	RetryPolicyConfig struct {
		InitialInterval    int
		BackoffCoefficient float64
		MaximumInterval    int
		MaximumAttempts    int32
	}
	//WorkerConfig contains worker related configuration
	WorkerConfig struct {
		MaxConcurrentActivityTaskPollers int
	}
	//LoggerConfig contains logging related configuration
	LoggerConfig struct {
		Level string
	}
	//WishFrontendConfig contains wish fe related configuration
	WishFrontendConfig struct {
		Host    string
		Timeout int
	}
)

// Init config from file
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

// GetEnvironment
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
