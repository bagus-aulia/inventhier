package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bagus-aulia/inventhier/internal/core/constants"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents standard JWT claims plus custom claims
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

// ExtractTokenFromHeader extracts JWT token from Authorization header
// Expected format: "Bearer <token>"
func ExtractTokenFromHeader(headers map[string][]string) (string, error) {
	authHeader := ""
	for _, value := range headers[constants.Authorization] {
		authHeader = value
		break
	}

	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	if parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization scheme: expected Bearer, got %s", parts[0])
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}

// ParseJWT parses JWT token without verification
// Use this only when you want to extract claims without validation
func ParseJWT(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token string is empty")
	}

	parser := jwt.NewParser(
		jwt.WithoutClaimsValidation(),
	)
	token, _, err := parser.ParseUnverified(
		tokenString,
		&JWTClaims{},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token claims type")
}

// ExtractUserID extracts user_id from JWT token
func ExtractUserID(tokenString string) (string, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	if claims.UserID == "" {
		return "", errors.New("user_id claim not found in token")
	}

	return claims.UserID, nil
}
