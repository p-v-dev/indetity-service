package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the JWT claims for the authentication service.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Service defines the interface for token generation and validation.
type Service interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// jwtService implements the TokenService interface using JWT.
type jwtService struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewService creates a new JWT token service.
func NewService(secret string, accessTTL, refreshTTL time.Duration) Service {
	return &jwtService{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// GenerateAccessToken generates a new access token.
func (s *jwtService) GenerateAccessToken(userID string) (string, error) {
	// TODO: Implement access token generation logic.
	// - Create Claims with UserID and appropriate RegisteredClaims (e.g., ExpiresAt, IssuedAt).
	// - Sign the token with the secret.
	// - Return the signed token string.
	return "", nil
}

// GenerateRefreshToken generates a new refresh token.
func (s *jwtService) GenerateRefreshToken(userID string) (string, error) {
	// TODO: Implement refresh token generation logic.
	// - Create Claims with UserID and appropriate RegisteredClaims (e.g., ExpiresAt, IssuedAt).
	// - Sign the token with the secret.
	// - Return the signed token string.
	return "", nil
}

// ValidateToken validates a given token string and returns its claims.
func (s *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	// TODO: Implement token validation logic.
	// - Parse the token string.
	// - Validate the signature using the secret.
	// - Extract and return the Claims.
	// - Handle potential errors (e.g., invalid token, expired token).
	return nil, nil
}
