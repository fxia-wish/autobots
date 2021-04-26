package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ContextLogic/autobots/pkg/clients"
	"github.com/ContextLogic/autobots/pkg/config"
	"github.com/ContextLogic/autobots/pkg/workflows"
	"github.com/ContextLogic/autobots/pkg/workflows/dummy"
	dummy_models "github.com/ContextLogic/autobots/pkg/workflows/dummy/models"
	"github.com/ContextLogic/autobots/pkg/workflows/wishcashpayment"
	"github.com/ContextLogic/autobots/pkg/workflows/wishcashpayment/models"
	wish_cash_payment_models "github.com/ContextLogic/autobots/pkg/workflows/wishcashpayment/models"
	s "github.com/ContextLogic/go-base-service/pkg/service"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	common "go.temporal.io/api/common/v1"
	enums "go.temporal.io/api/enums/v1"
	history "go.temporal.io/api/history/v1"
	ws "go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type (
	// Handlers collection
	Handlers struct {
		Config    *config.Config
		Clients   *clients.Clients
		Workflows workflows.Workflows
	}
)

// New register service APIs
func New(config *config.Config, clients *clients.Clients, service *s.Service, workflows workflows.Workflows) *Handlers {
	h := &Handlers{
		Config:    config,
		Clients:   clients,
		Workflows: workflows,
	}

	service.Mux.HandleFunc("/health", h.Health()).Methods("GET")
	service.Mux.HandleFunc("/place-order-sync", h.PlaceOrderSync()).Methods("POST")
	service.Mux.HandleFunc("/place-order-async", h.PlaceOrderAsync()).Methods("POST")
	service.Mux.HandleFunc("/reset", h.Reset()).Methods("POST")
	service.Mux.HandleFunc("/shipped", h.Shipped()).Methods("POST")
	service.Mux.HandleFunc("/start-wish-cash-payment", h.StartWishCashPayment()).Methods("POST")
	return h
}

