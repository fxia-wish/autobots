package temporal

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ContextLogic/autobots/pkg/auth"
	"github.com/ContextLogic/autobots/pkg/config"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	//Temporal client
	Temporal struct {
		DefaultClients        map[string]DefaultClients
		WorkflowServiceClient workflowservice.WorkflowServiceClient
		NamespaceClient       client.NamespaceClient
	}
	// DefaultClients contains workflow client and worker client
	DefaultClients struct {
		Client client.Client
		Worker worker.Worker
	}
)

// New create new temporal client from config
func New(config *config.TemporalConfig) (t *Temporal, err error) {
	t = &Temporal{
		DefaultClients: make(map[string]DefaultClients),
	}

	/*
		tlsCACertFile := "/Users/bgao/Downloads/DigiCertCA.crt"
		var rpool *x509.CertPool
		pemBytes, err := ioutil.ReadFile(tlsCACertFile)
		if err != nil {
			return nil, err
		}

		rpool = x509.NewCertPool()
		rpool.AppendCertsFromPEM(pemBytes)
	*/
	connOption := client.ConnectionOptions{
		TLS: &tls.Config{
			ServerName: "internode.dev.temporal.i.wish.com",
			//RootCAs:    rpool,
		},
		DisableHealthCheck: true,
	}

	headerProvider := &auth.Provider{}

	for k, v := range config.Clients {
		c, err := client.NewClient(
			client.Options{
				HostPort:          config.HostPort,
				Namespace:         k,
				ConnectionOptions: connOption,
				HeadersProvider:   headerProvider,
			},
		)
		if err != nil {
			return nil, err
		}
		w := worker.New(c, fmt.Sprintf("%s_%s", config.TaskQueuePrefix, k), worker.Options{
			MaxConcurrentActivityTaskPollers: v.Worker.MaxConcurrentActivityTaskPollers,
		})
		t.DefaultClients[k] = DefaultClients{c, w}
	}

	conn, err := grpc.Dial(config.HostPort, grpc.WithTransportCredentials(credentials.NewTLS(connOption.TLS)))
	if err != nil {
		return nil, err
	}
	t.WorkflowServiceClient = workflowservice.NewWorkflowServiceClient(conn)

	t.NamespaceClient, err = client.NewNamespaceClient(
		client.Options{
			HostPort:          config.HostPort,
			ConnectionOptions: connOption,
			HeadersProvider:   headerProvider,
		},
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// RegisterNamespace register namespace
func (t *Temporal) RegisterNamespace(namespace string, retention int) error {
	r := time.Duration(retention) * time.Hour * 24
	err := t.NamespaceClient.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        namespace,
		WorkflowExecutionRetentionPeriod: &r,
	})
	if err != nil && err.Error() != "Namespace already exists." {
		return err
	}
	return nil
}
