package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var (
	//Version is the service version
	Version = "unset"
	//Git is the git
	Git     = "unset"
	rootCmd = &cobra.Command{
		Use:   "autobots",
		Short: "This is a template service",
		Long:  `This is a go template service based on go-base-service"`,
	}
)

//Execute is Execute
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Int("grpc_port", 8081, "Port to run grpc server on")
	viper.BindPFlag("base_config.server_config.grpc_port", rootCmd.PersistentFlags().Lookup("grpc_port"))

	rootCmd.PersistentFlags().Int("http_port", 8080, "Port to run http server on")
	viper.BindPFlag("base_config.http_port", rootCmd.PersistentFlags().Lookup("http_port"))

	rootCmd.PersistentFlags().StringP("log_level", "l", "debug", "log level")
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log_level"))

	rootCmd.PersistentFlags().StringP("config_file", "c", "config/service.json", "config file path")
	viper.BindPFlag("config_file", rootCmd.PersistentFlags().Lookup("config_file"))

	rootCmd.PersistentFlags().String("sentry_dsn", "", "Your project's dsn from sentry")
	viper.BindPFlag("sentry_config.sentry_dsn", rootCmd.PersistentFlags().Lookup("sentry_dsn"))

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(viper.GetString("config_file"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configuration file: ", viper.ConfigFileUsed())
	} else {
		panic(err)
	}
}
