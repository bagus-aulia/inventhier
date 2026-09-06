package ports

import (
	"context"

	dto "github.com/bagus-aulia/inventhier/internal/core/dto/user"
)

// UserClient is a driven port defining the contract for communicating with the user service.
// It can be implemented via gRPC or REST - the implementation detail is hidden from business logic.
type UserClient interface {
	GetUser(ctx context.Context, uuid string) (*dto.User, error)
}
