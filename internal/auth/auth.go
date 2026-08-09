package auth

import (
	"kmed/api/ent"
	"kmed/api/internal/httpx"
	"net/http"
)

type Handler struct {
	svc Service
}

// TODO: Needs to be securely stored outside and read
const secretKeyHex = "ab0f159045e95df5f9fc84194a6662cb0fc589c3c9fb323d876275a3d25a7525"

func NewHandler(client *ent.Client) *Handler {
	svc, err := newService(client, secretKeyHex)
	if err != nil {
		panic(err)
	}
	return &Handler{svc: svc}
}

func (h *Handler) Router() *http.ServeMux {
	authMiddleware := NewMiddleware()

	authRouter := http.NewServeMux()
	authRouter.HandleFunc("POST /login", h.login)
	authRouter.HandleFunc("POST /refresh", h.refresh)
	authRouter.Handle("POST /logout", authMiddleware.Authenticate(http.HandlerFunc(h.logout)))

	return authRouter
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r LoginRequest) Validate() error {
	return httpx.Validator.Struct(r)
}

const (
	RefreshTokenCookieName = "refresh_token"
	RefreshCookiePath      = "/v1/auth/refresh"
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
	// Generate tokens
	tokens := h.svc.generateTokens(&UserClaims{
		ID: user.ID, Email: user.Email, Role: user.Role,
	})

	//TODO: Create session

	// Set refresh token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    tokens.RefreshToken,
		Path:     RefreshCookiePath,
		HttpOnly: true,
		Secure:   false, // false for local HTTP development
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24,
	})

	// access token returned in the JSON response
	httpx.JSON(w, http.StatusOK, httpx.Response{
		Message:     ptr("logged in"),
		AccessToken: ptr(tokens.AccessToken),
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
	// Parse and validate the refresh token
	userClaims, err := h.svc.parseToken(refreshToken)
	if err != nil {
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
		httpx.JSONError(w, http.StatusUnauthorized, httpx.ErrorResponse{Message: "invalid refresh token"})
		return
	}

	// TODO: Validate the refresh token against the session - preferably a hash of the refresh token should be stored

	accessToken := h.svc.generateAccessToken(userClaims)

	httpx.JSON(w, http.StatusOK, httpx.Response{
		AccessToken: ptr(accessToken),
	})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	_, _ = r.Context().Value(UserContextKey).(UserClaims)

	//TODO: Delete active session for given user and revoke all tokens

	// Unset the refresh token cookie
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
	httpx.JSON(w, http.StatusOK, httpx.Response{
		Message: ptr("logged out"),
	})
}

func ptr[T any](v T) *T {
	return &v
}
