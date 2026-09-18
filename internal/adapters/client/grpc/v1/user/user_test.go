package user_test

import (
	"context"
	"errors"
	"log"
	"testing"

	userCli "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/v1/user"
	grpc_pb "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/v1/user/pb"
	grpc_pbMocks "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/v1/user/pb/mocks"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	t.Run("GetUser success", func(t *testing.T) {
		setup := func() *grpc_pbMocks.MockUserServer {
			mockServer := new(grpc_pbMocks.MockUserServer)
			mockServer.On("GetUser", mock.Anything, mock.Anything).Return(&grpc_pb.Data{
				Status: 200,
				Error:  nil,
				Data: &grpc_pb.User{
					Name:     "new-name",
					Uuid:     "example-uuid",
					Avatar:   "http://avatar.io/sample.tiff",
					Position: "cashier",
				},
			}, nil)

			return mockServer
		}

		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(grpc_pbMocks.Dialer(setup)))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		var uuid string
		err = faker.FakeData(&uuid)
		assert.NoError(t, err)

		userClient := userCli.NewGRPCUserClient(conn)
		result, errRes := userClient.GetUser(context.Background(), uuid)

		assert.NoError(t, errRes)
		assert.NotNil(t, result)
		assert.Equal(t, "new-name", result.Name)
		assert.Equal(t, "example-uuid", result.UUID)
	})

	t.Run("GetUser error response", func(t *testing.T) {
		setup := func() *grpc_pbMocks.MockUserServer {
			mockServer := new(grpc_pbMocks.MockUserServer)
			mockServer.On("GetUser", mock.Anything, mock.Anything).Return(&grpc_pb.Data{
				Status: 404,
				Error: &grpc_pb.Error{
					Message: "data not found",
					Reason:  "not_found",
				},
				Data: nil,
			}, errors.New("example error"))

			return mockServer
		}

		conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(grpc_pbMocks.Dialer(setup)))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		var uuid string
		err = faker.FakeData(&uuid)
		assert.NoError(t, err)

		userClient := userCli.NewGRPCUserClient(conn)
		result, errRes := userClient.GetUser(context.Background(), uuid)

		assert.Error(t, errRes)
		assert.Nil(t, result)
	})
}
