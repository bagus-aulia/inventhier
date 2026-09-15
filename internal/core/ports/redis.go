package ports

import (
	"context"
	"time"
)

// RedisInterface interface for Redis Repository
type RedisInterface interface {
	GetRedisData(ctx context.Context, key string, destination interface{}) error
	GetRedisString(ctx context.Context, key string) (string, error)
	SetRedisData(ctx context.Context, key string, data string, timeout time.Duration) error
	DelRedisData(ctx context.Context, keys []string) error
}
