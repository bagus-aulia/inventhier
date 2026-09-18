package ports

import (
	"context"
	"time"
)

// RedisInterface interface for Redis Helper
type RedisInterface interface {
	GetRedisData(ctx context.Context, key string, destination interface{}) error
	SetRedisData(ctx context.Context, key string, data string, timeout time.Duration) error
	DelRedisData(ctx context.Context, key string) error
}
