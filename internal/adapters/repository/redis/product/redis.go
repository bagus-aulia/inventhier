package product

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/product"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
)

type redisCache struct {
	redisClient ports.RedisInterface
	productRepo ports.Product
	cfg         *config.Config
}

// NewRedisCache creates a new caching adapter using the go-redis client.
func NewRedisCache(
	redisClient ports.RedisInterface,
	productRepo ports.Product,
	cfg *config.Config,
) ports.ProductCache {
	return &redisCache{
		redisClient: redisClient,
		productRepo: productRepo,
		cfg:         cfg,
	}
}

func (c *redisCache) GetProductBySKU(ctx context.Context, sku string) (*dto.Product, error) {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("repository", "redis.product").
		Str("function", "GetProductBySKU").
		Str("sku", sku).
		Logger()

	var data dto.Product

	// Get Room Entity From REdis
	redisKey := "product:sku:" + sku
	redisTimeOut := time.Duration(c.cfg.RedisDefaultTimeout)
	err := c.redisClient.GetRedisData(ctx, redisKey, &data)
	if err == nil {
		return &data, nil
	}

	productData, err := c.productRepo.GetProductBySKU(ctx, sku)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get product data from sql")
		return nil, err
	}

	// set redis data
	dataJSON, _ := json.Marshal(productData)
	err = c.redisClient.SetRedisData(ctx, redisKey, string(dataJSON), redisTimeOut*time.Second)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to set redis cache with product data")
	}

	return productData, nil
}
