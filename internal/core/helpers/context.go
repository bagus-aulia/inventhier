package helpers

import (
	"context"

	"github.com/bagus-aulia/inventhier/internal/core/constants"
)

// GetRequestID extracts request ID from context
// Returns the request ID if available, otherwise returns empty string
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	requestID, ok := ctx.Value(constants.XRequestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, constants.XRequestIDKey, requestID)
}
