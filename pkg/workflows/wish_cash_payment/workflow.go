package wish_cash_payment

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"time"

	"github.com/ContextLogic/autobots/pkg/clients"
	"github.com/ContextLogic/autobots/pkg/config"
	"github.com/ContextLogic/autobots/pkg/workflows/wish_cash_payment/models"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	// WishCashPaymentWorkflow config and clients
	WishCashPaymentWorkflow struct {
		Config     *config.TemporalClientConfig
		Clients    *clients.Clients
		Activities *WishCashPaymentActivities
	}
	// WishCashPaymentActivities
	WishCashPaymentActivities struct {
		Clients *clients.Clients
	}
)

// create wish cash payment workflow
func NewWishCashPaymentWorkflow(config *config.TemporalClientConfig, clients *clients.Clients) *WishCashPaymentWorkflow {
	return &WishCashPaymentWorkflow{
		Config:  config,
		Clients: clients,
		Activities: &WishCashPaymentActivities{
			Clients: clients,
		},
	}
}

// create wish cash order request
func (a *WishCashPaymentActivities) WishCashPaymentCreateOrder(ctx context.Context, input interface{}) (interface{}, error) {
	objInput := &models.WishCashPaymentWorkflowContext{}
	err := GetInputObject(input, objInput)
	if err != nil {
		return nil, err
	}
	a.Clients.Logger.Info("==========calling wish-fe to create order: started==========")
	a.Clients.Logger.WithFields(logrus.Fields{"headers": objInput.Header, "body": string(objInput.Body)}).Info("create order request info")

	bytes, err := a.Clients.WishFrontend.Post(objInput.Header, objInput.Body, "api/temporal-payment/create-order")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentCreateOrderResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}
	response.Context.Header = objInput.Header
	response.Context.Body = []byte(fmt.Sprintf("%s&transaction_id=%s", string(objInput.Body), response.Data.TransactionID))

	a.Clients.Logger.Info("==========calling wish-fe to create order: finished==========")
	a.Clients.Logger.WithFields(logrus.Fields{"response": response}).Info("create order response info")

	return response, nil
}

// create request to clear cart
func (a *WishCashPaymentActivities) WishCashPaymentClearCart(ctx context.Context, input interface{}) (interface{}, error) {
	objInput := &models.WishCashPaymentCreateOrderResponse{}
	err := GetInputObject(input, objInput)
	if err != nil {
		return nil, err
	}

	a.Clients.Logger.Info("==========calling wish-fe to clear cart: started==========")
	a.Clients.Logger.WithFields(logrus.Fields{"headers": objInput.Context.Header, "body": string(objInput.Context.Body)}).Info("clear cart request info")

	bytes, err := a.Clients.WishFrontend.Post(objInput.Context.Header, objInput.Context.Body, "api/temporal-payment/clear-cart")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentClearCartResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	response.Context.Header = objInput.Context.Header
	response.Context.Body = objInput.Context.Body
	a.Clients.Logger.Info("==========calling wish-fe to clear cart: finished==========")
	a.Clients.Logger.WithFields(logrus.Fields{"response": response}).Info("clear cart response info")

	return response, nil
}

// create approve payment request
func (a *WishCashPaymentActivities) WishCashPaymentApprovePayment(ctx context.Context, input interface{}) (interface{}, error) {
	objInput := &models.WishCashPaymentClearCartResponse{}
	err := GetInputObject(input, objInput)
	if err != nil {
		return nil, err
	}

	a.Clients.Logger.Info("==========calling wish-fe to approve payment: started==========")
	a.Clients.Logger.WithFields(logrus.Fields{"headers": objInput.Context.Header, "body": string(objInput.Context.Body)}).Info("approve payment request info")

	bytes, err := a.Clients.WishFrontend.Post(objInput.Context.Header, objInput.Context.Body, "api/temporal-payment/approve-payment")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentApprovePaymentResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	response.Context.Header = objInput.Context.Header
	response.Context.Body = objInput.Context.Body
	a.Clients.Logger.Info("==========calling wish-fe to approve payment: finished==========")
	a.Clients.Logger.WithFields(logrus.Fields{"response": response}).Info("approve payment response info")

	return response, nil
}

