package product_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	"github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/v1/user/mocks"
	productLogMocks "github.com/bagus-aulia/inventhier/internal/adapters/repository/mongodb/product_log/mocks"
	productCacheMocks "github.com/bagus-aulia/inventhier/internal/adapters/repository/redis/product/mocks"
	productRepoMocks "github.com/bagus-aulia/inventhier/internal/adapters/repository/sql/product/mocks"
	"github.com/bagus-aulia/inventhier/internal/core/constants"
	"github.com/bagus-aulia/inventhier/internal/core/dto/product"
	userDTO "github.com/bagus-aulia/inventhier/internal/core/dto/user"
	productService "github.com/bagus-aulia/inventhier/internal/core/services/v1/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	cfg = &config.Config{
		ContextTimeout: 5,
	}

	mockProduct = &product.Product{
		ID:        1,
		SKU:       "SKU-001",
		Name:      "Test Product",
		Price:     100000,
		Supplier:  "Supplier A",
		Stock:     50,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockStockInPayload = product.StockInPayload{
		StaffUUID: "staff-uuid-123",
		StockInReq: product.StockInReq{
			ProductSKU: "SKU-001",
			ProductQty: 10,
		},
	}
)

// TestNewProductService tests the constructor
func TestNewProductService(t *testing.T) {
	t.Run("creates new product service instance", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		assert.NotNil(t, svc)
	})
}

