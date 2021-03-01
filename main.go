package main

import (
	c "github.com/ContextLogic/go-base-service/pkg/config"
	s "github.com/ContextLogic/go-base-service/pkg/service"
	"github.com/ContextLogic/autobots/clients"
	"github.com/ContextLogic/autobots/config"
	"github.com/ContextLogic/autobots/handlers"
	"github.com/ContextLogic/autobots/workflows"
	"go.temporal.io/sdk/worker"
)

func main() {
	config, err := config.Init()
	if err != nil {
		panic(err)
	}

	clients, err := clients.Init(config)
	if err != nil {
		panic(err)
	}
	defer clients.Temporal.Client.Close()

	s, err := s.NewService(
		&c.Config{
			ServiceName:     config.Service.ServiceName,
			HTTPPort:        config.Service.HTTP.Port,
			ShutDownTimeOut: &config.Service.ShutDownTimeOut,
			ShutDownDelay:   &config.Service.ShutDownDelay,
			ServerConfig: c.ServerConfig{
				GRPCPort: config.Service.GRPC.Port,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	workflows := workflows.New(clients)
	err = workflows.Register()
	if err != nil {
		panic(err)
	}
	go clients.Temporal.Worker.Run(worker.InterruptCh())
	handlers.New(config, clients, s, workflows)
	if err = s.Start(); err != nil {
		panic(err)
	}
}
