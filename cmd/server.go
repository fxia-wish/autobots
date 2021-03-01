package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	log "github.com/ContextLogic/wish-logger/pkg"

	"github.com/ContextLogic/autobots/config"

	m "github.com/ContextLogic/wish-metric/pkg"

	"github.com/ContextLogic/autobots/pkg/service"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/ContextLogic/wish-sentry-go/pkg/wishsentry"
	"time"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "start a autobots server",
		Long:  `start a autobots server`,
		Run: func(cmd *cobra.Command, args []string) {
			execServerCmd(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

func execServerCmd(cmd *cobra.Command) {
	cfg, err := config.UnmarshalConfig()
	if err != nil {
		fmt.Printf("Failed to unmarshal config file: %s\n", err)
		panic(err)
	}
	// initialize a global metric
	m.InitGlobalMetrics(cfg.BaseConfig.ServiceName)


	// enable histogram metrics for latencies.
	// for more documentations of server side metrics, please refer to
	// https://godoc.org/github.com/grpc-ecosystem/go-grpc-prometheus
	grpc_prometheus.EnableHandlingTimeHistogram()

	// initialize wish-sentry client to pick up the sentry config
	// for more documentations please refer to https://github.com/ContextLogic/wish-sentry-go
	wishsentry.Init(&cfg.SentryConfig)
	defer wishsentry.Flush(5 * time.Second)

	// create a hello service
	s, err := service.NewHelloService(cfg)
	if err != nil {
		panic(err)
	}

	// start the hello service server
	err = s.Base.Start()
	if err != nil {
		log.Fatal(err)
	}
}
