package workflows

import (
	"github.com/ContextLogic/autobots/clients"
	"github.com/ContextLogic/autobots/config"
	dummy "github.com/ContextLogic/autobots/workflows/dummy"
	"github.com/ContextLogic/autobots/workflows/wish_cash_payment"
)

type (
	Workflows map[string]Workflow
	Workflow  interface {
		Register() error
	}
)

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

func (w Workflows) Register() error {
	for _, wf := range w {
		err := wf.Register()
		if err != nil {
			return err
		}
	}
	return nil
}
