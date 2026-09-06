package ports

import (
	"context"

	"github.com/bagus-aulia/inventhier/internal/core/domain"
)

// ============================================================================
// PRODUCT DRIVEN PORTS (Adapters that Product service depends on)
// ============================================================================

// ProductRepository is a driven port defining relational database operations (e.g. MySQL, PostgreSQL).
type ProductRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, product *domain.Product) error
	List(ctx context.Context) ([]domain.Product, error)
}

// ProductCache is a driven port defining caching operations (e.g. Redis).
type ProductCache interface {
	Get(ctx context.Context, id string) (*domain.Product, error)
	Set(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
}

// ProductAuditLogger is a driven port defining audit logging or document logging (e.g. MongoDB).
type ProductAuditLogger interface {
	LogAction(ctx context.Context, action string, productID string, details string) error
}

// ============================================================================
// PRODUCT DRIVING PORT (Service interface exposed to outside world)
// ============================================================================

// ProductService is a driving port defining business use cases.
type ProductService interface {
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	CreateProduct(ctx context.Context, name string, sku string, price float64, stock int) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, error)
}
