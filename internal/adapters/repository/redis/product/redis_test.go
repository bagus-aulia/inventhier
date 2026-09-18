package product_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	"github.com/bagus-aulia/inventhier/internal/adapters/helpers/redis/mocks"
	"github.com/bagus-aulia/inventhier/internal/adapters/repository/redis/product"
	sqlProductMocks "github.com/bagus-aulia/inventhier/internal/adapters/repository/sql/product/mocks"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	cfg = &config.Config{
		RedisDefaultTimeout: 10,
	}

	mockProduct = &dto.Product{
		ID:        1,
		SKU:       "SKU-001",
		Name:      "Test Product",
		Price:     100000,
		Supplier:  "Supplier A",
		Stock:     50,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	redisKey = "product:sku:SKU-001"
)

// TestNewRedisCache tests the constructor
func TestNewRedisCache(t *testing.T) {
	t.Run("creates new redis cache instance", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		assert.NotNil(t, cache)
	})
}

// TestGetProductBySKU tests the GetProductBySKU method
func TestGetProductBySKU(t *testing.T) {
	ctx := context.Background()

	t.Run("returns product from cache when data exists in redis", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		// Mock redis to return cached product
		productJSON, _ := json.Marshal(mockProduct)
		redisClient.On("GetRedisData", ctx, redisKey, mock.MatchedBy(func(dest interface{}) bool {
			return true
		})).Run(func(args mock.Arguments) {
			dest := args.Get(2)
			json.Unmarshal(productJSON, dest)
		}).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		result, err := cache.GetProductBySKU(ctx, "SKU-001")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockProduct.SKU, result.SKU)
		assert.Equal(t, mockProduct.Name, result.Name)

		// Verify redis was called
		redisClient.AssertCalled(t, "GetRedisData", ctx, redisKey, mock.Anything)
		// Verify productRepo was NOT called since cache hit
		productRepo.AssertNotCalled(t, "GetProductBySKU")
	})

	t.Run("fetches from database and sets cache when redis cache misses", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		// Mock redis to return error (cache miss)
		redisClient.On("GetRedisData", ctx, redisKey, mock.MatchedBy(func(dest interface{}) bool {
			return true
		})).Return(errors.New("cache miss"))

		// Mock productRepo to return product
		productRepo.On("GetProductBySKU", ctx, "SKU-001").Return(mockProduct, nil)

		// Mock redis SET to cache the data
		productJSON, _ := json.Marshal(mockProduct)
		redisClient.On("SetRedisData", ctx, redisKey, string(productJSON), time.Duration(cfg.RedisDefaultTimeout)*time.Second).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		result, err := cache.GetProductBySKU(ctx, "SKU-001")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockProduct.SKU, result.SKU)

		// Verify both redis and productRepo were called
		redisClient.AssertCalled(t, "GetRedisData", ctx, redisKey, mock.Anything)
		productRepo.AssertCalled(t, "GetProductBySKU", ctx, "SKU-001")
		redisClient.AssertCalled(t, "SetRedisData", ctx, redisKey, mock.Anything, time.Duration(cfg.RedisDefaultTimeout)*time.Second)
	})

	t.Run("returns error when database query fails", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		// Mock redis to return error (cache miss)
		redisClient.On("GetRedisData", ctx, redisKey, mock.MatchedBy(func(dest interface{}) bool {
			return true
		})).Return(errors.New("cache miss"))

		// Mock productRepo to return error
		dbErr := errors.New("database connection error")
		productRepo.On("GetProductBySKU", ctx, "SKU-001").Return(nil, dbErr)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		result, err := cache.GetProductBySKU(ctx, "SKU-001")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, dbErr, err)

		// Verify productRepo was called
		productRepo.AssertCalled(t, "GetProductBySKU", ctx, "SKU-001")
		// Verify SetRedisData was NOT called due to database error
		redisClient.AssertNotCalled(t, "SetRedisData")
	})

	t.Run("caches data even if redis SET fails", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		// Mock redis to return error (cache miss)
		redisClient.On("GetRedisData", ctx, redisKey, mock.MatchedBy(func(dest interface{}) bool {
			return true
		})).Return(errors.New("cache miss"))

		// Mock productRepo to return product
		productRepo.On("GetProductBySKU", ctx, "SKU-001").Return(mockProduct, nil)

		// Mock redis SET to fail
		productJSON, _ := json.Marshal(mockProduct)
		redisClient.On("SetRedisData", ctx, redisKey, string(productJSON), time.Duration(cfg.RedisDefaultTimeout)*time.Second).Return(errors.New("redis set failed"))

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		result, err := cache.GetProductBySKU(ctx, "SKU-001")

		// Still returns product from database despite redis SET failure
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, mockProduct.SKU, result.SKU)

		// Verify all calls were made
		productRepo.AssertCalled(t, "GetProductBySKU", ctx, "SKU-001")
		redisClient.AssertCalled(t, "SetRedisData", ctx, redisKey, mock.Anything, time.Duration(cfg.RedisDefaultTimeout)*time.Second)
	})

	t.Run("returns product with different SKU", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		otherProduct := &dto.Product{
			ID:       2,
			SKU:      "SKU-002",
			Name:     "Another Product",
			Price:    50000,
			Supplier: "Supplier B",
			Stock:    100,
		}
		otherRedisKey := "product:sku:SKU-002"

		// Mock redis to return cached product
		productJSON, _ := json.Marshal(otherProduct)
		redisClient.On("GetRedisData", ctx, otherRedisKey, mock.MatchedBy(func(dest interface{}) bool {
			return true
		})).Run(func(args mock.Arguments) {
			dest := args.Get(2)
			json.Unmarshal(productJSON, dest)
		}).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		result, err := cache.GetProductBySKU(ctx, "SKU-002")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, otherProduct.SKU, result.SKU)
		assert.Equal(t, otherProduct.Name, result.Name)
	})

	t.Run("handles empty SKU string", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		emptyRedisKey := "product:sku:"
		redisClient.On("GetRedisData", ctx, emptyRedisKey, mock.MatchedBy(func(dest interface{}) bool {
			return true
		})).Return(errors.New("not found"))

		productRepo.On("GetProductBySKU", ctx, "").Return(nil, errors.New("invalid SKU"))

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		result, err := cache.GetProductBySKU(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("handles special characters in SKU", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		specialSKU := "SKU-001-SPECIAL@#$"
		specialRedisKey := "product:sku:" + specialSKU
		specialProduct := &dto.Product{
			SKU:   specialSKU,
			Name:  "Special Product",
			Price: 75000,
		}

		productJSON, _ := json.Marshal(specialProduct)
		redisClient.On("GetRedisData", ctx, specialRedisKey, mock.MatchedBy(func(dest interface{}) bool {
			return true
		})).Run(func(args mock.Arguments) {
			dest := args.Get(2)
			json.Unmarshal(productJSON, dest)
		}).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		result, err := cache.GetProductBySKU(ctx, specialSKU)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, specialSKU, result.SKU)
	})
}

