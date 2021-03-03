package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ContextLogic/autobots/clients"
	"github.com/ContextLogic/autobots/config"
	"github.com/ContextLogic/autobots/models"
	"github.com/ContextLogic/autobots/workflows"
	ct "github.com/ContextLogic/autobots/workflows/dummy"
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
	Handlers struct {
		Config    *config.Config
		Clients   *clients.Clients
		Workflows workflows.Workflows
	}
)

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
	return h
}

func (h *Handlers) Unmarshal(req *http.Request) (*models.Order, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	order := models.Order{}
	err = json.Unmarshal(b, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (h *Handlers) UnmarshalResetRequest(req *http.Request) (*models.ResetRequest, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	r := models.ResetRequest{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (h *Handlers) UnmarshalShippedNotificationRequest(req *http.Request) (*models.ShippedNotificationRequest, error) {
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	r := models.ShippedNotificationRequest{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (h *Handlers) Health() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handlers) PlaceOrderSync() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		order, err := h.Unmarshal(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		ts := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
		we, err := h.Clients.Temporal.Client.ExecuteWorkflow(
			context.Background(),
			client.StartWorkflowOptions{
				ID:        strings.Join([]string{ct.GetNamespace(), ts}, "_"),
				TaskQueue: h.Config.Clients.Temporal.TaskQueue,
			},
			h.Workflows[ct.GetNamespace()].Entry,
			order,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		res := models.OrderResponse{}
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
		json.NewEncoder(w).Encode(&models.PlaceOrderSyncResponse{
			OrderResponse: &res,
			WorkflowID:    we.GetID(),
			RunID:         we.GetRunID(),
		})
	}
}

func (h *Handlers) PlaceOrderAsync() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		order, err := h.Unmarshal(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		ts := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-")
		we, err := h.Clients.Temporal.Client.ExecuteWorkflow(
			context.Background(),
			client.StartWorkflowOptions{
				ID:        strings.Join([]string{"dummy", ts}, "_"),
				TaskQueue: h.Config.Clients.Temporal.TaskQueue,
			},
			h.Workflows["dummy"].Entry,
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
		json.NewEncoder(w).Encode(&models.PlaceOrderAsyncResponse{
			WorkflowID: we.GetID(),
			RunID:      we.GetRunID(),
		})
	}
}

func (h *Handlers) Shipped() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rr, err := h.UnmarshalShippedNotificationRequest(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		histories, err := h.Clients.Temporal.Frontend.GetWorkflowExecutionHistory(
			context.Background(),
			&ws.GetWorkflowExecutionHistoryRequest{
				Namespace: ct.GetNamespace(),
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
				if event.GetActivityTaskScheduledEventAttributes().GetActivityType().GetName() == "Shipping" {
					rr.ActivityID = event.GetActivityTaskScheduledEventAttributes().GetActivityId()
					break
				}
			}
		}

		err = h.Clients.Temporal.Client.CompleteActivityByID(
			context.Background(),
			ct.GetNamespace(),
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

func (h *Handlers) Reset() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rr, err := h.UnmarshalResetRequest(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		histories, err := h.Clients.Temporal.Frontend.GetWorkflowExecutionHistory(
			context.Background(),
			&ws.GetWorkflowExecutionHistoryRequest{
				Namespace: ct.GetNamespace(),
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

		_, err = h.Clients.Temporal.Frontend.ResetWorkflowExecution(
			context.Background(),
			&ws.ResetWorkflowExecutionRequest{
				Namespace: ct.GetNamespace(),
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
