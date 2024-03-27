package interfaces

import (
	"github.com/coderero/erochat-server/api/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	// Generate generates a new token.
	GenerateTokens(email string, userId uuid.UUID) (string, string, error)

	// ValidateToken validates a token.
	ValidateToken(tokenString string) (bool, error)

	// GetClaims gets the claims from a token.
	GetClaims(tokenString string) (jwt.MapClaims, error)

	// GenerateToken generates a token.
	GenerateToken(email string, userId uuid.UUID, tokenType service.TokenType) (string, error)

	// RefreshToken refreshes a token.
	RefreshToken(refreshToken string) (string, error)
}