// TestDelProductCache tests the DelProductCache method
func TestDelProductCache(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully deletes product cache", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		// Mock redis delete to succeed
		redisClient.On("DelRedisData", ctx, redisKey).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		err := cache.DelProductCache(ctx, "SKU-001")

		assert.NoError(t, err)
		redisClient.AssertCalled(t, "DelRedisData", ctx, redisKey)
	})

	t.Run("returns error when redis delete fails", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		// Mock redis delete to fail
		deleteErr := errors.New("redis delete failed")
		redisClient.On("DelRedisData", ctx, redisKey).Return(deleteErr)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		err := cache.DelProductCache(ctx, "SKU-001")

		assert.Error(t, err)
		assert.Equal(t, deleteErr, err)
		redisClient.AssertCalled(t, "DelRedisData", ctx, redisKey)
	})

	t.Run("handles deletion of non-existent key", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		// Mock redis delete to succeed even if key doesn't exist
		nonExistentKey := "product:sku:NONEXISTENT"
		redisClient.On("DelRedisData", ctx, nonExistentKey).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		err := cache.DelProductCache(ctx, "NONEXISTENT")

		assert.NoError(t, err)
		redisClient.AssertCalled(t, "DelRedisData", ctx, nonExistentKey)
	})

	t.Run("deletes cache for different SKUs independently", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		sku1 := "SKU-001"
		sku2 := "SKU-002"
		key1 := "product:sku:" + sku1
		key2 := "product:sku:" + sku2

		// Mock redis delete for both keys
		redisClient.On("DelRedisData", ctx, key1).Return(nil)
		redisClient.On("DelRedisData", ctx, key2).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		err1 := cache.DelProductCache(ctx, sku1)
		err2 := cache.DelProductCache(ctx, sku2)

		assert.NoError(t, err1)
		assert.NoError(t, err2)

		// Verify both deletes were called with correct keys
		redisClient.AssertCalled(t, "DelRedisData", ctx, key1)
		redisClient.AssertCalled(t, "DelRedisData", ctx, key2)
	})

	t.Run("handles empty SKU string in delete", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		emptyRedisKey := "product:sku:"
		redisClient.On("DelRedisData", ctx, emptyRedisKey).Return(nil)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		err := cache.DelProductCache(ctx, "")

		assert.NoError(t, err)
		redisClient.AssertCalled(t, "DelRedisData", ctx, emptyRedisKey)
	})

	t.Run("connection error during delete", func(t *testing.T) {
		redisClient := &mocks.RedisInterface{}
		productRepo := &sqlProductMocks.Product{}

		connErr := errors.New("connection timeout")
		redisClient.On("DelRedisData", ctx, redisKey).Return(connErr)

		cache := product.NewRedisCache(redisClient, productRepo, cfg)

		err := cache.DelProductCache(ctx, "SKU-001")

		assert.Error(t, err)
		assert.Equal(t, connErr, err)
	})
}

