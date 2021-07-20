package authetication

import (
	"context"
	"flag"
	"fmt"

	"github.com/ContextLogic/authn/pkg/authn"
	"github.com/ContextLogic/authn/pkg/common"
)

type WishAuthHeadersProvider interface {
	GetHeaders(ctx context.Context) (map[string]string, error)
}

type AuthProvider struct {
}

func (p *AuthProvider) GetHeaders(ctx context.Context) (map[string]string, error) {
	env := flag.String("env", string(common.EnvLocal), "environment")
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

	token, err := requester.GetToken(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("got token: %v\n", token)
	var ret map[string]string
	ret["Authorization"] = token
	return ret, nil
}
