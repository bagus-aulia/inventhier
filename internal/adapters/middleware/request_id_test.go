package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bagus-aulia/inventhier/internal/core/helpers"
	"github.com/stretchr/testify/assert"
)

// TestRequestIDMiddleware tests the request ID middleware
func TestRequestIDMiddleware(t *testing.T) {
	testCases := []struct {
		name              string
		headerRequestID   string
		expectedHeader    string
		shouldGenerateNew bool
		description       string
	}{
		{
			name:              "uses X-Request-ID header if provided",
			headerRequestID:   "custom-request-id-123",
			expectedHeader:    "custom-request-id-123",
			shouldGenerateNew: false,
			description:       "Should use provided request ID from header",
		},
		{
			name:              "generates new request ID if not provided",
			headerRequestID:   "",
			expectedHeader:    "", // We can't predict the UUID, but it should be set
			shouldGenerateNew: true,
			description:       "Should generate new UUID when no header provided",
		},
		{
			name:              "handles empty header value",
			headerRequestID:   "",
			expectedHeader:    "",
			shouldGenerateNew: true,
			description:       "Should generate UUID for empty header",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request ID is in context
				requestID := helpers.GetRequestID(r.Context())
				assert.NotEmpty(t, requestID, "Request ID should be in context")

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(requestID))
			})

			middleware := RequestIDMiddleware(handler)

			// Create request with optional header
			req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
			if tc.headerRequestID != "" {
				req.Header.Set("X-Request-ID", tc.headerRequestID)
			}

			rr := httptest.NewRecorder()

			// Act
			middleware.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, http.StatusOK, rr.Code)

			responseRequestID := rr.Header().Get("X-Request-ID")
			assert.NotEmpty(t, responseRequestID, "Response should have X-Request-ID header")

			if !tc.shouldGenerateNew {
				assert.Equal(t, tc.expectedHeader, responseRequestID)
			}
		})
	}
}

// TestRequestIDMiddleware_ContextPropagation tests that request ID is properly propagated
func TestRequestIDMiddleware_ContextPropagation(t *testing.T) {
	t.Run("request ID is accessible in context throughout handler chain", func(t *testing.T) {
		// Arrange
		customRequestID := "test-req-propagation-123"
		receivedRequestID := ""

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract request ID from context
			receivedRequestID = helpers.GetRequestID(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		middleware := RequestIDMiddleware(handler)

		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		req.Header.Set("X-Request-ID", customRequestID)

		rr := httptest.NewRecorder()

		// Act
		middleware.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, customRequestID, receivedRequestID)
	})
}

// TestRequestIDMiddleware_HeaderInResponse tests that request ID is returned in response
func TestRequestIDMiddleware_HeaderInResponse(t *testing.T) {
	t.Run("response includes X-Request-ID header", func(t *testing.T) {
		// Arrange
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := RequestIDMiddleware(handler)

		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		rr := httptest.NewRecorder()

		// Act
		middleware.ServeHTTP(rr, req)

		// Assert
		responseRequestID := rr.Header().Get("X-Request-ID")
		assert.NotEmpty(t, responseRequestID, "Response should include X-Request-ID header")
	})
}

// TestRequestIDMiddleware_MultipleRequests tests that each request gets unique ID
func TestRequestIDMiddleware_MultipleRequests(t *testing.T) {
	t.Run("multiple requests get different request IDs", func(t *testing.T) {
		// Arrange
		requestIDs := make([]string, 3)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		middleware := RequestIDMiddleware(handler)

		// Act - Create 3 requests without X-Request-ID header
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
			rr := httptest.NewRecorder()
			middleware.ServeHTTP(rr, req)
			requestIDs[i] = rr.Header().Get("X-Request-ID")
		}

		// Assert - All request IDs should be unique
		assert.NotEmpty(t, requestIDs[0])
		assert.NotEmpty(t, requestIDs[1])
		assert.NotEmpty(t, requestIDs[2])
		assert.NotEqual(t, requestIDs[0], requestIDs[1])
		assert.NotEqual(t, requestIDs[1], requestIDs[2])
		assert.NotEqual(t, requestIDs[0], requestIDs[2])
	})
}

// BenchmarkRequestIDMiddleware benchmarks the middleware
func BenchmarkRequestIDMiddleware(b *testing.B) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestIDMiddleware(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		rr := httptest.NewRecorder()
		middleware.ServeHTTP(rr, req)
	}
}
