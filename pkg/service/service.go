package service

import (
	base "github.com/ContextLogic/go-base-service/pkg/service"
	v1 "github.com/ContextLogic/hello-service/api/proto_gen/contextlogic/hello_service/v1"
	"github.com/ContextLogic/hello-service/config"
	svc "github.com/ContextLogic/hello-service/pkg/contextlogic/hello_service/v1"
	"google.golang.org/grpc/reflection"
)

//HelloService is HelloService
type HelloService struct {
	Base *base.Service
}

// NewHelloService creates a new hello service
func NewHelloService(cfg *config.Config) (*HelloService, error) {

	// create a base service
	baseService, err := base.NewService(&cfg.BaseConfig)
	if err != nil {
		return nil, err
	}

	// register the rpc service to the base service
	v1.RegisterGreeterServer(baseService.Server.GRPCServer, &svc.Greeter{})
	v1.RegisterStreamServer(baseService.Server.GRPCServer, &svc.Stream{})

	// enable grpc server reflection to help local development with grpc cli tools
	// e.g. grpcurl
	// NOTE: this is only for development/debugging purpose.
	// You do not need server reflection in production.
	if cfg.AppConfig.EnableReflection {
		reflection.Register(baseService.Server.GRPCServer)
	}

	/*
		An example of how to use grpc_gateway to serve http requests.
		i.e. create a grpc gateway as the reverse proxy to prpxy http requests to rpc
	*/
	//gwMux := runtime.NewServeMux()
	//v1.RegisterGreeterHandlerFromEndpoint(
	//	context.Background(), gwMux,
	//	fmt.Sprintf(":%d", baseService.Config.ServerConfig.GRPCPort),
	//	[]grpc.DialOption{
	//		grpc.WithInsecure(),
	//	},
	//)
	//v1.RegisterStreamHandlerFromEndpoint(
	//	context.Background(), gwMux,
	//	fmt.Sprintf(":%d", baseService.Config.ServerConfig.GRPCPort),
	//	[]grpc.DialOption{
	//		grpc.WithInsecure(),
	//	},
	//)
	//baseService.Mux.PathPrefix("/").Handler(gwMux)

	return &HelloService{baseService}, nil
}
