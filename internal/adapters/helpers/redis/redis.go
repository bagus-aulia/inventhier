package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type conRedisRepository struct {
	client redis.UniversalClient
	locker *redislock.Client
	cfg    *config.Config
}

// NewRedisRepository is a publisher of Redis Repository
func NewRedisRepository(redisCli redis.UniversalClient, cfg *config.Config) ports.RedisInterface {
	lockerCli := redislock.New(redisCli)

	return &conRedisRepository{
		client: redisCli,
		locker: lockerCli,
		cfg:    cfg,
	}
}

// GetRedisData to get redis cache
func (r *conRedisRepository) GetRedisData(ctx context.Context, key string, destination interface{}) error {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("adapter.helpers", "redis").
		Str("function", "GetRedisData").
		Str("key", key).
		Logger()

	redisResult, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Warn().
				Err(err).
				Msg("Failed to get redis data")
		}
		return err
	}

	if redisResult != "" {
		err := json.Unmarshal([]byte(redisResult), destination)
		if err != nil {
			logger.Warn().
				Err(err).
				Msg("Failed to unmarshal redis data")
			return err
		}
	}

	return nil
}

// SetRedisData to store redis cache
func (r *conRedisRepository) SetRedisData(ctx context.Context, key string, data string, timeout time.Duration) error {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("adapter.helpers", "redis").
		Str("function", "SetRedisData").
		Str("key", key).
		Str("data", data).
		Logger()

	if timeout == 0 {
		timeout = time.Duration(r.cfg.RedisDefaultTimeout) * time.Second
	}

	err := r.client.Set(ctx, key, data, timeout).Err()
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Failed to set redis data")
		return err
	}

	return nil
}

// DelRedisData to delete redis cache
func (r *conRedisRepository) DelRedisData(ctx context.Context, key string) error {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("adapter.helpers", "redis").
		Str("function", "DelRedisData").
		Str("key", key).
		Logger()

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		logger.Warn().
			Err(err).
			Msg("Failed to delete redis key")
		return err
	}

	return nil
}
