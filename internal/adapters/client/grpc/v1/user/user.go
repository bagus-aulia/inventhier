package user

import (
	"context"
	"errors"

	grpc_pb "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/v1/user/pb"
	"github.com/bagus-aulia/inventhier/internal/core/constants"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/user"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/bagus-aulia/inventhier/internal/core/ports"
	"google.golang.org/grpc"
)

type grpcUserClient struct {
	grpcClient *grpc.ClientConn
}

// NewGRPCUserClient creates a new user gRPC client that implements the UserClient port.
func NewGRPCUserClient(grpcClient *grpc.ClientConn) ports.UserClient {
	return &grpcUserClient{
		grpcClient: grpcClient,
	}
}

// GetUser retrieves a user by UUID from the user service via gRPC.
func (c *grpcUserClient) GetUser(ctx context.Context, uuid string) (*dto.User, error) {
	logger := helpers.GetZerologWithContext(ctx).
		With().
		Str("client", "grpc.user.v1").
		Str("function", "GetUser").
		Str("uuid", uuid).
		Logger()

	userConn := grpc_pb.NewUserV1HandlerClient(c.grpcClient)
	resp, err := userConn.GetUser(ctx, &grpc_pb.Params{
		Uuid: uuid,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get user")

		return nil, err
	}

	if resp.Error != nil {
		err = constants.ErrInternalServer
		if resp.Error.Reason != "" {
			err = errors.New(resp.Error.Reason)
		}

		logger.Error().
			Err(err).
			Msg("There is an error on get user service")

		return nil, err
	}

	if resp.Data == nil {
		return nil, constants.ErrNotFound
	}

	return &dto.User{
		UUID:     resp.Data.Uuid,
		Name:     resp.Data.Name,
		Avatar:   resp.Data.Avatar,
		Position: resp.Data.Position,
	}, nil
}
