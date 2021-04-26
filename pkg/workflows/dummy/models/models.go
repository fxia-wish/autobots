package models

type (
	// order detail
	Order struct {
		ProductID       string `json:"product_id"`
		CustomerID      string `json:"customer_id"`
		ShippingAddress string `json:"shipping_address"`
	}
	// order response detail
	OrderResponse struct {
		Order  *Order `json:"order"`
		Status string `json:"status"`
	}
	// place order sync response detail
	PlaceOrderSyncResponse struct {
		OrderResponse *OrderResponse `json:"order_response"`
		WorkflowID    string         `json:"workflow_id"`
		RunID         string         `json:"run_id"`
	}
	// place order async response detail
	PlaceOrderAsyncResponse struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
	}
	//reset reqeust detail
	ResetRequest struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
	}
	//ship notification reqeust detail
	ShippedNotificationRequest struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
		ActivityID string
		Result     *OrderResponse `json:"result"`
	}
)
