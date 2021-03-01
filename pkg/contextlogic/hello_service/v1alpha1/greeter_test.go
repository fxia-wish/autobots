package v1alpha1

import (
	"context"
	"testing"

	"github.com/ContextLogic/hello-service/api/proto_gen/contextlogic/hello_service/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func testHelloService() (*Greeter, error) {
	return &Greeter{}, nil
}

func TestCreateGreet(t *testing.T) {

	server, err := testHelloService()
	assert.Emptyf(t, err, "Unable to get test server: %v", err)

	greet := v1alpha1.Greet{Id: 1}
	resp, err := server.CreateGreet(context.Background(), &v1alpha1.CreateGreetRequest{
		Greet: &greet,
	})
	assert.Emptyf(t, err, "Error calling CreateGreet: %v", err)
	assert.Equalf(t, greet.Id, resp.GreetID, "CreateGreet response returned an incorrect ID")
}

func TestReadGreet(t *testing.T) {
	server, err := testHelloService()
	assert.Emptyf(t, err, "Unable to get test server: %v", err)

	greet := v1alpha1.Greet{Id: 1}
	server.CreateGreet(context.Background(), &v1alpha1.CreateGreetRequest{
		Greet: &greet,
	})

	resp, err := server.ReadGreet(context.Background(), &v1alpha1.ReadGreetRequest{
		GreetID: 1,
	})
	assert.Emptyf(t, err, "Error calling ReadGreet: %v", err)
	assert.Equal(t, greet.Id, resp.Greet.Id, "ReadGreet response returned an incorrect ID")

}

func TestReadAllGreets(t *testing.T) {
	server, err := testHelloService()
	assert.Emptyf(t, err, "Unable to get test server: %v", err)

	greets := v1alpha1.Greets{
		Greets: []*v1alpha1.Greet{
			{Id: 1},
			{Id: 2},
		},
	}
	for _, g := range greets.Greets {
		server.CreateGreet(context.Background(), &v1alpha1.CreateGreetRequest{
			Greet: g,
		})
	}

	resp, err := server.ReadAllGreets(context.Background(), &v1alpha1.ReadAllGreetsRequest{})
	assert.Emptyf(t, err, "Error calling ReadAllGreets: %v", err)
	assert.ElementsMatch(t, greets.Greets, resp.Greets.Greets, "ReadAllGreets response returned wrong list of greets")
}

func TestUpdateGreet(t *testing.T) {
	server, err := testHelloService()
	assert.Emptyf(t, err, "Unable to get test server: %v", err)

	greetBefore := v1alpha1.Greet{Id: 1, GreetMessage: &v1alpha1.Greet_Message{Msg: "1"}}
	server.CreateGreet(context.Background(), &v1alpha1.CreateGreetRequest{
		Greet: &greetBefore,
	})

	greetAfter := v1alpha1.Greet{Id: 1, GreetMessage: &v1alpha1.Greet_Message{Msg: "2"}}
	resp, err := server.UpdateGreet(context.Background(), &v1alpha1.UpdateGreetRequest{
		GreetID: greetBefore.Id,
		Greet:   &greetAfter,
	})
	assert.Emptyf(t, err, "Error calling UpdateGreet: %v", err)
	assert.Equal(t, greetAfter.Id, resp.GreetID, "UpdateGreet response returned incorrect ID")
}

func TestDeleteGreet(t *testing.T) {
	server, err := testHelloService()
	assert.Emptyf(t, err, "Unable to get test server: %v", err)

	greet := v1alpha1.Greet{Id: 1, GreetMessage: &v1alpha1.Greet_Message{Msg: "1"}}
	server.CreateGreet(context.Background(), &v1alpha1.CreateGreetRequest{
		Greet: &greet,
	})
	resp, err := server.DeleteGreet(context.Background(), &v1alpha1.DeleteGreetRequest{
		GreetID: greet.Id,
	})
	assert.Emptyf(t, err, "Error calling DeleteGreet: %v", err)
	assert.Equal(t, greet.Id, resp.GreetID, "DeleteGreet response returned incorrect ID")
}
