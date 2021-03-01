package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/ContextLogic/autobots/clients"
	ct "github.com/ContextLogic/autobots/workflows/commerce_transaction"
)

type (
	Workflows map[string]Workflow
	Workflow  interface {
		Entry(workflow.Context, interface{}) (interface{}, error)
		Register() error
	}
)

func New(clients *clients.Clients) Workflows {
	return map[string]Workflow{
		ct.GetNamespace(): ct.NewCommerceTransactionWorkflow(clients),
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
