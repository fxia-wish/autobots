package main

import (
	"github.com/ContextLogic/autobots/pkg/clients"
	"github.com/ContextLogic/autobots/pkg/config"
	"github.com/ContextLogic/autobots/pkg/handlers"
	"github.com/ContextLogic/autobots/pkg/workflows"
	c "github.com/ContextLogic/go-base-service/pkg/config"
	s "github.com/ContextLogic/go-base-service/pkg/service"
	"go.temporal.io/sdk/worker"
)

func main() {
	config, err := config.Init(config.GetEnvironment())
	if err != nil {
		panic(err)
	}

	clients, err := clients.Init(config)
	if err != nil {
		panic(err)
	}

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

	workflows := workflows.New(config, clients)
	err = workflows.Register()
	if err != nil {
		panic(err)
	}
	for k := range workflows {
		go clients.Temporal.DefaultClients[k].Worker.Run(worker.InterruptCh())
	}
	handlers.New(config, clients, s, workflows)
	if err = s.Start(); err != nil {
		panic(err)
	}
	for k := range workflows {
		defer clients.Temporal.DefaultClients[k].Client.Close()
	}
}