// Unmarshal dummy request
func (h *Handlers) Unmarshal(req *http.Request) (*dummy_models.Order, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	order := dummy_models.Order{}
	err = json.Unmarshal(b, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UnmarshalResetRequest
func (h *Handlers) UnmarshalResetRequest(req *http.Request) (*dummy_models.ResetRequest, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	r := dummy_models.ResetRequest{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

//UnmarshalShippedNotificationRequest
func (h *Handlers) UnmarshalShippedNotificationRequest(req *http.Request) (*dummy_models.ShippedNotificationRequest, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	r := dummy_models.ShippedNotificationRequest{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Health check
func (h *Handlers) Health() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// PlaceOrderSync
func (h *Handlers) PlaceOrderSync() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		order, err := h.Unmarshal(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		we, err := h.Clients.Temporal.DefaultClients[dummy.GetNamespace()].Client.ExecuteWorkflow(
			context.Background(),
			client.StartWorkflowOptions{
				ID:        strings.Join([]string{dummy.GetNamespace(), strconv.Itoa(int(time.Now().Unix()))}, "_"),
				TaskQueue: fmt.Sprintf("%s_%s", h.Config.Clients.Temporal.TaskQueuePrefix, dummy.GetNamespace()),
			},
			h.Workflows[dummy.GetNamespace()].(*dummy.DummyWorkflow).DummyWorkflow,
			order,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		res := dummy_models.OrderResponse{}
		err = we.Get(context.Background(), &res)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		h.Clients.Logger.WithFields(logrus.Fields{
			"Response":   res,
			"WorkflowID": we.GetID(),
			"RunID":      we.GetRunID(),
		}).Info("workflow execution completed")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&dummy_models.PlaceOrderSyncResponse{
			OrderResponse: &res,
			WorkflowID:    we.GetID(),
			RunID:         we.GetRunID(),
		})
	}
}

// PlaceOrderAsync
func (h *Handlers) PlaceOrderAsync() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		order, err := h.Unmarshal(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		we, err := h.Clients.Temporal.DefaultClients[dummy.GetNamespace()].Client.ExecuteWorkflow(
			context.Background(),
			client.StartWorkflowOptions{
				ID:        strings.Join([]string{dummy.GetNamespace(), strconv.Itoa(int(time.Now().Unix()))}, "_"),
				TaskQueue: fmt.Sprintf("%s_%s", h.Config.Clients.Temporal.TaskQueuePrefix, dummy.GetNamespace()),
			},
			h.Workflows[dummy.GetNamespace()].(*dummy.DummyWorkflow).DummyWorkflow,
			order,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		h.Clients.Logger.WithFields(logrus.Fields{
			"WorkflowID": we.GetID(),
			"RunID":      we.GetRunID(),
		}).Info("workflow execution started")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&dummy_models.PlaceOrderAsyncResponse{
			WorkflowID: we.GetID(),
			RunID:      we.GetRunID(),
		})
	}
}

// Shipped return function that parse shipped notification
func (h *Handlers) Shipped() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rr, err := h.UnmarshalShippedNotificationRequest(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		histories, err := h.Clients.Temporal.WorkflowServiceClient.GetWorkflowExecutionHistory(
			context.Background(),
			&ws.GetWorkflowExecutionHistoryRequest{
				Namespace: dummy.GetNamespace(),
				Execution: &common.WorkflowExecution{
					WorkflowId: rr.WorkflowID,
					RunId:      rr.RunID,
				},
			},
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, event := range histories.History.Events {
			if event.GetEventType() == enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED {
				if event.GetActivityTaskScheduledEventAttributes().GetActivityType().GetName() == "DummyShipping" {
					rr.ActivityID = event.GetActivityTaskScheduledEventAttributes().GetActivityId()
					break
				}
			}
		}

		err = h.Clients.Temporal.DefaultClients[dummy.GetNamespace()].Client.CompleteActivityByID(
			context.Background(),
			dummy.GetNamespace(),
			rr.WorkflowID,
			rr.RunID,
			rr.ActivityID,
			rr.Result,
			nil,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return

		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode("ok")
	}
}

// Reset return function that reset workflow execution
func (h *Handlers) Reset() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rr, err := h.UnmarshalResetRequest(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		histories, err := h.Clients.Temporal.WorkflowServiceClient.GetWorkflowExecutionHistory(
			context.Background(),
			&ws.GetWorkflowExecutionHistoryRequest{
				Namespace: dummy.GetNamespace(),
				Execution: &common.WorkflowExecution{
					WorkflowId: rr.WorkflowID,
					RunId:      rr.RunID,
				},
			},
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var lastWorkflowTask *history.HistoryEvent
		for _, event := range histories.History.Events {
			if event.GetEventType() == enums.EVENT_TYPE_ACTIVITY_TASK_FAILED {
				break
			}
			if event.GetEventType() == enums.EVENT_TYPE_WORKFLOW_TASK_COMPLETED {
				lastWorkflowTask = event
			}
		}

		_, err = h.Clients.Temporal.WorkflowServiceClient.ResetWorkflowExecution(
			context.Background(),
			&ws.ResetWorkflowExecutionRequest{
				Namespace: dummy.GetNamespace(),
				WorkflowExecution: &common.WorkflowExecution{
					WorkflowId: rr.WorkflowID,
					RunId:      rr.RunID,
				},
				Reason:                    "reset execution for failure",
				WorkflowTaskFinishEventId: lastWorkflowTask.GetEventId(),
				RequestId:                 uuid.New(),
			},
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode("ok")
	}
}

//StartWishCashPayment return function that triggers wish cash payment workflow
func (h *Handlers) StartWishCashPayment() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()

		data := &models.WishCashPaymentWorkflowContext{
			Header: req.Header,
			Body:   []byte(body),
		}

		h.Clients.Logger.Info("workflow execution triggered")
		we, err := h.Clients.Temporal.DefaultClients[wishcashpayment.GetNamespace()].Client.ExecuteWorkflow(
			context.Background(),
			client.StartWorkflowOptions{
				ID:        strings.Join([]string{wishcashpayment.GetNamespace(), strconv.Itoa(int(time.Now().Unix()))}, "_"),
				TaskQueue: fmt.Sprintf("%s_%s", h.Config.Clients.Temporal.TaskQueuePrefix, wishcashpayment.GetNamespace()),
			},
			h.Workflows[wishcashpayment.GetNamespace()].(*wishcashpayment.WishCashPaymentWorkflow).WishCashPaymentWorkflow,
			data,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		response := &wish_cash_payment_models.WishCashPaymentResponse{}
		err = we.Get(context.Background(), &response)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		response.WorkflowID = we.GetID()
		response.RunID = we.GetRunID()

		h.Clients.Logger.WithFields(logrus.Fields{
			"api response": response,
		}).Info("workflow execution completed")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}
}
