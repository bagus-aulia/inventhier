package services

import (
	"context"
	"fmt"
	"time"

	"github.com/bagus-aulia/inventhier/internal/core/domain"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

type productService struct {
	repo  ports.ProductRepository
	cache ports.ProductCache
	log   ports.ProductAuditLogger
}

// NewProductService creates a new instance of the ProductService driving port implementation.
func NewProductService(repo ports.ProductRepository, cache ports.ProductCache, log ports.ProductAuditLogger) ports.ProductService {
	return &productService{
		repo:  repo,
		cache: cache,
		log:   log,
	}
}

func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	// 1. Try to read from cache (Redis)
	if cachedProduct, err := s.cache.Get(ctx, id); err == nil {
		_ = s.log.LogAction(ctx, "GET_PRODUCT", id, "Cache Hit")
		return cachedProduct, nil
	}

	// 2. Cache miss, query SQL database
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		_ = s.log.LogAction(ctx, "GET_PRODUCT_FAILED", id, err.Error())
		return nil, err
	}

	// 3. Update cache (Redis)
	_ = s.cache.Set(ctx, product)
	_ = s.log.LogAction(ctx, "GET_PRODUCT", id, "Cache Miss, Retrieved from SQL & Cached")

	return product, nil
}

func (s *productService) CreateProduct(ctx context.Context, name string, sku string, price float64, stock int) (*domain.Product, error) {
	product := &domain.Product{
		ID:        "generate-uuid-here", // Mock ID generation
		Name:      name,
		SKU:       sku,
		Price:     price,
		Stock:     stock,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 1. Save to SQL database
	if err := s.repo.Create(ctx, product); err != nil {
		_ = s.log.LogAction(ctx, "CREATE_PRODUCT_FAILED", "", err.Error())
		return nil, err
	}

	// 2. Store in cache (Redis)
	_ = s.cache.Set(ctx, product)

	// 3. Log document action to MongoDB
	_ = s.log.LogAction(ctx, "CREATE_PRODUCT", product.ID, fmt.Sprintf("Created SKU: %s", sku))

	return product, nil
}

func (s *productService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	_ = s.log.LogAction(ctx, "LIST_PRODUCTS", "", fmt.Sprintf("Retrieved %d products", len(products)))
	return products, nil
}
