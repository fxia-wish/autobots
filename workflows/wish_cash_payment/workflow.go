package wish_cash_payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"time"

	"github.com/ContextLogic/autobots/clients"
	"github.com/ContextLogic/autobots/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	WishCashPaymentWorkflow struct {
		Clients    *clients.Clients
		Activities *WishCashPaymentActivities
	}
	WishCashPaymentActivities struct {
		Clients *clients.Clients
		Root    string
	}
)

func NewWishCashPaymentWorkflow(clients *clients.Clients) *WishCashPaymentWorkflow {
	_, filename, _, _ := runtime.Caller(0)
	return &WishCashPaymentWorkflow{
		Clients: clients,
		Activities: &WishCashPaymentActivities{
			Clients: clients,
			Root:    path.Join(path.Dir(filename), "../.."),
		},
	}
}

func (a *WishCashPaymentActivities) WishCashPaymentCreateOrder(ctx context.Context, h http.Header, body []byte) (*models.WishCashPaymentCreateOrderResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "http://lshu.corp.contextlogic.com/api/temporal-payment/create-order", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = h
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &models.WishCashPaymentCreateOrderResponse{}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (a *WishCashPaymentActivities) WishCashPaymentApprovePayment(ctx context.Context, h http.Header, body []byte) (*models.WishCashPaymentApprovePaymentResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "http://lshu.corp.contextlogic.com/api/temporal-payment/approve-payment", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = h
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &models.WishCashPaymentApprovePaymentResponse{}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (w *WishCashPaymentWorkflow) Register() error {
	err := w.Clients.Temporal.RegisterNamespace(GetNamespace())
	if err != nil && err.Error() != "Namespace already exists." {
		return err
	}

	w.Clients.Temporal.Worker.RegisterWorkflow(w.WishCashPaymentWorkflow)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.WishCashPaymentCreateOrder)
	w.Clients.Temporal.Worker.RegisterActivity(w.Activities.WishCashPaymentApprovePayment)

	return nil
}

func (w *WishCashPaymentWorkflow) WishCashPaymentWorkflow(ctx workflow.Context, h http.Header, body []byte) (interface{}, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.0,
			MaximumInterval:    time.Second,
			MaximumAttempts:    10,
		},
	})
	createOrderResponse := &models.WishCashPaymentCreateOrderResponse{}
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentCreateOrder, h, body).Get(ctx, createOrderResponse); err != nil {
		return nil, err
	}

	if createOrderResponse.Data.TransactionID == "" {
		return nil, errors.New("transaction not found")
	}

	bodyStr := fmt.Sprintf("%s&transaction_id=%s", string(body), createOrderResponse.Data.TransactionID)
	approvePaymentResponse := &models.WishCashPaymentApprovePaymentResponse{}
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentApprovePayment, h, []byte(bodyStr)).Get(ctx, approvePaymentResponse); err != nil {
		return nil, err
	}
	return approvePaymentResponse, nil
}

func GetNamespace() string {
	return path.Base(reflect.TypeOf(WishCashPaymentWorkflow{}).PkgPath())
}
