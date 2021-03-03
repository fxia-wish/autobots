package temporal

import (
	"context"
	"time"

	"github.com/ContextLogic/autobots/config"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
)

type (
	Temporal struct {
		Client    client.Client
		Worker    worker.Worker
		Frontend  workflowservice.WorkflowServiceClient
		Namespace client.NamespaceClient
	}
	Options struct {
		ClientOptions client.Options
		WorkerOptions worker.Options
	}
)

func New(config *config.TemporalConfig, options *Options) (t *Temporal, err error) {
	t = &Temporal{}
	t.Client, err = client.NewClient(options.ClientOptions)
	if err != nil {
		return nil, err
	}

	t.Worker = worker.New(t.Client, config.TaskQueue, options.WorkerOptions)

	conn, err := grpc.Dial(
		options.ClientOptions.HostPort,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	t.Frontend = workflowservice.NewWorkflowServiceClient(conn)

	t.Namespace, err = client.NewNamespaceClient(
		client.Options{HostPort: options.ClientOptions.HostPort},
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Temporal) RegisterNamespace(namespace string) error {
	retention := 1 * time.Hour * 24
	err := t.Namespace.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        namespace,
		WorkflowExecutionRetentionPeriod: &retention,
	})
	return err
}
