package workflows

import (
	"github.com/ContextLogic/autobots/pkg/clients"
	dummy "github.com/ContextLogic/autobots/pkg/workflows/dummy"
	"github.com/ContextLogic/autobots/pkg/workflows/wish_cash_payment"
)

type (
	Workflows map[string]Workflow
	Workflow  interface {
		Register() error
	}
)

func New(clients *clients.Clients) Workflows {
	return map[string]Workflow{
		dummy.GetNamespace():             dummy.NewDummyWorkflow(clients),
		wish_cash_payment.GetNamespace(): wish_cash_payment.NewWishCashPaymentWorkflow(clients),
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
