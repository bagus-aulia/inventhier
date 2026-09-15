package middleware

import (
	"net/http"

	"github.com/bagus-aulia/inventhier/internal/core/constants"
	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to the request context
// If X-Request-ID header is present, it uses that; otherwise generates a new UUID
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get request ID from X-Request-ID header
		requestID := r.Header.Get(constants.XRequestIDHeader)

		// If not provided, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to response header for client reference
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to context
		ctx := helpers.WithRequestID(r.Context(), requestID)

		// Continue with the next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
