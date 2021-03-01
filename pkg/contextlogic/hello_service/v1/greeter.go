package v1

import (
	"context"
	"sync"

	"github.com/ContextLogic/hello-service/api/proto_gen/contextlogic/hello_service/v1"
	m "github.com/ContextLogic/wish-metric/pkg"
	"github.com/pkg/errors"
)

var (
	greets = make(map[int64]*v1.Greet)
	mutex  sync.Mutex
)

// Greeter is Greeter
type Greeter struct {
}

//CreateGreet is the code implementation of the function defined in the proto files
func (g *Greeter) CreateGreet(ctx context.Context, req *v1.CreateGreetRequest) (res *v1.CreateGreetResponse, err error) {
	mutex.Lock()
	greets[req.Greet.Id] = req.Greet
	mutex.Unlock()
	return &v1.CreateGreetResponse{
		GreetID: req.Greet.Id,
	}, nil
}

//ReadGreet is the code implementation of the function defined in the proto files
func (g *Greeter) ReadGreet(ctx context.Context, req *v1.ReadGreetRequest) (res *v1.ReadGreetResponse, err error) {
	mutex.Lock()
	greet, ok := greets[req.GreetID]
	mutex.Unlock()
	if !ok {
		return nil, errors.Errorf("requested Greet not found given Id: %d", req.GreetID)
	}

	return &v1.ReadGreetResponse{
		Greet: greet,
	}, nil

}

//UpdateGreet is the code implementation of the function defined in the proto files
func (g *Greeter) UpdateGreet(ctx context.Context, req *v1.UpdateGreetRequest) (res *v1.UpdateGreetResponse, err error) {
	mutex.Lock()
	greets[req.Greet.Id] = req.Greet
	mutex.Unlock()
	return &v1.UpdateGreetResponse{
		GreetID: req.Greet.Id,
	}, nil

}

//DeleteGreet is the code implementation of the function defined in the proto files
func (g *Greeter) DeleteGreet(ctx context.Context, req *v1.DeleteGreetRequest) (res *v1.DeleteGreetResponse, err error) {
	mutex.Lock()
	delete(greets, req.GreetID)
	mutex.Unlock()
	return &v1.DeleteGreetResponse{
		GreetID: req.GreetID,
	}, nil

}

//ReadAllGreets is the code implementation of the function defined in the proto files
func (g *Greeter) ReadAllGreets(ctx context.Context, req *v1.ReadAllGreetsRequest) (res *v1.ReadAllGreetsResponse, err error) {

	// example of using a counter from the global metric
	m.Counter("greets_count", "method").WithLabelValues("read_all").Inc()

	mutex.Lock()
	var greetList []*v1.Greet
	for _, g := range greets {
		greetList = append(greetList, g)
	}
	mutex.Unlock()
	return &v1.ReadAllGreetsResponse{
		Greets: &v1.Greets{
			Greets: greetList,
		},
	}, nil
}
