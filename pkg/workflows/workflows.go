package workflows

import (
	"github.com/ContextLogic/autobots/pkg/clients"
	"github.com/ContextLogic/autobots/pkg/config"
	dummy "github.com/ContextLogic/autobots/pkg/workflows/dummy"
	"github.com/ContextLogic/autobots/pkg/workflows/wish_cash_payment"
)

type (
	// workflow map
	Workflows map[string]Workflow
	// workflow interface
	Workflow interface {
		Register() error
	}
)

// init workflow map
func New(config *config.Config, clients *clients.Clients) Workflows {
	return map[string]Workflow{
		dummy.GetNamespace(): dummy.NewDummyWorkflow(
			config.Clients.Temporal.Clients[dummy.GetNamespace()],
			clients,
		),
		wish_cash_payment.GetNamespace(): wish_cash_payment.NewWishCashPaymentWorkflow(
			config.Clients.Temporal.Clients[wish_cash_payment.GetNamespace()],
			clients,
		),
	}
}

// register workflow
func (w Workflows) Register() error {
	for _, wf := range w {
		err := wf.Register()
		if err != nil {
			return err
		}
	}
	return nil
}
