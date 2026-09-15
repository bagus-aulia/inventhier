package auth_test

import (
	"testing"

	"github.com/bagus-aulia/inventhier/internal/adapters/helpers/auth"
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
		header      string
		expected    string
		expectError bool
		description string
	}{
		{
			name:        "extracts valid token from header",
			header:      "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expectError: false,
			description: "Should extract token from valid Bearer header",
		},
		{
			name:        "handles empty header",
			header:      "",
			expected:    "",
			expectError: true,
			description: "Should return error for empty header",
		},
		{
			name:        "handles invalid format without Bearer scheme",
			header:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expected:    "",
			expectError: true,
			description: "Should return error when Bearer scheme is missing",
		},
		{
			name:        "handles invalid scheme",
			header:      "Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected:    "",
			expectError: true,
			description: "Should return error for invalid scheme",
		},
		{
			name:        "handles empty token",
			header:      "Bearer ",
			expected:    "",
			expectError: true,
			description: "Should return error for empty token",
		},
		{
			name:        "handles extra whitespace",
			header:      "Bearer   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expectError: false,
			description: "Should handle extra whitespace",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := auth.ExtractTokenFromHeader(tc.header)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, token)
			}
		})
	}
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

	t.Run("returns error for empty token", func(t *testing.T) {
		_, err := auth.ParseJWT("")
		assert.Error(t, err)
	})

	t.Run("returns error for invalid token format", func(t *testing.T) {
		_, err := auth.ParseJWT("invalid.token")
		assert.Error(t, err)
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
	})
}

// BenchmarkParseJWT benchmarks JWT parsing
func BenchmarkParseJWT(b *testing.B) {
	claims := jwt.MapClaims{
		"sub":   "user123",
		"email": "test@example.com",
	}

	tokenString := createTestToken("secret", claims, jwt.SigningMethodHS256)

	for b.Loop() {
		_, _ = auth.ParseJWT(tokenString)
	}
}
