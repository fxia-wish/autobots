package models

type (
	WishCashPaymentResponseData struct {
		Msg           string `json:"msg"`
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentCreateOrderResponseData struct {
		TransactionID    string `json:"transaction_id"`
		FraudActionTaken string `json:"fraud_action_taken"`
	}
	WishCashPaymentClearCartResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentApprovePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentDeclinePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentCreateOrderResponse struct {
		Msg         string                                 `json:"msg"`
		Code        int                                    `json:"code"`
		Data        WishCashPaymentCreateOrderResponseData `json:"data"`
		SweeperUUID string                                 `json:"sweeper_uuid"`
	}
	WishCashPaymentClearCartResponse struct {
		Msg         string                               `json:"msg"`
		Code        int                                  `json:"code"`
		Data        WishCashPaymentClearCartResponseData `json:"data"`
		SweeperUUID string                               `json:"sweeper_uuid"`
	}
	WishCashPaymentApprovePaymentResponse struct {
		Msg         string                                    `json:"msg"`
		Code        int                                       `json:"code"`
		Data        WishCashPaymentApprovePaymentResponseData `json:"data"`
		SweeperUUID string                                    `json:"sweeper_uuid"`
	}
	WishCashPaymentDeclinePaymentResponse struct {
		Msg         string                                    `json:"msg"`
		Code        int                                       `json:"code"`
		Data        WishCashPaymentDeclinePaymentResponseData `json:"data"`
		SweeperUUID string                                    `json:"sweeper_uuid"`
	}
	WishCashPaymentResponse struct {
		Data       WishCashPaymentResponseData `json:"data"`
		WorkflowID string                      `json:"workflow_id"`
		RunID      string                      `json:"run_id"`
	}
)
