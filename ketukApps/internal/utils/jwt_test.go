package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	// Set secret for testing
	SetJWTSecret("test-secret-key")

	t.Run("GenerateToken", func(t *testing.T) {
		token, err := GenerateToken(1, "test@example.com", "user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate generated token
		claims, err := ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint(1), claims.UserID)
		assert.Equal(t, "test@example.com", claims.Email)
		assert.Equal(t, "user", claims.Role)
	})

	t.Run("GenerateRefreshToken", func(t *testing.T) {
		token, err := GenerateRefreshToken(1, "test@example.com", "user")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate refresh token
		claims, err := ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint(1), claims.UserID)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		claims, err := ValidateToken("invalid.token.string")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		// Manually create an expired token
		claims := JWTClaims{
			UserID: 1,
			Email:  "test@example.com",
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)

		_, err := ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has invalid claims")
	})
}
