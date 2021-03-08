package temporal

import (
	"context"
	"fmt"
	"time"

	"github.com/ContextLogic/autobots/pkg/config"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

type (
	Temporal struct {
		DefaultClients        map[string]DefaultClients
		WorkflowServiceClient workflowservice.WorkflowServiceClient
		NamespaceClient       client.NamespaceClient
	}

	DefaultClients struct {
		Client client.Client
		Worker worker.Worker
	}
)

func New(config *config.TemporalConfig) (t *Temporal, err error) {
	t = &Temporal{
		DefaultClients: make(map[string]DefaultClients),
	}

	for k, v := range config.Clients {
		c, err := client.NewClient(
			client.Options{HostPort: config.HostPort, Namespace: k},
		)
		if err != nil {
			return nil, err
		}
		w := worker.New(c, fmt.Sprintf("%s_%s", config.TaskQueuePrefix, k), worker.Options{
			MaxConcurrentActivityTaskPollers: v.Worker.MaxConcurrentActivityTaskPollers,
		})
		t.DefaultClients[k] = DefaultClients{c, w}
	}

	conn, err := grpc.Dial(config.HostPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	t.WorkflowServiceClient = workflowservice.NewWorkflowServiceClient(conn)

	t.NamespaceClient, err = client.NewNamespaceClient(
		client.Options{HostPort: config.HostPort},
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Temporal) RegisterNamespace(config *config.TemporalConfig) error {
	for k, v := range config.Clients {
		retention := time.Duration(v.Retention) * time.Hour * 24
		err := t.NamespaceClient.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
			Namespace:                        k,
			WorkflowExecutionRetentionPeriod: &retention,
		})
		if err != nil && err.Error() != "Namespace already exists." {
			return err
		}
	}
	return nil
}
