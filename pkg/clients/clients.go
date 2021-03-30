package clients

import (
	"github.com/ContextLogic/autobots/pkg/clients/logger"
	"github.com/ContextLogic/autobots/pkg/clients/temporal"
	"github.com/ContextLogic/autobots/pkg/clients/wish_frontend"
	"github.com/ContextLogic/autobots/pkg/config"
	"github.com/sirupsen/logrus"
)

// Clients defines the clients object
type Clients struct {
	Temporal     *temporal.Temporal
	WishFrontend *wish_frontend.WishFrontend
	Logger       *logrus.Logger
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
		Temporal:     temporal,
		WishFrontend: wish_frontend.New(config.Clients.WishFrontend),
		Logger:       logger,
	}, nil
}
