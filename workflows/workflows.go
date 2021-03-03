package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/ContextLogic/autobots/clients"
	dummy "github.com/ContextLogic/autobots/workflows/dummy"
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
		dummy.GetNamespace(): dummy.NewDummyWorkflow(clients),
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
