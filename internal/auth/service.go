package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"kmed/api/ent"
	"kmed/api/ent/user"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service interface {
	validateCredentials(credentials LoginRequest) (*ent.User, error)
	generateAccessToken(user *UserClaims) string
	generateRefreshToken(user *UserClaims) string
	generateTokens(user *UserClaims) *TokenResponse
	parseToken(token string) (*UserClaims, error)
	getUser(id uuid.UUID) (*ent.User, error)
}

type service struct {
	client   *ent.Client
	tokenSvc tokenService
}

func newService(client *ent.Client, secretKeyHex string) (Service, error) {
	secretKey, err := paseto.V4SymmetricKeyFromHex(secretKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}

	tokenSvc := newTokenService(secretKey)
	return &service{client: client, tokenSvc: tokenSvc}, nil
}

func (s service) validateCredentials(credentials LoginRequest) (*ent.User, error) {
	// Get user with given email
	u, err := s.client.User.Query().
		Where(user.EmailEQ(credentials.Email)).
		First(context.Background())
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	// Compare the password hashes
	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(credentials.Password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}
	// Return user if valid credentials
	return u, nil
}

// getUser returns a user with the given id
func (s service) getUser(id uuid.UUID) (*ent.User, error) {
	u, err := s.client.User.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s service) generateAccessToken(user *UserClaims) string {
	return s.tokenSvc.generateToken(user, AccessToken)
}

func (s service) generateRefreshToken(user *UserClaims) string {
	return s.tokenSvc.generateToken(user, RefreshToken)
}

func (s service) generateTokens(user *UserClaims) *TokenResponse {
	accessToken := s.generateAccessToken(user)
	refreshToken := s.generateRefreshToken(user)

	return &TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}
}

func (s service) parseToken(token string) (*UserClaims, error) {
	parsedToken, err := s.tokenSvc.parseToken(token)
	if err != nil {
		return nil, err
	}
	return extractUserClaims(parsedToken), nil
}

// extractUserClaims extracts the claims in the token
func extractUserClaims(parsedToken *paseto.Token) *UserClaims {
	// Already validated so no errors expected.
	ID, _ := parsedToken.GetString("user_id")
	email, _ := parsedToken.GetString("email")
	role, _ := parsedToken.GetString("role")
	scope, _ := parsedToken.GetString("scope")

	return &UserClaims{
		Email: email,
		ID:    uuid.MustParse(ID),
		Role:  user.Role(role),
		Scope: Scope(scope),
	}
}

// hashRefreshToken hashes a the refresh token to be stored in the sessions table
func hashRefreshToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
