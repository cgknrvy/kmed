package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"kmed/api/ent"
)

func TestLogin(t *testing.T) {
	loginData := []LoginRequest{
		{
			Email:    "test@kmed.com",
			Password: "12345678",
		},
		{
			Email:    "test1@kmed.com",
			Password: "12345678",
		},
	}
	svc := &ServiceStub{
		users: []*ent.User{
			{
				Email:    loginData[0].Email,
				Password: loginData[0].Password,
			},
			{
				Email:    loginData[1].Email,
				Password: loginData[1].Password,
			},
		},
	}
	authHandler := &Handler{svc}

	t.Run("returns token after user login", func(t *testing.T) {
		request := newLoginRequest(loginData[0])
		response := httptest.NewRecorder()
		authHandler.Router().ServeHTTP(response, request)

		assert.Contains(t, svc.tokensGenerated, loginData[0].Email)
		assert.Contains(t, svc.credValidated, loginData[0].Email)
		assert.Equal(
			t,
			response.Code,
			http.StatusOK,
			"should successfully login",
		)
		cookies := response.Header().Values("Set-Cookie")
		assert.Equal(t, len(cookies), 1, "expected refresh token to be set")
		assert.Contains(t, cookies[0], RefreshTokenCookieName)
	})

	t.Run("returns error when invalid email", func(t *testing.T) {
		request := newLoginRequest(LoginRequest{Email: "one@email.com", Password: loginData[0].Password})
		response := httptest.NewRecorder()
		authHandler.Router().ServeHTTP(response, request)

		assert.NotContains(t, svc.tokensGenerated, "one@email.com")
		assert.NotContains(t, svc.tokensGenerated, "one@email.com")
		assert.Equal(
			t,
			response.Code,
			http.StatusUnauthorized,
			"login should fail with code 401",
		)
	})

	t.Run("returns error when invalid password", func(t *testing.T) {
		request := newLoginRequest(LoginRequest{Email: loginData[1].Email, Password: "invalid-password"})
		response := httptest.NewRecorder()
		authHandler.Router().ServeHTTP(response, request)

		assert.NotContains(t, svc.tokensGenerated, loginData[1].Email)
		assert.NotContains(t, svc.tokensGenerated, loginData[1].Email)
		assert.Equal(
			t,
			response.Code,
			http.StatusUnauthorized,
			"login should fail with code 401",
		)
	})
}

func newLoginRequest(payload LoginRequest) *http.Request {
	data, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", "/login", strings.NewReader(string(data)))
	request.Header.Set("Content-Type", "application/json")
	return request
}

type ServiceStub struct {
	users           []*ent.User
	tokensGenerated []string
	credValidated   []string
}

func (s *ServiceStub) generateAccessToken(user *UserClaims) string {
	return ""
}

func (s *ServiceStub) generateRefreshToken(user *UserClaims) string {
	return ""
}

func (s *ServiceStub) validateCredentials(credentials LoginRequest) (*ent.User, error) {
	if s.users[0].Email == credentials.Email && s.users[0].Password == credentials.Password {
		s.credValidated = append(s.credValidated, credentials.Email)
		return &ent.User{
			Email:    credentials.Email,
			Password: credentials.Password,
		}, nil
	}
	return &ent.User{}, errors.New("invalid credentials")
}

func (s *ServiceStub) generateTokens(user *UserClaims) *TokenResponse {
	s.tokensGenerated = append(s.tokensGenerated, user.Email)
	return &TokenResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}
}

func (s *ServiceStub) parseToken(token string) (*UserClaims, error) {
	return nil, nil
}
