package auth

import (
	"errors"
	"fmt"
	"time"

	"kmed/api/ent/user"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type Scope string

const (
	Full           Scope = "full"
	PasswordChange Scope = "password_change"
)

type UserClaims struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  user.Role `json:"role"`
	Scope Scope     `json:"scope"`
}

type tokenService struct {
	key           paseto.V4SymmetricKey
	parser        paseto.Parser
	token         paseto.Token
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

const (
	AccessTTL         = 15 * time.Minute
	PasswordChangeTTL = 15 * time.Minute
	RefreshTTL        = 24 * time.Hour
)

const (
	Issuer   = "kmed-auth"
	Audience = "kmed-api"
)

func newTokenService(key paseto.V4SymmetricKey) tokenService {
	// token with the common claims set
	token := paseto.NewToken()
	token.SetIssuer(Issuer)
	token.SetAudience(Audience)

	// parser with rules added
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ForAudience(Audience))
	parser.AddRule(paseto.IssuedBy(Issuer))

	return tokenService{
		key:           key,
		parser:        parser,
		token:         token,
		accessExpiry:  AccessTTL,
		refreshExpiry: RefreshTTL,
	}
}

type tokenType string

const (
	AccessToken  tokenType = "access_token"
	RefreshToken tokenType = "refresh_token"
)

func (t tokenService) generateToken(userClaims *UserClaims, tt tokenType) string {
	token := t.token // Copy the default token
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetSubject(userClaims.ID.String())

	switch tt {
	case AccessToken:
		token.SetExpiration(time.Now().Add(t.accessExpiry))
		// Custom claims
		token.SetString("user_id", userClaims.ID.String())
		token.SetString("email", userClaims.Email)
		token.SetString("role", string(userClaims.Role))
	case RefreshToken:
		token.SetExpiration(time.Now().Add(t.refreshExpiry))
	}

	return token.V4Encrypt(t.key, nil)
}

func (t tokenService) parseToken(token string) (*paseto.Token, error) {
	t.parser.AddRule(paseto.ValidAt(time.Now()))

	parsedToken, err := t.parser.ParseV4Local(t.key, token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if err := t.validateCustomClaims(parsedToken); err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	return parsedToken, nil
}

func (t tokenService) validateCustomClaims(token *paseto.Token) error {
	userID, err := token.GetString("user_id")
	if err != nil || userID == "" {
		return errors.New("missing user id in claims")
	}
	// Ensure userID is a valid UUID
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("invalid user id in claims")
	}

	email, err := token.GetString("email")
	if err != nil || email == "" {
		return errors.New("missing email in claims")
	}

	role, err := token.GetString("role")
	if err != nil {
		return errors.New("missing role in claims")
	}
	if err := user.RoleValidator(user.Role(role)); err != nil {
		return errors.New("invalid role in claims")
	}

	return nil
}
