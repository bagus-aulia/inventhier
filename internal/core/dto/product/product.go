package product

import "time"

// Product is struct for product
type Product struct {
	ID        int32     `json:"id,omitempty" bson:"id,omitempty"`
	SKU       string    `json:"sku" bson:"sku"`
	Name      string    `json:"name" bson:"name"`
	Price     int       `json:"price" bson:"price"`
	Supplier  string    `json:"supplier" bson:"supplier"`
	Stock     int       `json:"stock" bson:"stock"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// StockInReq is struct for restock product endpoint request
type StockInReq struct {
	ProductSKU string `json:"product_sku"`
	ProductQty int    `json:"product_qty"`
}

// StockInPayload is struct for restock product
type StockInPayload struct {
	StaffUUID string `json:"staff_uuid"`
	StockInReq
}
