package auth

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/ContextLogic/authn/pkg/authn"
	"github.com/ContextLogic/authn/pkg/common"
	"github.com/ContextLogic/autobots/pkg/config"
)

//WishAuthHeadersProvider interface
type WishAuthHeadersProvider interface {
	//GetHeaders func
	GetHeaders(ctx context.Context) (map[string]string, error)
}

//Provider struct
type Provider struct {
}

//GetHeaders use authn library to get auth token
func (p *Provider) GetHeaders(ctx context.Context) (map[string]string, error) {
	var token string
	var err error
	env := string(config.GetEnvironment())

	if env == "local" {
		t := &authn.K8sIDTokenJSON{
			Issuer:   "iss",
			Audience: "aud",
			Duration: 10 * time.Hour,
			Kid:      "testk8s",
			Subject:  "autobots",
			Groups:   []string{"autobots"},
		}
		token, err = authn.NewTestToken(t.GetTokenMap())
		if err != nil {
			return nil, err
		}
	} else {
		env := flag.String("env", string(env), "environment")
		isTest := flag.Bool("test", false, "a flag for unit test, test token issued if true")
		flag.Parse()

		var (
			requester *authn.TokenRequester
			err       error
		)
		cfg, err := authn.NewKubernetesRequesterConfig(common.Environment(*env), *isTest)
		if err != nil {
			fmt.Printf("failed to create kubernetes requester config: %v\n", err)
		}
		requester, err = authn.NewKubernetesRequester(cfg)
		if err != nil {
			fmt.Printf("failed to construct requester: %v\n", err)
		}

		token, err = requester.GetToken(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("got token: %s, env:%s\n", token, env)
	return map[string]string{"authorization": token, "authorization-extras": env}, nil
}