// create decline payent request
func (a *WishCashPaymentActivities) WishCashPaymentDeclinePayment(ctx context.Context, input interface{}) (interface{}, error) {
	objInput := &models.WishCashPaymentCreateOrderResponse{}
	err := GetInputObject(input, objInput)
	if err != nil {
		return nil, err
	}

	a.Clients.Logger.Info("==========calling wish-fe to decline payment: started==========")
	a.Clients.Logger.WithFields(logrus.Fields{"headers": objInput.Context.Header, "body": string(objInput.Context.Body)}).Info("decline payment request info")

	bytes, err := a.Clients.WishFrontend.Post(objInput.Context.Header, objInput.Context.Body, "api/temporal-payment/decline-payment")
	if err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentDeclinePaymentResponse{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}
	response.Context.Header = objInput.Context.Header
	response.Context.Body = objInput.Context.Body
	a.Clients.Logger.Info("==========calling wish-fe to decline payment: finished==========")
	a.Clients.Logger.WithFields(logrus.Fields{"response": response}).Info("decline payment response info")

	return response, nil
}

// register workflow namespace
func (w *WishCashPaymentWorkflow) Register() error {
	if err := w.Clients.Temporal.RegisterNamespace(GetNamespace(), w.Config.Retention); err != nil {
		return err
	}

	worker := w.Clients.Temporal.DefaultClients[GetNamespace()].Worker
	worker.RegisterWorkflow(w.WishCashPaymentWorkflow)
	worker.RegisterActivity(w.Activities.WishCashPaymentCreateOrder)
	worker.RegisterActivity(w.Activities.WishCashPaymentClearCart)
	worker.RegisterActivity(w.Activities.WishCashPaymentApprovePayment)
	worker.RegisterActivity(w.Activities.WishCashPaymentDeclinePayment)

	return nil
}

// execute wish cash payment workflow in code
func (w *WishCashPaymentWorkflow) WishCashPaymentWorkflow(ctx workflow.Context, input interface{}) (interface{}, error) {
	c := w.Config.Activities
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * time.Duration(c.StartToCloseTimeout),
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * time.Duration(c.RetryPolicy.InitialInterval),
			BackoffCoefficient: c.RetryPolicy.BackoffCoefficient,
			MaximumInterval:    time.Second * time.Duration(c.RetryPolicy.MaximumInterval),
			MaximumAttempts:    c.RetryPolicy.MaximumAttempts,
		},
	})

	createOrderResponse := &models.WishCashPaymentCreateOrderResponse{}
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentCreateOrder, input).Get(ctx, createOrderResponse); err != nil {
		return nil, err
	}
	if createOrderResponse.Data.FraudActionTaken != "" {
		declinePaymentResponse := &models.WishCashPaymentDeclinePaymentResponse{}
		createOrderResponse.Context.Body = []byte(fmt.Sprintf("%s&fraud_action_taken=%s&transaction_id=%s", string(createOrderResponse.Context.Body), createOrderResponse.Data.FraudActionTaken, createOrderResponse.Data.TransactionID))
		if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentDeclinePayment, createOrderResponse).Get(ctx, declinePaymentResponse); err != nil {
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
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentClearCart, createOrderResponse).Get(ctx, clearCartResponse); err != nil {
		return nil, err
	}

	approvePaymentResponse := &models.WishCashPaymentApprovePaymentResponse{}
	if err := workflow.ExecuteActivity(ctx, w.Activities.WishCashPaymentApprovePayment, createOrderResponse).Get(ctx, approvePaymentResponse); err != nil {
		return nil, err
	}

	response := &models.WishCashPaymentResponse{
		Data: models.WishCashPaymentResponseData{
			Msg:           approvePaymentResponse.Msg,
			Code:          approvePaymentResponse.Code,
			TransactionID: approvePaymentResponse.Data.TransactionID,
		},
	}
	w.Clients.Logger.WithField("response", response).Info("workflow response info")
	return response, nil
}

// get workflow namespace
func GetNamespace() string {
	return path.Base(reflect.TypeOf(WishCashPaymentWorkflow{}).PkgPath())
}

// get activity map for workflow
func GetActivityMap(w *WishCashPaymentWorkflow) map[string]func(ctx context.Context, input interface{}) (interface{}, error) {
	activityMap := map[string]func(ctx context.Context, input interface{}) (interface{}, error){
		"WishCashPaymentCreateOrder":    w.Activities.WishCashPaymentCreateOrder,
		"WishCashPaymentClearCart":      w.Activities.WishCashPaymentClearCart,
		"WishCashPaymentApprovePayment": w.Activities.WishCashPaymentApprovePayment,
		"WishCashPaymentDeclinePayment": w.Activities.WishCashPaymentDeclinePayment,
	}
	return activityMap
}

// unmarshal workflow input from interface to objects
func GetInputObject(input interface{}, objInput interface{}) error {
	data, _ := json.Marshal(input)
	err := json.Unmarshal(data, &objInput)
	if err != nil {
		return err
	}
	return nil
}
