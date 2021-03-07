package wish_cash_payment

import (
	"bytes"
	"context"
	"encoding/json"
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

	bytes, err := HTTPCall(h, data, "create-order")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentCreateOrderResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (a *WishCashPaymentActivities) WishCashPaymentClearCart(ctx context.Context, h http.Header, data []byte) (*models.WishCashPaymentClearCartResponse, error) {
	a.Clients.Logger.WithFields(logrus.Fields{"headers": h, "body": string(data)}).Info("==========calling wish-frontend to clear cart==========")

	bytes, err := HTTPCall(h, data, "clear-cart")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentClearCartResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (a *WishCashPaymentActivities) WishCashPaymentApprovePayment(ctx context.Context, h http.Header, data []byte) (*models.WishCashPaymentApprovePaymentResponse, error) {
	a.Clients.Logger.WithFields(logrus.Fields{"headers": h, "body": string(data)}).Info("==========calling wish-frontend to approve payment==========")

	bytes, err := HTTPCall(h, data, "approve-payment")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentApprovePaymentResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (a *WishCashPaymentActivities) WishCashPaymentDeclinePayment(ctx context.Context, h http.Header, data []byte) (*models.WishCashPaymentDeclinePaymentResponse, error) {
	a.Clients.Logger.WithFields(logrus.Fields{"headers": h, "body": string(data)}).Info("==========calling wish-frontend to decline payment==========")

	bytes, err := HTTPCall(h, data, "decline-payment")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentDeclinePaymentResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (w *WishCashPaymentWorkflow) Register() error {
	worker := w.Clients.Temporal.DefaultClients[GetNamespace()].Worker
	worker.RegisterWorkflow(w.WishCashPaymentWorkflow)
	worker.RegisterActivity(w.Activities.WishCashPaymentCreateOrder)
	worker.RegisterActivity(w.Activities.WishCashPaymentClearCart)
	worker.RegisterActivity(w.Activities.WishCashPaymentApprovePayment)
	worker.RegisterActivity(w.Activities.WishCashPaymentDeclinePayment)

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

	if createOrderResponse.Data.FraudActionTaken != "" {
		declinePaymentResponse := &models.WishCashPaymentDeclinePaymentResponse{}
		body := fmt.Sprintf("%s&fraud_action_taken=%s", string(data), createOrderResponse.Data.FraudActionTaken)
		if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentDeclinePayment, h, []byte(body)).Get(ctx, declinePaymentResponse); err != nil {
			return nil, err
		}
		return &models.WishCashPaymentResponse{
			Data: models.WishCashPaymentResponseData{
				Msg:           declinePaymentResponse.Msg,
				TransactionID: declinePaymentResponse.Data.TransactionID,
			},
		}, nil
	}

	if createOrderResponse.Data.TransactionID == "" {
		return &models.WishCashPaymentResponse{
			Data: models.WishCashPaymentResponseData{
				Msg:           createOrderResponse.Msg,
				TransactionID: createOrderResponse.Data.TransactionID,
			},
		}, nil
	}

	clearCartResponse := &models.WishCashPaymentClearCartResponse{}
	body := fmt.Sprintf("%s&transaction_id=%s", string(data), createOrderResponse.Data.TransactionID)
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentClearCart, h, []byte(body)).Get(ctx, clearCartResponse); err != nil {
		return nil, err
	}

	approvePaymentResponse := &models.WishCashPaymentApprovePaymentResponse{}
	body = fmt.Sprintf("%s&transaction_id=%s", string(data), clearCartResponse.Data.TransactionID)
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentApprovePayment, h, []byte(body)).Get(ctx, approvePaymentResponse); err != nil {
		return nil, err
	}
	return &models.WishCashPaymentResponse{
		Data: models.WishCashPaymentResponseData{
			Msg:           approvePaymentResponse.Msg,
			TransactionID: approvePaymentResponse.Data.TransactionID,
		},
	}, nil
}

func GetNamespace() string {
	return path.Base(reflect.TypeOf(WishCashPaymentWorkflow{}).PkgPath())
}

func HTTPCall(h http.Header, data []byte, api string) ([]byte, error) {
	url := fmt.Sprintf("http://lshu.corp.contextlogic.com/api/temporal-payment/%s", api)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header = h
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
