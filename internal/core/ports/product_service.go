package ports

import (
	"context"

	productDTO "github.com/bagus-aulia/inventhier/internal/core/dto/product"
)

// ProductService is a driving port defining business use cases.
type ProductService interface {
	StockIn(c context.Context, payload productDTO.StockInPayload) error
}
