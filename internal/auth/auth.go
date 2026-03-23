package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/p-v-dev/identity-service/internal/token"
	"github.com/p-v-dev/identity-service/internal/user"
)

// -- Interfaces consumed by this package --
// Defined here, not in the packages that implement them (see agents.md conventions).

// userRepository is the subset of persistence operations auth needs.
type userRepository interface {
	CreateUser(ctx context.Context, u user.User) (user.User, error)
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
}

// tokenService is the subset of token operations auth needs.
type tokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(tokenString string) (*token.Claims, error)
}

// cacheService is the subset of cache operations auth needs.
type cacheService interface {
	StoreRefreshToken(ctx context.Context, userID, refreshToken string, ttl time.Duration) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

// -- Request / Response types --

// RegisterRequest is the payload for POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the payload for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is returned after a successful login.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ValidateResponse is returned by GET /auth/validate.
// Other services use this to confirm identity without querying the DB.
type ValidateResponse struct {
	UserID string `json:"user_id"`
}

// -- Service --

// Service holds the auth business logic and its dependencies.
// All state lives outside the process (Postgres, Redis) — this struct is stateless.
type Service struct {
	users      userRepository
	tokens     tokenService
	cache      cacheService
	refreshTTL time.Duration // mirrors JWT_REFRESH_TTL; needed to set Redis key expiry
}

// NewService constructs a Service. Dependencies are injected; no globals.
func NewService(
	users userRepository,
	cache cacheService,
	tokens tokenService,
	refreshTTL time.Duration,
) *Service {
	return &Service{
		users:      users,
		tokens:     tokens,
		cache:      cache,
		refreshTTL: refreshTTL,
	}
}

// Register validates input, hashes the password, and persists a new user.
// Returns the created user on success (password field is always empty in the response).
func (s *Service) Register(ctx context.Context, req RegisterRequest) (user.User, error) {
	// TODO: Validate req.Email (format check) and req.Password (minimum length/strength).
	// TODO: Hash req.Password with bcrypt — hint: bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// TODO: Build a user.User{Email: req.Email, Password: hashedPassword} and call s.users.CreateUser.
	// TODO: On duplicate email, return a clear sentinel error (e.g. ErrEmailTaken) so the handler can respond 409.
	_ = fmt.Errorf // silence import until implemented
	return user.User{}, nil
}

// Login verifies credentials and issues a short-lived access token + long-lived refresh token.
// The refresh token is persisted in Redis so it can be revoked.
func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	// TODO: Fetch user via s.users.GetUserByEmail — return ErrInvalidCredentials on not-found (do not leak existence).
	// TODO: Compare req.Password against stored hash — hint: bcrypt.CompareHashAndPassword(...)
	// TODO: Generate access token  — s.tokens.GenerateAccessToken(user.ID)
	// TODO: Generate refresh token — s.tokens.GenerateRefreshToken(user.ID)
	// TODO: Persist refresh token  — s.cache.StoreRefreshToken(ctx, user.ID, refreshToken, s.refreshTTL)
	// TODO: Return LoginResponse{AccessToken: ..., RefreshToken: ...}
	return LoginResponse{}, nil
}

// ValidateToken checks the JWT signature and returns the embedded user identity.
// Called by GET /auth/validate — the endpoint other services rely on.
func (s *Service) ValidateToken(ctx context.Context, tokenString string) (ValidateResponse, error) {
	// TODO: Call s.tokens.ValidateToken(tokenString) to parse and verify the signature.
	// TODO: Map the returned *token.Claims to ValidateResponse{UserID: claims.UserID}.
	// TODO: Optionally check the blacklist via cacheService if logout support is added.
	return ValidateResponse{}, nil
}
