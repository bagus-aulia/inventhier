package redis_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/bagus-aulia/inventhier/config"
	redisHelper "github.com/bagus-aulia/inventhier/internal/adapters/helpers/redis"
	"github.com/bagus-aulia/inventhier/internal/core/constants"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

var (
	cfg = &config.Config{
		RedisDefaultTimeout: 10,
	}
)

func TestGetRedisData(t *testing.T) {
	ctx := context.TODO()
	key := "testKey"
	data := "testData"
	dataJSON, _ := json.Marshal(data)

	var destination string

	t.Run("success", func(t *testing.T) {
		// Set up mock redis connection
		clusterClient, clusterMock := redismock.NewClusterMock()
		repo := redisHelper.NewRedisRepository(clusterClient, cfg)

		// Set up the mock behavior
		clusterMock.Regexp().ExpectGet(key).SetVal(string(dataJSON))

		// Call the method

		err := repo.GetRedisData(ctx, key, &destination)

		// Assertions
		assert.NoError(t, err)
	})

	t.Run("error unmarshall", func(t *testing.T) {
		// Set up mock redis connection
		clusterClient, clusterMock := redismock.NewClusterMock()
		repo := redisHelper.NewRedisRepository(clusterClient, cfg)

		// Set up the mock behavior
		clusterMock.Regexp().ExpectGet(key).SetVal(data)

		// Call the method
		err := repo.GetRedisData(ctx, key, &destination)

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("error get redis", func(t *testing.T) {
		// Set up mock redis connection
		clusterClient, clusterMock := redismock.NewClusterMock()
		repo := redisHelper.NewRedisRepository(clusterClient, cfg)

		// Set up the mock behavior
		clusterMock.Regexp().ExpectGet(key).SetErr(constants.ErrRedisNil)

		// Call the method
		err := repo.GetRedisData(ctx, key, &destination)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, constants.ErrRedisNil, err)
	})
}

func TestSetRedisData(t *testing.T) {
	ctx := context.TODO()
	key := "testKey"
	data := "testData"
	timeout := 10 * time.Second

	t.Run("success", func(t *testing.T) {
		// Set up mock redis connection
		clusterClient, clusterMock := redismock.NewClusterMock()
		repo := redisHelper.NewRedisRepository(clusterClient, cfg)

		// Set up the mock behavior
		clusterMock.Regexp().ExpectSet(key, data, timeout).SetVal(data)

		// Call the method
		err := repo.SetRedisData(ctx, key, data, timeout)

		// Assertions
		assert.NoError(t, err)
	})

	t.Run("error set redis", func(t *testing.T) {
		// Set up mock redis connection
		clusterClient, clusterMock := redismock.NewClusterMock()
		repo := redisHelper.NewRedisRepository(clusterClient, cfg)

		// Set up the mock behavior
		clusterMock.Regexp().ExpectSet(key, data, timeout).SetErr(constants.ErrRedisNil)

		// Call the method
		err := repo.SetRedisData(ctx, key, data, timeout)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, constants.ErrRedisNil, err)
	})
}

func TestDelRedisData(t *testing.T) {
	ctx := context.TODO()
	key := "testKey"

	t.Run("success", func(t *testing.T) {
		// Set up mock redis connection
		clusterClient, clusterMock := redismock.NewClusterMock()
		repo := redisHelper.NewRedisRepository(clusterClient, cfg)

		// Set up the mock behavior
		clusterMock.Regexp().ExpectDel(key).SetVal(1)

		// Call the method
		err := repo.DelRedisData(ctx, key)

		// Assertions
		assert.NoError(t, err)
	})

	t.Run("error del redis", func(t *testing.T) {
		// Set up mock redis connection
		clusterClient, clusterMock := redismock.NewClusterMock()
		repo := redisHelper.NewRedisRepository(clusterClient, cfg)

		// Set up the mock behavior
		clusterMock.Regexp().ExpectDel(key).SetErr(constants.ErrRedisNil)

		// Call the method
		err := repo.DelRedisData(ctx, key)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, constants.ErrRedisNil, err)
	})
}