// Benchmarks
func BenchmarkGetProductBySKU_CacheHit(b *testing.B) {
	redisClient := &mocks.RedisInterface{}
	productRepo := &sqlProductMocks.Product{}
	ctx := context.Background()

	productJSON, _ := json.Marshal(mockProduct)
	redisClient.On("GetRedisData", ctx, redisKey, mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(2)
		json.Unmarshal(productJSON, dest)
	}).Return(nil)

	cache := product.NewRedisCache(redisClient, productRepo, cfg)

	b.ResetTimer()
	for b.Loop() {
		cache.GetProductBySKU(ctx, "SKU-001")
	}
}

func BenchmarkGetProductBySKU_CacheMiss(b *testing.B) {
	redisClient := &mocks.RedisInterface{}
	productRepo := &sqlProductMocks.Product{}
	ctx := context.Background()

	redisClient.On("GetRedisData", ctx, redisKey, mock.Anything).Return(errors.New("cache miss"))
	productRepo.On("GetProductBySKU", ctx, "SKU-001").Return(mockProduct, nil)

	productJSON, _ := json.Marshal(mockProduct)
	redisClient.On("SetRedisData", ctx, redisKey, string(productJSON), time.Duration(cfg.RedisDefaultTimeout)*time.Second).Return(nil)

	cache := product.NewRedisCache(redisClient, productRepo, cfg)

	b.ResetTimer()
	for b.Loop() {
		cache.GetProductBySKU(ctx, "SKU-001")
	}
}

func BenchmarkDelProductCache(b *testing.B) {
	redisClient := &mocks.RedisInterface{}
	productRepo := &sqlProductMocks.Product{}
	ctx := context.Background()

	redisClient.On("DelRedisData", ctx, redisKey).Return(nil)

	cache := product.NewRedisCache(redisClient, productRepo, cfg)

	b.ResetTimer()
	for b.Loop() {
		cache.DelProductCache(ctx, "SKU-001")
	}
}
