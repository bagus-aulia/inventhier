package v1_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	user_client "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/user/v1"
// 	grpc_pb "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/user/v1/pb"
// 	grpc_pb_mocks "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/user/v1/pb/mocks"
// 	"github.com/bagus-aulia/inventhier/internal/core/constants"
// 	dto "github.com/bagus-aulia/inventhier/internal/core/dto/user"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"google.golang.org/grpc"
// )

// func TestGetUser(t *testing.T) {
// 	cfg := config.LoadConfig()

// 	uuid := "user-1"

// 	t.Run("error resp grpc", func(t *testing.T) {
// 		// grpc mock conn
// 		grpcConn, err := grpc_pb_mocks.MockGrpcConn()
// 		assert.NoError(t, err)
// 		defer grpcConn.Close()

// 		mockRedisConn := new(redisConn_mock.RedisMock)

// 		mockRedisConn.On("GetRedisData", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(redis.Nil)

// 		u := liveConn.NewLiveService(grpcConn, mockRedisConn, cfg)

// 		res, err := u.GetLivestream(c, portalSlug, slug, uuid)
// 		assert.Error(t, err)
// 		assert.Nil(t, res)
// 	})

// 	t.Run("success redis", func(t *testing.T) {
// 		// grpc mock conn
// 		grpcConn, err := live_rpcMock.MockGrpcConn()
// 		assert.NoError(t, err)
// 		defer grpcConn.Close()

// 		mockRedisConn := new(redisConn_mock.RedisMock)

// 		mockRedisConn.On("GetRedisData", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(nil)

// 		u := liveConn.NewLiveService(grpcConn, mockRedisConn, cfg)

// 		res, err := u.GetLivestream(c, portalSlug, slug, uuid)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, res)
// 		assert.Equal(t, int32(200), res.Status)
// 	})

// 	t.Run("success", func(t *testing.T) {
// 		// grpc mock conn
// 		grpcConn, err := live_rpcMock.MockGrpcConn()
// 		assert.NoError(t, err)
// 		defer grpcConn.Close()

// 		mockRedisConn := new(redisConn_mock.RedisMock)

// 		mockRedisConn.On("GetRedisData", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(redis.Nil)

// 		mockRedisConn.On("SetRedisData", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)

// 		u := liveConn.NewLiveService(grpcConn, mockRedisConn, cfg)

// 		res, err := u.GetLivestream(c, portalSlug, slug, uuid)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, res)
// 	})
// }
