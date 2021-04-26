package workflows

import (
	"github.com/ContextLogic/autobots/pkg/clients"
	"github.com/ContextLogic/autobots/pkg/config"
	dummies "github.com/ContextLogic/autobots/pkg/workflows/dummy"
	wishcashpayments "github.com/ContextLogic/autobots/pkg/workflows/wishcashpayment"
)

type (
	// Workflows contains workflow map
	Workflows map[string]Workflow
	// Workflow interface definition
	Workflow interface {
		Register() error
	}
)

// New init workflow map
func New(config *config.Config, clients *clients.Clients) Workflows {
	return map[string]Workflow{
		dummies.GetNamespace(): dummies.NewDummyWorkflow(
			config.Clients.Temporal.Clients[dummies.GetNamespace()],
			clients,
		),
		wishcashpayments.GetNamespace(): wishcashpayments.NewWishCashPaymentWorkflow(
			config.Clients.Temporal.Clients[wishcashpayments.GetNamespace()],
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
