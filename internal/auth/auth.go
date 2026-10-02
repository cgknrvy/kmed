package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"kmed/api/ent"
	"kmed/api/internal/httpx"
)

type Handler struct {
	svc            Service
	authMiddleware Middleware
	sessions       SessionRepo
}

func NewHandler(client *ent.Client, authMiddleware Middleware, secretKeyHex string) *Handler {
	svc, err := newService(client, secretKeyHex)
	if err != nil {
		panic(err)
	}
	sessions := SessionRepo{client: client}
	return &Handler{svc: svc, authMiddleware: authMiddleware, sessions: sessions}
}

func (h *Handler) Router() *http.ServeMux {
	authRouter := http.NewServeMux()
	authRouter.HandleFunc("POST /login", h.login)
	authRouter.HandleFunc("POST /refresh", h.refresh)
	authRouter.Handle("POST /logout", h.authMiddleware.Authenticate(http.HandlerFunc(h.logout)))

	return authRouter
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r LoginRequest) Validate() error {
	return httpx.Validator.Struct(r)
}

const (
	RefreshTokenCookieName = "refresh_token"
	RefreshCookiePath      = "/api/v1/auth/refresh"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	// Get the credentials
	credentials, ok := httpx.Parse[LoginRequest](w, r)
	if !ok {
		return
	}
	// Validate the credentials
	user, err := h.svc.validateCredentials(credentials)
	if err != nil {
		httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{
			Message: "invalid credentials",
		})
		return
	}

	h.issueSession(r.Context(), w, user)
}

func (h *Handler) issueSession(ctx context.Context, w http.ResponseWriter, user *ent.User) {
	scope := FullScope
	if user.MustChangePassword {
		scope = PasswordChangeScope
	}

	tokens := h.svc.generateTokens(&UserClaims{
		ID: user.ID, Email: user.Email, Role: user.Role, Scope: scope,
	})

	hash := hashRefreshToken(tokens.RefreshToken)
	if err := h.sessions.Create(ctx, user.ID, hash, time.Now().Add(RefreshTTL)); err != nil {
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: fmt.Sprintf("could not persist session: %s", err.Error())},
		)
		return
	}

	// Set refresh token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    tokens.RefreshToken,
		Path:     RefreshCookiePath,
		HttpOnly: true,
		Secure:   false, // false for local HTTP development
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(RefreshTTL),
	})

	httpx.JSON(w, http.StatusOK, httpx.Response{
		Message:            ptr("logged in"),
		AccessToken:        &tokens.AccessToken,
		MustChangePassword: ptr(user.MustChangePassword),
	})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{Message: "unauthorized"})
		return
	}

	refreshToken := cookie.Value

	t, err := h.sessions.FindByHash(r.Context(), hashRefreshToken(refreshToken))
	if err != nil {
		unsetCookie(w)
		httpx.JSONError(
			w,
			http.StatusUnauthorized,
			httpx.ErrorResponse{Message: "invalid refresh token"},
		)
		return
	}

	if t.RevokedAt != nil {
		h.sessions.RevokeAllForUser(r.Context(), t.UserID)
		unsetCookie(w)
		httpx.JSONError(
			w,
			http.StatusUnauthorized,
			httpx.ErrorResponse{Message: "refresh token already used"},
		)
		return
	}

	if time.Now().After(t.ExpiresAt) {
		unsetCookie(w)
		httpx.JSONError(
			w,
			http.StatusUnauthorized,
			httpx.ErrorResponse{Message: "refresh token expired"},
		)
		return
	}

	// Parse and validate the refresh token
	_, err = h.svc.parseToken(refreshToken)
	if err != nil {
		unsetCookie(w)
		httpx.JSONError(
			w,
			http.StatusUnauthorized,
			httpx.ErrorResponse{Message: "invalid refresh token"},
		)
		return
	}

	user, err := h.svc.getUser(t.UserID)
	if err != nil {
		unsetCookie(w)
		httpx.JSONError(
			w,
			http.StatusUnauthorized,
			httpx.ErrorResponse{Message: "user not found"},
		)
		return
	}

	// Rotate: revoke the old one before issuing a new session.
	if err := h.sessions.Revoke(r.Context(), t.ID); err != nil {
		unsetCookie(w)
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: "failed to revoke session"},
		)
		return
	}

	h.issueSession(r.Context(), w, user)
}

func unsetCookie(w http.ResponseWriter) {
	// Unset the cookie in the client
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     RefreshTokenCookieName,
			Value:    "",
			Path:     RefreshCookiePath,
			HttpOnly: true,
			Secure:   false, // false for local HTTP development
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		},
	)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := r.Context().Value(UserContextKey).(UserClaims)
	if ok {
		// Delete all user sessions once they logout
		// This may make other logged in instances to be logged out as well
		h.sessions.DeleteAllForUser(r.Context(), claims.ID)
	}

	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err == nil && cookie.Value != "" {
		if t, err := h.sessions.FindByHash(
			r.Context(),
			hashRefreshToken(cookie.Value),
		); err == nil {
			h.sessions.Revoke(r.Context(), t.ID)
		}
	}

	// Unset the refresh token cookie
	unsetCookie(w)
	httpx.JSON(w, http.StatusOK, httpx.Response{
		Message: ptr("logged out"),
	})
}

func ptr[T any](v T) *T {
	return &v
}
