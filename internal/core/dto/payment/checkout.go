package payment

// CheckoutPayload is payload request for POST Checkout client endpoint
type CheckoutPayload struct {
	TransactionID string            `json:"transaction_id"`
	StaffID       string            `json:"staff_id"`
	Products      []CheckoutProduct `json:"products"`
}

// CheckoutProduct is product data for checkout payload
type CheckoutProduct struct {
	ProductID string `json:"product_id"`
	Quantity  string `json:"qty"`
}

// CheckoutResp is response for POST Checkout client endpoint
type CheckoutResp struct {
	TransactionStatus string `json:"transaction_status"`
}
