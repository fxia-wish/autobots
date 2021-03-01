package cmd

import (
	"fmt"
	"net/http"

	"github.com/ContextLogic/go-base-service/pkg/client"
	"github.com/ContextLogic/go-base-service/pkg/service"
	"github.com/ContextLogic/hello-service/config"
	"github.com/ContextLogic/hello-service/api/proto_gen/contextlogic/hello_service/v1"
	"github.com/spf13/cobra"
)

var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "start a hello-service client",
		Long:  `start a hello-service client`,
		Run: func(cmd *cobra.Command, args []string) {
			execClientCmd(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(clientCmd)
}

func execClientCmd(cmd *cobra.Command) {
	cfg, err := config.UnmarshalConfig()
	if err != nil {
		fmt.Printf("Failed to unmarshal config file: %s \n", err)
		panic(err)
	}
	exampleClient(cfg)

}

// This is very basic example of a client for the purpose of the client cmd example only. Please place any application specific code in the pkg folder.
func exampleClient(cfg *config.Config) {
	svc, err := service.NewService(&cfg.BaseConfig)
	if err != nil {
		panic(err)
	}

	cc, ctx, err := client.Dial("hello-service-dev", client.WithClientConfig(&cfg.BaseConfig.ClientConfig))
	if err != nil {
		panic(err)
	}

	greeterClient := v1.NewGreeterClient(cc)

	svc.Mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		res, err := greeterClient.ReadAllGreets(ctx, &v1.ReadAllGreetsRequest{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		if res != nil {
			w.Write([]byte(fmt.Sprintf("%v", res.Greets.Greets)))
		} else {
			w.Write([]byte("empty"))
		}
	})

	err = svc.Start()
	if err != nil {
		panic(err)
	}
}
