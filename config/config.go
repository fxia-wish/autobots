package config

import (
	base_config "github.com/ContextLogic/go-base-service/pkg/config"
	"github.com/ContextLogic/wish-sentry-go/pkg/wishsentry"
	"github.com/spf13/viper"
)

// AppConfig holds the configuration of the application
type AppConfig struct {
	EnableReflection bool `mapstructure:"enable_reflection"`
}

//Config is the top level service config
type Config struct {
	BaseConfig base_config.Config `mapstructure:"base_config"`

	SentryConfig wishsentry.Config `mapstructure:"sentry_config"`

	//Add you app's config here
	AppConfig `mapstructure:"app_config"`
}

//UnmarshalConfig unmarshals config fromo viper to Config Struct
func UnmarshalConfig() (*Config, error) {
	// unmarshal to server config
	var config Config
	err := viper.Unmarshal(&config)
	return &config, err
}
