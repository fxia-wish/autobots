package models

type (
	Order struct {
		ProductID       string `json:"product_id"`
		CustomerID      string `json:"customer_id"`
		ShippingAddress string `json:"shipping_address"`
	}
	OrderResponse struct {
		Order  *Order `json:"order"`
		Status string `json:"status"`
	}
	PlaceOrderSyncResponse struct {
		OrderResponse *OrderResponse `json:"order_response"`
		WorkflowID    string         `json:"workflow_id"`
		RunID         string         `json:"run_id"`
	}
	PlaceOrderAsyncResponse struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
	}
	ResetRequest struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
	}
	ShippedNotificationRequest struct {
		WorkflowID string `json:"workflow_id"`
		RunID      string `json:"run_id"`
		ActivityID string
		Result     *OrderResponse `json:"result"`
	}
)