package workflows

import (
	"github.com/ContextLogic/autobots/pkg/clients"
	"github.com/ContextLogic/autobots/pkg/config"
	dummy "github.com/ContextLogic/autobots/pkg/workflows/dummy"
	"github.com/ContextLogic/autobots/pkg/workflows/wishcashpayment"
)

type (
	// Workflows contains workflow map
	Workflows map[string]Workflow
	// workflow interface definition
	Workflow interface {
		Register() error
	}
)

// New init workflow map
func New(config *config.Config, clients *clients.Clients) Workflows {
	return map[string]Workflow{
		dummy.GetNamespace(): dummy.NewDummyWorkflow(
			config.Clients.Temporal.Clients[dummy.GetNamespace()],
			clients,
		),
		wishcashpayment.GetNamespace(): wishcashpayment.NewWishCashPaymentWorkflow(
			config.Clients.Temporal.Clients[wishcashpayment.GetNamespace()],
			clients,
		),
	}
}

// Register workflow
func (w Workflows) Register() error {
	for _, wf := range w {
		err := wf.Register()
		if err != nil {
			return err
		}
	}
	return nil
}
