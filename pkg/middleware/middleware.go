package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/p-v-dev/identity-service/internal/token"
)

// contextKey is an unexported type for context keys in this package.
// Using a custom type avoids collisions with keys from other packages.
type contextKey string

const claimsKey contextKey = "claims"

// Claims holds the validated token data injected into the request context.
// It is a deliberately thin struct — the middleware only exposes what HTTP
// handlers need, not the full JWT internals from the token package.
type Claims struct {
	UserID string
}

// tokenValidator is the interface this middleware needs from the token package.
// Defined here (in the consumer) per project convention — small, focused interfaces.
// It matches token.Service so the concrete *token.jwtService satisfies it directly.
type tokenValidator interface {
	ValidateToken(tokenString string) (*token.Claims, error)
}

// Authenticate returns a chi-compatible middleware that:
//  1. Reads the "Authorization: Bearer <token>" header.
//  2. Validates the token via the provided tokenValidator.
//  3. Maps the result to a middleware.Claims and injects it into the request context.
//  4. Responds 401 Unauthorized if the header is missing, malformed, or the token is invalid.
//
// Usage with chi:
//
//	r.With(middleware.Authenticate(tokenSvc)).Get("/auth/validate", handler.Validate)
func Authenticate(v tokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := extractBearerToken(r)
			if err != nil {
				http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
				return
			}

			tokenClaims, err := v.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Map token.Claims → middleware.Claims.
			// This keeps downstream handlers decoupled from the token package internals.
			ctx := withClaims(r.Context(), Claims{UserID: tokenClaims.UserID})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext retrieves the Claims injected by Authenticate from the context.
// Returns false if no claims are present (e.g. the middleware was not applied to the route).
func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey).(Claims)
	return c, ok
}

// withClaims stores Claims in the context under the package-private key.
func withClaims(ctx context.Context, c Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

// extractBearerToken reads the Authorization header and strips the "Bearer " prefix.
// Returns a sentinel error if the header is absent or does not follow the expected format.
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errMissingHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", errMalformedHeader
	}

	tok := strings.TrimSpace(parts[1])
	if tok == "" {
		return "", errMalformedHeader
	}

	return tok, nil
}

// Sentinel errors — unexported because callers only care about the HTTP status code.
var (
	errMissingHeader   = sentinelError("missing Authorization header")
	errMalformedHeader = sentinelError("Authorization header must be: Bearer <token>")
)

type sentinelError string

func (e sentinelError) Error() string { return string(e) }
