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
	"time"

	"github.com/ContextLogic/autobots/clients"
	"github.com/ContextLogic/autobots/workflows/wish_cash_payment/models"
	"github.com/sirupsen/logrus"
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
	}
)

func NewWishCashPaymentWorkflow(clients *clients.Clients) *WishCashPaymentWorkflow {
	return &WishCashPaymentWorkflow{
		Clients: clients,
		Activities: &WishCashPaymentActivities{
			Clients: clients,
		},
	}
}

func (a *WishCashPaymentActivities) WishCashPaymentCreateOrder(ctx context.Context, h http.Header, data []byte) (*models.WishCashPaymentCreateOrderResponse, error) {
	a.Clients.Logger.WithFields(logrus.Fields{"headers": h, "body": string(data)}).Info("==========calling wish-frontend to create order==========")
	req, err := http.NewRequest(http.MethodPost, "http://lshu.corp.contextlogic.com/api/temporal-payment/create-order", bytes.NewBuffer(data))
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

func (a *WishCashPaymentActivities) WishCashPaymentApprovePayment(ctx context.Context, h http.Header, data []byte) (*models.WishCashPaymentApprovePaymentResponse, error) {
	a.Clients.Logger.WithFields(logrus.Fields{"headers": h, "body": string(data)}).Info("==========calling wish-frontend to approve payment==========")
	req, err := http.NewRequest(http.MethodPost, "http://lshu.corp.contextlogic.com/api/temporal-payment/approve-payment", bytes.NewBuffer(data))
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
	worker := w.Clients.Temporal.DefaultClients[GetNamespace()].Worker
	worker.RegisterWorkflow(w.WishCashPaymentWorkflow)
	worker.RegisterActivity(w.Activities.WishCashPaymentCreateOrder)
	worker.RegisterActivity(w.Activities.WishCashPaymentApprovePayment)

	return nil
}

func (w *WishCashPaymentWorkflow) WishCashPaymentWorkflow(ctx workflow.Context, h http.Header, data []byte) (interface{}, error) {
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
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentCreateOrder, h, data).Get(ctx, createOrderResponse); err != nil {
		return nil, err
	}

	if createOrderResponse.Data.TransactionID == "" {
		return nil, errors.New("transaction not found")
	}

	body := fmt.Sprintf("%s&transaction_id=%s", string(data), createOrderResponse.Data.TransactionID)
	approvePaymentResponse := &models.WishCashPaymentApprovePaymentResponse{}
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentApprovePayment, h, []byte(body)).Get(ctx, approvePaymentResponse); err != nil {
		return nil, err
	}
	return approvePaymentResponse, nil
}

func GetNamespace() string {
	return path.Base(reflect.TypeOf(WishCashPaymentWorkflow{}).PkgPath())
}
