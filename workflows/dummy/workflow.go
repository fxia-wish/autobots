package dummy

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"

	"github.com/ContextLogic/autobots/clients"
	"github.com/ContextLogic/autobots/models"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	DummyWorkflow struct {
		Clients    *clients.Clients
		Activities *DummyActivities
	}
	DummyActivities struct {
		Clients *clients.Clients
		Root    string
	}
)

func NewDummyWorkflow(clients *clients.Clients) *DummyWorkflow {
	_, filename, _, _ := runtime.Caller(0)
	return &DummyWorkflow{
		Clients: clients,
		Activities: &DummyActivities{
			Clients: clients,
			Root:    path.Join(path.Dir(filename), "../.."),
		},
	}
}

func (a *DummyActivities) ReadProfile(path string) (map[string]string, error) {
	profile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer profile.Close()

	value, _ := ioutil.ReadAll(profile)
	results := make(map[string]string)
	json.Unmarshal([]byte(value), &results)
	return results, nil
}

func (a *DummyActivities) CreateOrder(ctx context.Context, order models.Order) error {
	profile, err := a.ReadProfile(path.Join(a.Root, "flags/order.json"))
	if err != nil {
		return err
	}

	a.Clients.Logger.WithField("Order", order).Info("==========calling order service==========")
	if profile["status"] == "valid" {
		time.Sleep(time.Second * 2)
		a.Clients.Logger.WithField("Order", order).Info("==========order is created==========")
		return nil
	}
	a.Clients.Logger.WithField("Order", order).Info("==========order is failed to create==========")
	return errors.New("failed to create order")
}

func (a *DummyActivities) ApprovePayment(ctx context.Context, order models.Order) (*models.OrderResponse, error) {
	profile, err := a.ReadProfile(path.Join(a.Root, "flags/payment.json"))
	if err != nil {
		return nil, err
	}

	a.Clients.Logger.WithField("Order", order).Info("==========calling payment service==========")
	switch profile["status"] {
	case "valid":
		time.Sleep(time.Second * 5)
		a.Clients.Logger.WithField("Order", order).Info("==========payment is approved==========")
		return &models.OrderResponse{
			Order:  &order,
			Status: "succeeded",
		}, nil
	case "invalid":
		a.Clients.Logger.WithField("Order", order).Info("==========payment is declined: invalid==========")
		return &models.OrderResponse{
			Order:  &order,
			Status: "invalid_payment",
		}, nil
	default:
		a.Clients.Logger.WithField("Order", order).Info("==========payment is declined: unknown==========")
		return nil, errors.New("failed to process payment")
	}
}

func (a *DummyActivities) Shipping(ctx context.Context, order models.Order) (*models.OrderResponse, error) {
	profile, err := a.ReadProfile(path.Join(a.Root, "flags/shipping.json"))
	if err != nil {
		return nil, err
	}

	a.Clients.Logger.WithField("Order", order).Info("==========calling shipping service==========")
	switch profile["status"] {
	case "valid":
		time.Sleep(time.Second * 5)
		a.Clients.Logger.WithField("Order", order).Info("==========shipping is initiated==========")
		return &models.OrderResponse{
			Order:  &order,
			Status: "succeeded",
		}, activity.ErrResultPending
	case "invalid":
		a.Clients.Logger.WithField("Order", order).Info("==========shipping is failed: invalid_address==========")
		return &models.OrderResponse{
			Order:  &order,
			Status: "invalid_address",
		}, nil
	default:
		a.Clients.Logger.WithField("Order", order).Info("==========shipping is failed: unknown==========")
		return nil, errors.New("failed to ship")
	}
}

func (a *DummyActivities) Shipped(ctx context.Context, order models.Order) (*models.OrderResponse, error) {
	a.Clients.Logger.WithField("Order", order).Info("==========shipped==========")
	return &models.OrderResponse{
		Order:  &order,
		Status: "succeeded",
	}, nil
}

func (a *DummyActivities) DeclineOrder(ctx context.Context, order models.Order) (*models.OrderResponse, error) {
	a.Clients.Logger.WithField("Order", order).Info("==========calling order declining service==========")
	time.Sleep(time.Second * 5)
	a.Clients.Logger.WithField("Order", order).Info("==========order is declined==========")
	return &models.OrderResponse{
		Order:  &order,
		Status: "order is declined",
	}, nil
}

func (a *DummyActivities) RefundOrder(ctx context.Context, order models.Order) (*models.OrderResponse, error) {
	a.Clients.Logger.WithField("Order", order).Info("==========calling refund service==========")
	time.Sleep(time.Second * 5)
	a.Clients.Logger.WithField("Order", order).Info("==========refund is initiated==========")
	return &models.OrderResponse{
		Order:  &order,
		Status: "order is refunded",
	}, nil
}

func (w *DummyWorkflow) Register() error {
	err := w.Clients.Temporal.RegisterNamespace(GetNamespace())
	if err != nil && err.Error() != "Namespace already exists." {
		return err
	}

	w.Clients.Temporal.Worker.RegisterWorkflow(w.Entry)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.CreateOrder)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.ApprovePayment)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.Shipping)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.Shipped)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.DeclineOrder)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.RefundOrder)

	return nil
}

func (w *DummyWorkflow) Entry(ctx workflow.Context, input interface{}) (interface{}, error) {
	data, _ := json.Marshal(input)
	order := models.Order{}
	err := json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}

	response := models.OrderResponse{}
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second,
			MaximumAttempts:    10,
		},
	})
	if workflow.ExecuteActivity(ctx, w.Activities.CreateOrder, order).Get(ctx, nil); err != nil {
		return nil, err
	}
	err = workflow.ExecuteActivity(ctx, w.Activities.ApprovePayment, order).Get(ctx, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "succeeded" {
		workflow.ExecuteActivity(ctx, w.Activities.DeclineOrder, order).Get(ctx, nil)
		return response, nil
	}
	err = workflow.ExecuteActivity(ctx, w.Activities.Shipping, order).Get(ctx, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "succeeded" {
		workflow.ExecuteActivity(ctx, w.Activities.DeclineOrder, order).Get(ctx, nil)
		workflow.ExecuteActivity(ctx, w.Activities.RefundOrder, order).Get(ctx, &response)
	} else {
		workflow.ExecuteActivity(ctx, w.Activities.Shipped, order).Get(ctx, &response)
	}
	return response, nil
}

func GetNamespace() string {
	return path.Base(reflect.TypeOf(DummyWorkflow{}).PkgPath())
}
