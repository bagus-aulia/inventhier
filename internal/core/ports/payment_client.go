package ports

import (
	"context"

	"github.com/bagus-aulia/go-tools/tools/handler/response"
	dto "github.com/bagus-aulia/inventhier/internal/core/dto/payment"
)

// PaymentClient is a driven port defining the contract for communicating with the payment service.
// It can be implemented via gRPC or REST - the implementation detail is hidden from business logic.
type PaymentClient interface {
	Checkout(ctx context.Context, payload dto.CheckoutPayload) (response.Response[*dto.CheckoutResp], error)
}
