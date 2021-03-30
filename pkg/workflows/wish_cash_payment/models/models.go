package models

type (
	WishCashPaymentCreateOrderResponseData struct {
		TransactionID    string `json:"transaction_id"`
		FraudActionTaken string `json:"fraud_action_taken"`
	}
	WishCashPaymentCreateOrderResponse struct {
		Msg         string                                 `json:"msg"`
		Code        int                                    `json:"code"`
		Data        WishCashPaymentCreateOrderResponseData `json:"data"`
		SweeperUUID string                                 `json:"sweeper_uuid"`
	}

	WishCashPaymentClearCartResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentClearCartResponse struct {
		Msg         string                               `json:"msg"`
		Code        int                                  `json:"code"`
		Data        WishCashPaymentClearCartResponseData `json:"data"`
		SweeperUUID string                               `json:"sweeper_uuid"`
	}

	WishCashPaymentApprovePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentApprovePaymentResponse struct {
		Msg         string                                    `json:"msg"`
		Code        int                                       `json:"code"`
		Data        WishCashPaymentApprovePaymentResponseData `json:"data"`
		SweeperUUID string                                    `json:"sweeper_uuid"`
	}

	WishCashPaymentDeclinePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentDeclinePaymentResponse struct {
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
		Data       WishCashPaymentResponseData `json:"data"`
		WorkflowID string                      `json:"workflow_id"`
		RunID      string                      `json:"run_id"`
	}
)
