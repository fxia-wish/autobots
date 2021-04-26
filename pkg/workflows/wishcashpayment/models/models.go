package models

import "net/http"

type (
	// WishCashPaymentCreateOrderResponseData data model
	WishCashPaymentCreateOrderResponseData struct {
		TransactionID    string `json:"transaction_id"`
		FraudActionTaken string `json:"fraud_action_taken"`
	}
	// WishCashPaymentCreateOrderResponse data model
	WishCashPaymentCreateOrderResponse struct {
		Context     WishCashPaymentWorkflowContext         `json:"wishCashPaymentWorkflowContext"`
		Msg         string                                 `json:"msg"`
		Code        int                                    `json:"code"`
		Data        WishCashPaymentCreateOrderResponseData `json:"data"`
		SweeperUUID string                                 `json:"sweeper_uuid"`
	}
	// WishCashPaymentClearCartResponseData data model
	WishCashPaymentClearCartResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	//WishCashPaymentClearCartResponse data model
	WishCashPaymentClearCartResponse struct {
		Context     WishCashPaymentWorkflowContext       `json:"wishCashPaymentWorkflowContext"`
		Msg         string                               `json:"msg"`
		Code        int                                  `json:"code"`
		Data        WishCashPaymentClearCartResponseData `json:"data"`
		SweeperUUID string                               `json:"sweeper_uuid"`
	}
	//WishCashPaymentApprovePaymentResponseData data model
	WishCashPaymentApprovePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	//WishCashPaymentApprovePaymentResponse data model
	WishCashPaymentApprovePaymentResponse struct {
		Context     WishCashPaymentWorkflowContext            `json:"wishCashPaymentWorkflowContext"`
		Msg         string                                    `json:"msg"`
		Code        int                                       `json:"code"`
		Data        WishCashPaymentApprovePaymentResponseData `json:"data"`
		SweeperUUID string                                    `json:"sweeper_uuid"`
	}
	//WishCashPaymentDeclinePaymentResponseData data model
	WishCashPaymentDeclinePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	//WishCashPaymentDeclinePaymentResponse data model
	WishCashPaymentDeclinePaymentResponse struct {
		Context     WishCashPaymentWorkflowContext            `json:"wishCashPaymentWorkflowContext"`
		Msg         string                                    `json:"msg"`
		Code        int                                       `json:"code"`
		Data        WishCashPaymentDeclinePaymentResponseData `json:"data"`
		SweeperUUID string                                    `json:"sweeper_uuid"`
	}
	//WishCashPaymentResponseData data model
	WishCashPaymentResponseData struct {
		Msg           string `json:"msg"`
		Code          int    `json:"code"`
		TransactionID string `json:"transaction_id"`
	}
	//WishCashPaymentResponse data model
	WishCashPaymentResponse struct {
		Context    WishCashPaymentWorkflowContext `json:"wishCashPaymentWorkflowContext"`
		Data       WishCashPaymentResponseData    `json:"data"`
		WorkflowID string                         `json:"workflow_id"`
		RunID      string                         `json:"run_id"`
	}
	//WishCashPaymentWorkflowContext data model
	WishCashPaymentWorkflowContext struct {
		Header http.Header `json:"header"`
		Body   []byte      `json:"body"`
	}
)
