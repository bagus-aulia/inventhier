package user

import (
	"context"

	"google.golang.org/grpc"
)

// GRPCUserServiceClient wraps the generated UserV1HandlerClient for easier mocking
type GRPCUserServiceClient interface {
	GetUser(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Data, error)
}
