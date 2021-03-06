package models

type (
	WishCashPaymentCreateOrderResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentApprovePaymentResponseData struct {
		TransactionID string `json:"transaction_id"`
	}
	WishCashPaymentCreateOrderResponse struct {
		Msg         string                                  `json:"msg"`
		Code        int                                     `json:"code"`
		Data        *WishCashPaymentCreateOrderResponseData `json:"data"`
		SweeperUUID string                                  `json:"sweeper_uuid"`
	}
	WishCashPaymentApprovePaymentResponse struct {
		Msg         string                                     `json:"msg"`
		Code        int                                        `json:"code"`
		Data        *WishCashPaymentApprovePaymentResponseData `json:"data"`
		SweeperUUID string                                     `json:"sweeper_uuid"`
	}
	WishCashPaymentResponse struct {
		TransactionID string `json:"transaction_id"`
		WorkflowID    string `json:"workflow_id"`
		RunID         string `json:"run_id"`
	}
)
