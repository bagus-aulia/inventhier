package helpers

import (
	"context"

	"github.com/bagus-aulia/inventhier/internal/core/constants"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GetZerologWithContext returns a zerolog logger with context values (request ID and user agent)
func GetZerologWithContext(ctx context.Context) zerolog.Logger {
	logger := log.Logger

	if ctx != nil {
		// Extract request ID from context
		if xRequestID := GetRequestID(ctx); xRequestID != "" {
			logger = logger.With().
				Interface(constants.XRequestIDHeader, xRequestID).
				Logger()
		}
	}

	return logger
}