// TestStockIn tests the StockIn method
func TestStockIn(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully stocks in product", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Mock user client to return user
		userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
			Return(&userDTO.User{UUID: mockStockInPayload.StaffUUID}, nil)

		// Mock product cache to return product
		productCache.On("GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU).
			Return(mockProduct, nil)

		// Mock product repo to update stock
		productRepo.On("UpdateProductStock", mock.Anything, mockStockInPayload.ProductSKU, mockStockInPayload.ProductQty, mockStockInPayload.StaffUUID).
			Return(nil)

		// Mock product cache to delete cache
		productCache.On("DelProductCache", mock.Anything, mockStockInPayload.ProductSKU).
			Return(nil)

		// Mock product log to store log
		productLog.On("StoreProductLog", mock.Anything, mock.MatchedBy(func(log product.ProductLog) bool {
			return log.UUID == mockStockInPayload.StaffUUID &&
				log.Status == constants.ProductStatusIn &&
				log.Quantity == mockStockInPayload.ProductQty
		})).Return(nil)

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		err := svc.StockIn(ctx, mockStockInPayload)

		assert.NoError(t, err)
		userClient.AssertCalled(t, "GetUser", mock.Anything, mockStockInPayload.StaffUUID)
		productCache.AssertCalled(t, "GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU)
		productRepo.AssertCalled(t, "UpdateProductStock", mock.Anything, mockStockInPayload.ProductSKU, mockStockInPayload.ProductQty, mockStockInPayload.StaffUUID)
		productCache.AssertCalled(t, "DelProductCache", mock.Anything, mockStockInPayload.ProductSKU)
		productLog.AssertCalled(t, "StoreProductLog", mock.Anything, mock.Anything)
	})

	t.Run("returns error when user validation fails", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Mock user client to return error
		userErr := errors.New("user not found")
		userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
			Return(nil, userErr)

		// Mock product log to store log (even on error)
		productLog.On("StoreProductLog", mock.Anything, mock.Anything).Return(nil)

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		err := svc.StockIn(ctx, mockStockInPayload)

		assert.Error(t, err)
		assert.Equal(t, userErr, err)
		userClient.AssertCalled(t, "GetUser", mock.Anything, mockStockInPayload.StaffUUID)
		productCache.AssertNotCalled(t, "GetProductBySKU")
		productRepo.AssertNotCalled(t, "UpdateProductStock")
	})

	t.Run("returns error when product retrieval fails", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Mock user client to return user
		userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
			Return(&userDTO.User{UUID: mockStockInPayload.StaffUUID}, nil)

		// Mock product cache to return error
		cacheErr := errors.New("product not found in cache")
		productCache.On("GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU).
			Return(nil, cacheErr)

		// Mock product log to store log
		productLog.On("StoreProductLog", mock.Anything, mock.Anything).Return(nil)

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		err := svc.StockIn(ctx, mockStockInPayload)

		assert.Error(t, err)
		assert.Equal(t, cacheErr, err)
		productCache.AssertCalled(t, "GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU)
		productRepo.AssertNotCalled(t, "UpdateProductStock")
	})

	t.Run("returns error when stock update fails", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Mock user client to return user
		userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
			Return(&userDTO.User{UUID: mockStockInPayload.StaffUUID}, nil)

		// Mock product cache to return product
		productCache.On("GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU).
			Return(mockProduct, nil)

		// Mock product repo to return error
		repoErr := errors.New("database error")
		productRepo.On("UpdateProductStock", mock.Anything, mockStockInPayload.ProductSKU, mockStockInPayload.ProductQty, mockStockInPayload.StaffUUID).
			Return(repoErr)

		// Mock product log to store log
		productLog.On("StoreProductLog", mock.Anything, mock.Anything).Return(nil)

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		err := svc.StockIn(ctx, mockStockInPayload)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		productRepo.AssertCalled(t, "UpdateProductStock", mock.Anything, mockStockInPayload.ProductSKU, mockStockInPayload.ProductQty, mockStockInPayload.StaffUUID)
		productCache.AssertNotCalled(t, "DelProductCache")
	})

	t.Run("returns error when cache deletion fails", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Mock user client to return user
		userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
			Return(&userDTO.User{UUID: mockStockInPayload.StaffUUID}, nil)

		// Mock product cache to return product
		productCache.On("GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU).
			Return(mockProduct, nil)

		// Mock product repo to update stock
		productRepo.On("UpdateProductStock", mock.Anything, mockStockInPayload.ProductSKU, mockStockInPayload.ProductQty, mockStockInPayload.StaffUUID).
			Return(nil)

		// Mock product cache to return error on delete
		delErr := errors.New("cache deletion failed")
		productCache.On("DelProductCache", mock.Anything, mockStockInPayload.ProductSKU).
			Return(delErr)

		// Mock product log to store log
		productLog.On("StoreProductLog", mock.Anything, mock.Anything).Return(nil)

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		err := svc.StockIn(ctx, mockStockInPayload)

		assert.Error(t, err)
		assert.Equal(t, delErr, err)
		productCache.AssertCalled(t, "DelProductCache", mock.Anything, mockStockInPayload.ProductSKU)
	})

	t.Run("stores product log even if log storage fails", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Mock user client to return user
		userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
			Return(&userDTO.User{UUID: mockStockInPayload.StaffUUID}, nil)

		// Mock product cache to return product
		productCache.On("GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU).
			Return(mockProduct, nil)

		// Mock product repo to update stock
		productRepo.On("UpdateProductStock", mock.Anything, mockStockInPayload.ProductSKU, mockStockInPayload.ProductQty, mockStockInPayload.StaffUUID).
			Return(nil)

		// Mock product cache to delete cache
		productCache.On("DelProductCache", mock.Anything, mockStockInPayload.ProductSKU).
			Return(nil)

		// Mock product log to return error
		logErr := errors.New("log storage failed")
		productLog.On("StoreProductLog", mock.Anything, mock.Anything).Return(logErr)

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		err := svc.StockIn(ctx, mockStockInPayload)

		// Should still succeed even if log storage fails (defer doesn't propagate error)
		assert.NoError(t, err)
		productLog.AssertCalled(t, "StoreProductLog", mock.Anything, mock.Anything)
	})

	t.Run("calculates correct current stock", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Use product with different initial stock
		initialStock := 100
		addedStock := 25
		expectedFinalStock := initialStock + addedStock

		productWithStock := &product.Product{
			ID:        1,
			SKU:       "SKU-002",
			Stock:     initialStock,
			CreatedAt: time.Now(),
		}

		payload := product.StockInPayload{
			StaffUUID: "staff-uuid-123",
			StockInReq: product.StockInReq{
				ProductSKU: "SKU-002",
				ProductQty: addedStock,
			},
		}

		// Mock user client to return user
		userClient.On("GetUser", mock.Anything, payload.StaffUUID).
			Return(&userDTO.User{UUID: payload.StaffUUID}, nil)

		// Mock product cache to return product
		productCache.On("GetProductBySKU", mock.Anything, payload.ProductSKU).
			Return(productWithStock, nil)

		// Mock product repo to update stock
		productRepo.On("UpdateProductStock", mock.Anything, payload.ProductSKU, payload.ProductQty, payload.StaffUUID).
			Return(nil)

		// Mock product cache to delete cache
		productCache.On("DelProductCache", mock.Anything, payload.ProductSKU).
			Return(nil)

		// Verify correct stock calculation in log
		productLog.On("StoreProductLog", mock.Anything, mock.MatchedBy(func(log product.ProductLog) bool {
			return log.CurrentStock == expectedFinalStock
		})).Return(nil)

		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			time.Duration(cfg.ContextTimeout)*time.Second,
		)

		err := svc.StockIn(ctx, payload)

		assert.NoError(t, err)
		productLog.AssertCalled(t, "StoreProductLog", mock.Anything, mock.MatchedBy(func(log product.ProductLog) bool {
			return log.CurrentStock == expectedFinalStock
		}))
	})

	t.Run("respects context timeout", func(t *testing.T) {
		productCache := &productCacheMocks.ProductCache{}
		productLog := &productLogMocks.ProductLogger{}
		productRepo := &productRepoMocks.Product{}
		userClient := &mocks.UserClient{}

		// Mock user client to simulate timeout
		userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
			Return(nil, context.DeadlineExceeded)

		// Mock product log to store log
		productLog.On("StoreProductLog", mock.Anything, mock.Anything).Return(nil)

		shortTimeout := time.Millisecond * 10
		svc := productService.NewProductService(
			productCache,
			productLog,
			productRepo,
			userClient,
			cfg,
			shortTimeout,
		)

		err := svc.StockIn(ctx, mockStockInPayload)

		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
	})
}

// Benchmarks
func BenchmarkStockIn(b *testing.B) {
	productCache := &productCacheMocks.ProductCache{}
	productLog := &productLogMocks.ProductLogger{}
	productRepo := &productRepoMocks.Product{}
	userClient := &mocks.UserClient{}

	userClient.On("GetUser", mock.Anything, mockStockInPayload.StaffUUID).
		Return(&userDTO.User{UUID: mockStockInPayload.StaffUUID}, nil)

	productCache.On("GetProductBySKU", mock.Anything, mockStockInPayload.ProductSKU).
		Return(mockProduct, nil)

	productRepo.On("UpdateProductStock", mock.Anything, mockStockInPayload.ProductSKU, mockStockInPayload.ProductQty, mockStockInPayload.StaffUUID).
		Return(nil)

	productCache.On("DelProductCache", mock.Anything, mockStockInPayload.ProductSKU).
		Return(nil)

	productLog.On("StoreProductLog", mock.Anything, mock.Anything).Return(nil)

	svc := productService.NewProductService(
		productCache,
		productLog,
		productRepo,
		userClient,
		cfg,
		time.Duration(cfg.ContextTimeout)*time.Second,
	)

	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		svc.StockIn(ctx, mockStockInPayload)
	}
}
