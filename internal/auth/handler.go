package auth

import (
	"encoding/json"
	"errors"
	"net/http"
)

// contextKey is an unexported type for context keys in this package.
// This prevents collisions with keys defined in other packages.
type contextKey string

const (
	// ClaimsContextKey is the key used to store validated JWT claims in the request context.
	// The Authenticate middleware writes here; the Validate handler reads from here.
	ClaimsContextKey contextKey = "claims"
)

// Handler holds the HTTP handlers for all auth routes.
// It depends only on the Service — it has no knowledge of the database or cache.
type Handler struct {
	svc *Service
}

// NewHandler creates a new Handler wrapping the given Service.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register handles POST /auth/register
//
// Request body:  { "email": "...", "password": "..." }
// Success:       201 Created  — { "id": "...", "email": "..." }
// Errors:        400 Bad Request | 409 Conflict | 500 Internal Server Error
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.svc.Register(r.Context(), req)
	if err != nil {
		// TODO: Distinguish between validation errors (400), duplicate email (409),
		//       and unexpected errors (500) using sentinel errors or error types.
		writeError(w, http.StatusInternalServerError, "could not register user")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// Login handles POST /auth/login
//
// Request body:  { "email": "...", "password": "..." }
// Success:       200 OK — { "access_token": "...", "refresh_token": "..." }
// Errors:        400 Bad Request | 401 Unauthorized | 500 Internal Server Error
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		// TODO: Distinguish between wrong credentials (401) and unexpected errors (500).
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Validate handles GET /auth/validate
//
// The token is validated upstream by the Authenticate middleware, which stores
// the claims in the request context. This handler just reads and returns them.
//
// Header:   Authorization: Bearer <token>
// Success:  200 OK — { "user_id": "..." }
// Errors:   401 Unauthorized (emitted by the middleware, not this handler)
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	// TODO: Extract the ValidateResponse stored in context by pkg/middleware.Authenticate.
	//       Example: claims, ok := r.Context().Value(ClaimsContextKey).(ValidateResponse)
	//       Return 401 if the assertion fails (should not happen if middleware is wired correctly).
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "token valid", // placeholder until middleware wires the claims
	})
}

// -- Helpers -----------------------------------------------------------------

// errorResponse is the standard error body returned by all handlers.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON serialises v as JSON and writes it with the given HTTP status code.
// It sets Content-Type to application/json before writing the header so that
// clients always receive a properly typed response.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// At this point the header is already sent, so we can only log.
		// TODO: Replace with a structured logger once one is wired in.
		_ = err
	}
}

// writeError writes a JSON error response with the given status and message.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

// decodeBody is a convenience wrapper that decodes JSON from r.Body into dst
// and returns a descriptive error if decoding fails.
// Usage: if err := decodeBody(r, &req); err != nil { writeError(...) }
func decodeBody(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errors.New("malformed JSON body: " + err.Error())
	}
	return nil
}
