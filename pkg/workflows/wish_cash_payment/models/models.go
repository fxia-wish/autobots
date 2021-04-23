package models

import "net/http"

type (
	WishCashPaymentCreateOrderResponseData struct {
		TransactionID    string `json:"transaction_id"`
		FraudActionTaken string `json:"fraud_action_taken"`
	}
	WishCashPaymentCreateOrderResponse struct {
		Header      http.Header                            `json:"header"`
		Body        []byte                                 `json:"body"`
		Msg         string                                 `json:"msg"`
		Code        int                                    `json:"code"`
		Data        WishCashPaymentCreateOrderResponseData `json:"data"`
		SweeperUUID string                                 `json:"sweeper_uuid"`
	}

	WishCashPaymentClearCartResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentClearCartResponse struct {
		Header      http.Header                          `json:"header"`
		Body        []byte                               `json:"body"`
		Msg         string                               `json:"msg"`
		Code        int                                  `json:"code"`
		Data        WishCashPaymentClearCartResponseData `json:"data"`
		SweeperUUID string                               `json:"sweeper_uuid"`
	}

	WishCashPaymentApprovePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentApprovePaymentResponse struct {
		Header      http.Header                               `json:"header"`
		Body        []byte                                    `json:"body"`
		Msg         string                                    `json:"msg"`
		Code        int                                       `json:"code"`
		Data        WishCashPaymentApprovePaymentResponseData `json:"data"`
		SweeperUUID string                                    `json:"sweeper_uuid"`
	}

	WishCashPaymentDeclinePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentDeclinePaymentResponse struct {
		Header      http.Header                               `json:"header"`
		Body        []byte                                    `json:"body"`
		Msg         string                                    `json:"msg"`
		Code        int                                       `json:"code"`
		Data        WishCashPaymentDeclinePaymentResponseData `json:"data"`
		SweeperUUID string                                    `json:"sweeper_uuid"`
	}

	WishCashPaymentResponseData struct {
		Msg           string `json:"msg"`
		Code          int    `json:"code"`
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentResponse struct {
		Header     http.Header                 `json:"header"`
		Body       []byte                      `json:"body"`
		Data       WishCashPaymentResponseData `json:"data"`
		WorkflowID string                      `json:"workflow_id"`
		RunID      string                      `json:"run_id"`
	}

	WishCashPaymentWorkflowInput struct {
		Header http.Header `json:"header"`
		Body   []byte      `json:"body"`
	}
)
