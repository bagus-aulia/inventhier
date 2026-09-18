package auth_test

import (
	"net/http"
	"testing"

	"github.com/bagus-aulia/inventhier/internal/adapters/helpers/auth"
	"github.com/bagus-aulia/inventhier/internal/core/constants"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test JWT token
func createTestToken(secret string, claims jwt.MapClaims, signingMethod jwt.SigningMethod) string {
	token := jwt.NewWithClaims(signingMethod, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// TestExtractTokenFromHeader tests extracting token from Authorization header
func TestExtractTokenFromHeader(t *testing.T) {
	testCases := []struct {
		name        string
		headers     map[string][]string
		expected    string
		expectError bool
		description string
	}{
		{
			name: "extracts valid token from header",
			headers: map[string][]string{
				constants.Authorization: {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
			},
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expectError: false,
			description: "Should extract token from valid Bearer header",
		},
		{
			name:        "handles empty header",
			headers:     map[string][]string{},
			expected:    "",
			expectError: true,
			description: "Should return error for empty header",
		},
		{
			name: "handles invalid format without Bearer scheme",
			headers: map[string][]string{
				constants.Authorization: {"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
			},
			expected:    "",
			expectError: true,
			description: "Should return error when Bearer scheme is missing",
		},
		{
			name: "handles invalid scheme",
			headers: map[string][]string{
				constants.Authorization: {"Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
			},
			expected:    "",
			expectError: true,
			description: "Should return error for invalid scheme",
		},
		{
			name: "handles empty token",
			headers: map[string][]string{
				constants.Authorization: {"Bearer "},
			},
			expected:    "",
			expectError: true,
			description: "Should return error for empty token",
		},
		{
			name: "handles extra whitespace",
			headers: map[string][]string{
				constants.Authorization: {"Bearer   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
			},
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expectError: false,
			description: "Should handle extra whitespace",
		},
		{
			name: "handles multiple authorization headers",
			headers: map[string][]string{
				constants.Authorization: {
					"Bearer token-123",
					"Bearer token-456", // Should use the first one
				},
			},
			expected:    "token-123",
			expectError: false,
			description: "Should use the first authorization header when multiple are present",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := auth.ExtractTokenFromHeader(tc.headers)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, token)
			}
		})
	}
}

// TestExtractTokenFromHeader_WithHTTPHeader tests with standard http.Header
func TestExtractTokenFromHeader_WithHTTPHeader(t *testing.T) {
	t.Run("extracts token from http.Header", func(t *testing.T) {
		header := http.Header{}
		header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U")

		token, err := auth.ExtractTokenFromHeader(header)

		assert.NoError(t, err)
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U", token)
	})
}

// TestParseJWT tests parsing JWT token without verification
func TestParseJWT(t *testing.T) {
	t.Run("parses valid JWT token", func(t *testing.T) {
		// Create a test token
		claims := jwt.MapClaims{
			"email":    "test@example.com",
			"user_id":  "user-id-123",
			"username": "testuser",
		}

		secret := "test-secret"
		tokenString := createTestToken(secret, claims, jwt.SigningMethodHS256)

		// Parse token
		parsedClaims, err := auth.ParseJWT(tokenString)

		assert.NoError(t, err)
		assert.NotNil(t, parsedClaims)
		assert.Equal(t, "test@example.com", parsedClaims.Email)
		assert.Equal(t, "user-id-123", parsedClaims.UserID)
		assert.Equal(t, "testuser", parsedClaims.Username)
	})

	t.Run("parses token with only user_id", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id": "user-789",
		}

		tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)
		parsedClaims, err := auth.ParseJWT(tokenString)

		assert.NoError(t, err)
		assert.Equal(t, "user-789", parsedClaims.UserID)
		assert.Empty(t, parsedClaims.Email)
	})

	t.Run("returns error for empty token", func(t *testing.T) {
		_, err := auth.ParseJWT("")
		assert.Error(t, err)
	})

	t.Run("returns error for invalid token format", func(t *testing.T) {
		_, err := auth.ParseJWT("invalid.token")
		assert.Error(t, err)
	})

	t.Run("returns error for malformed JWT", func(t *testing.T) {
		_, err := auth.ParseJWT("not-a-valid-jwt-at-all")
		assert.Error(t, err)
	})
}

// TestParseJWT_StandardClaims tests parsing standard JWT claims
func TestParseJWT_StandardClaims(t *testing.T) {
	t.Run("preserves standard claims", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":     "subject-123",
			"user_id": "user-123",
		}

		tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)
		parsedClaims, err := auth.ParseJWT(tokenString)

		assert.NoError(t, err)
		assert.Equal(t, "subject-123", parsedClaims.Subject)
		assert.Equal(t, "user-123", parsedClaims.UserID)
	})
}

// TestExtractUserID tests extracting user_id specifically
func TestExtractUserID(t *testing.T) {
	t.Run("extracts user_id successfully", func(t *testing.T) {
		claims := jwt.MapClaims{
			"user_id": "user-123",
		}

		tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)
		userID, err := auth.ExtractUserID(tokenString)

		assert.NoError(t, err)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("returns error when user_id is missing", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "user123",
		}

		tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)
		_, err := auth.ExtractUserID(tokenString)

		assert.Error(t, err)
		assert.Equal(t, "user_id claim not found in token", err.Error())
	})

	t.Run("returns error for empty token", func(t *testing.T) {
		_, err := auth.ExtractUserID("")
		assert.Error(t, err)
	})

	t.Run("returns error for invalid token", func(t *testing.T) {
		_, err := auth.ExtractUserID("invalid-token")
		assert.Error(t, err)
	})
}

// TestExtractUserID_VariousForms tests extracting user_id in various formats
func TestExtractUserID_VariousForms(t *testing.T) {
	testCases := []struct {
		name      string
		userID    string
		expectErr bool
	}{
		{"UUID format", "550e8400-e29b-41d4-a716-446655440000", false},
		{"Numeric string", "12345", false},
		{"Alphanumeric", "user_abc123", false},
		{"With special chars", "user@domain.com", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				"user_id": tc.userID,
			}

			tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)
			extractedID, err := auth.ExtractUserID(tokenString)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.userID, extractedID)
			}
		})
	}
}

// BenchmarkExtractTokenFromHeader benchmarks header extraction
func BenchmarkExtractTokenFromHeader(b *testing.B) {
	headers := map[string][]string{
		constants.Authorization: {"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"},
	}

	b.ResetTimer()
	for b.Loop() {
		_, _ = auth.ExtractTokenFromHeader(headers)
	}
}

// BenchmarkParseJWT benchmarks JWT parsing
func BenchmarkParseJWT(b *testing.B) {
	claims := jwt.MapClaims{
		"sub":   "user123",
		"email": "test@example.com",
	}

	tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)

	b.ResetTimer()
	for b.Loop() {
		_, _ = auth.ParseJWT(tokenString)
	}
}

// BenchmarkExtractUserID benchmarks user_id extraction
func BenchmarkExtractUserID(b *testing.B) {
	claims := jwt.MapClaims{
		"user_id": "user-123",
	}

	tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)

	b.ResetTimer()
	for b.Loop() {
		_, _ = auth.ExtractUserID(tokenString)
	}
}
