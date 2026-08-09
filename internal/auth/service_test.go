package auth

import (
	"context"
	"testing"

	"aidanwoods.dev/go-paseto"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"kmed/api/ent"
	entUser "kmed/api/ent/user"
)

const testSecretKeyHex = "ab0f159045e95df5f9fc84194a6662cb0fc589c3c9fb323d876275a3d25a7525"

func newTestClient(t *testing.T) *ent.Client {
	t.Helper()

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}

func createMockUser(t *testing.T) (*ent.Client, *ent.User) {
	t.Helper()
	client := newTestClient(t)

	password, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed hashing password: %v", err)
	}
	// Create dummy user
	mockUser, err := client.User.Create().
		SetName("Mocky").
		SetEmail("test@kmed.com").
		SetRole(entUser.RoleUser).
		SetPassword(string(password)).
		Save(context.Background())
	if err != nil {
		t.Fatalf("failed to create mock user: %v", err)
	}

	return client, mockUser
}

func TestValidateCredentials(t *testing.T) {
	client, mockUser := createMockUser(t)
	svc, err := newService(client, testSecretKeyHex)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	t.Run("valid credentials", func(t *testing.T) {
		// LoginRequest
		credentials := LoginRequest{
			Email:    mockUser.Email,
			Password: "12345678",
		}
		// validate credentials
		user, err := svc.validateCredentials(credentials)
		if assert.NoError(t, err, "expecting no error for valid credentials") {
			assert.Equal(t, mockUser.Email, user.Email, "expecting user email to match")
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		credentials := LoginRequest{
			Email:    "invalid-email@kmed.com",
			Password: "12345678",
		}
		user, err := svc.validateCredentials(credentials)
		assert.Nil(t, user, "expecting user to be nil")
		assert.Error(t, err, "expecting error for invalid email")
		assert.ErrorAs(t, err, &ErrInvalidCredentials)
	})

	t.Run("invalid password", func(t *testing.T) {
		credentials := LoginRequest{
			Email:    mockUser.Email,
			Password: "invalid-password",
		}
		user, err := svc.validateCredentials(credentials)
		assert.Nil(t, user, "expecting user to be nil")
		assert.Error(t, err, "expecting error for invalid email")
		assert.ErrorAs(t, err, &ErrInvalidCredentials)
	})
}

func TestGenerateToken(t *testing.T) {
	client, mockUser := createMockUser(t)
	svc, err := newService(client, testSecretKeyHex)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	t.Run("create token", func(t *testing.T) {
		token := svc.generateTokens(
			&UserClaims{ID: mockUser.ID, Email: mockUser.Email, Role: mockUser.Role},
		)
		assert.NotEmpty(t, token, "token should not be empty")

		secretKey, err := paseto.V4SymmetricKeyFromHex(testSecretKeyHex)
		if err != nil {
			t.Fatalf("failed to generate symmetric key: %v", err)
		}
		parsedToken, err := tokenService{}.parser.ParseV4Local(secretKey, token.AccessToken, nil)
		assert.NoError(t, err, "expecting no error for valid token")

		iss, err := parsedToken.GetIssuer()
		assert.NoError(t, err)
		assert.Equal(t, iss, Issuer)

		aud, err := parsedToken.GetAudience()
		assert.NoError(t, err)
		assert.Equal(t, aud, Audience)

		userID, err := parsedToken.GetString("user_id")
		assert.NoError(t, err)
		assert.Equal(t, userID, mockUser.ID.String())

		email, err := parsedToken.GetString("email")
		assert.NoError(t, err)
		assert.Equal(t, email, mockUser.Email)
	})
}

func TestParseToken(t *testing.T) {
	t.Run("parses valid token without error", func(t *testing.T) {})
	t.Run("parsing invalid token returns error", func(t *testing.T) {})
}
