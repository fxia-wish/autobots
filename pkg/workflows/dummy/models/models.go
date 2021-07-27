package models

type (
	// SimpleSignal detail
	SimpleSignal struct {
		Namespace  string `json:"namespace"`
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
		Name       string `json:"signal_name"`
		Value      string `json:"signal_val"`
	}

	// Order detail
	Order struct {
		ProductID       string `json:"product_id"`
		CustomerID      string `json:"customer_id"`
		ShippingAddress string `json:"shipping_address"`
	}
	// OrderResponse detail
	OrderResponse struct {
		Order  *Order `json:"order"`
		Status string `json:"status"`
	}
	// PlaceOrderSyncResponse detail
	PlaceOrderSyncResponse struct {
		OrderResponse *OrderResponse `json:"order_response"`
		WorkflowID    string         `json:"workflow_id"`
		RunID         string         `json:"run_id"`
	}
	// PlaceOrderAsyncResponse detail
	PlaceOrderAsyncResponse struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
	}
	//ResetRequest detail
	ResetRequest struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
	}
	//ShippedNotificationRequest detail
	ShippedNotificationRequest struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
		ActivityID string
		Result     *OrderResponse `json:"result"`
	}
)
