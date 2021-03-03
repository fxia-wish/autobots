package clients

import (
	"github.com/ContextLogic/autobots/clients/logger"
	"github.com/ContextLogic/autobots/clients/temporal"
	"github.com/ContextLogic/autobots/config"
	"github.com/sirupsen/logrus"
)

// Clients defines the clients object
type Clients struct {
	Temporal *temporal.Temporal
	Logger   *logrus.Logger
}

// Init initiates the clients object
func Init(config *config.Config) (*Clients, error) {
	logger, err := logger.New(config.Clients.Logger)
	if err != nil {
		return nil, err
	}

	temporal, err := temporal.New(config.Clients.Temporal)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Temporal: temporal,
		Logger:   logger,
	}, nil
}
