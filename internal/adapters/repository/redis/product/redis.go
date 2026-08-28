package product

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bagus-aulia/inventhier/internal/core/domain"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new caching adapter using the go-redis client.
func NewRedisCache(client *redis.Client) ports.ProductCache {
	return &redisCache{
		client: client,
	}
}

func (c *redisCache) Get(ctx context.Context, id string) (*domain.Product, error) {
	val, err := c.client.Get(ctx, "product:"+id).Result()
	if err != nil {
		return nil, err
	}

	var product domain.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *redisCache) Set(ctx context.Context, product *domain.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	// Cache with TTL of 1 hour
	return c.client.Set(ctx, "product:"+product.ID, data, 1*time.Hour).Err()
}

func (c *redisCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, "product:"+id).Err()
}
