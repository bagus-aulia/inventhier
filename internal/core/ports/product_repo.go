package ports

import (
	"context"

	dto "github.com/bagus-aulia/inventhier/internal/core/dto/product"
)

// ProductLogger is a driven port defining product logging (e.g. MongoDB).
type ProductLogger interface {
	StoreProductLog(ctx context.Context, log dto.ProductLog) error
}

// Product is a driven port defining product data
type Product interface {
	GetProductBySKU(ctx context.Context, sku string) (*dto.Product, error)
	UpdateProductStock(ctx context.Context, sku string, stockIn int, staffUUID string) error
	CreateProduct(ctx context.Context, data dto.Product) error
}

// ProductCache is a driven port defining caching operations (e.g. Redis).
type ProductCache interface {
	GetProductBySKU(ctx context.Context, sku string) (*dto.Product, error)
	DelProductCache(ctx context.Context, sku string) error
}
