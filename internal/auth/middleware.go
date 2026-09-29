package auth

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"kmed/api/ent/user"
	"kmed/api/internal/httpx"

	"aidanwoods.dev/go-paseto"
)

type ContextKey string

const UserContextKey ContextKey = "user"

type Middleware interface {
	Authenticate(next http.Handler) http.Handler
	RequireRole(next http.Handler, role []user.Role) http.Handler
	RequireUser(next http.Handler) http.Handler
	RequireAdmin(next http.Handler) http.Handler
	RequireLabTech(next http.Handler) http.Handler
	RequireDoctor(next http.Handler) http.Handler
	RequireAdminOrDoctor(next http.Handler) http.Handler
}

type middleware struct {
	svc tokenService
}

func NewMiddleware(secretKeyHex string) Middleware {
	secretKey, err := paseto.V4SymmetricKeyFromHex(secretKeyHex)
	if err != nil {
		panic(fmt.Errorf("failed to generate symmetric key: %w", err))
	}
	svc := newTokenService(secretKey)

	return &middleware{svc: svc}
}

func (m middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{
				Message: "missing authorization header",
			})
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{
				Message: "invalid authorization format",
			})
			return
		}
		tokenString := parts[1]

		// Parse token
		parsedToken, err := m.svc.parseToken(tokenString)
		if err != nil {
			httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{
				Message: "invalid token",
				Errors:  map[string]string{"parsing": err.Error()},
			})
			return
		}

		// Extract user claims. Already validated so no errors expected.
		userClaims := extractUserClaims(parsedToken)

		// TODO: Check if there is a session for given user. If session is expired
		// revoke the key and return an unauthorized error

		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			*userClaims,
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m middleware) RequireRole(next http.Handler, roles []user.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims, ok := r.Context().Value(UserContextKey).(UserClaims)
		if !ok {
			httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{
				Message: "unauthorized",
			})
			return
		}

		if !slices.Contains(roles, userClaims.Role) {
			httpx.JSONError(w, http.StatusForbidden, httpx.ErrorResponse{
				Message: "insufficient permissions",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m middleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireRole(next, []user.Role{user.RoleAdmin})
}

func (m middleware) RequireAdminOrDoctor(next http.Handler) http.Handler {
	return m.RequireRole(next, []user.Role{user.RoleAdmin, user.RoleDoctor})
}

func (m middleware) RequireDoctor(next http.Handler) http.Handler {
	return m.RequireRole(next, []user.Role{user.RoleDoctor})
}

func (m middleware) RequireLabTech(next http.Handler) http.Handler {
	return m.RequireRole(next, []user.Role{user.RoleLabTech})
}

func (m middleware) RequireUser(next http.Handler) http.Handler {
	return m.RequireRole(next, []user.Role{user.RoleUser})
}

// RequireFullScope checks for the scope of the user and only allows full scope
// to access the handler. The password_change scope cannot access all functionality.
func (m middleware) RequireFullScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userClaims, ok := r.Context().Value(UserContextKey).(UserClaims)
		if !ok {
			httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{
				Message: "unauthorized",
			})
			return
		}

		if userClaims.Scope != "full" {
			httpx.JSONError(w, http.StatusForbidden, httpx.ErrorResponse{
				Message: "password change required",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
