package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"kmed/api/ent"
	"kmed/api/ent/user"
	"kmed/api/internal/auth"
)

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

func TestSqliteStore_CreateUser(t *testing.T) {
	t.Run("creating a new user", func(t *testing.T) {
		client := newTestClient(t)
		svc := newService(client)
		newUser := CreateRequest{
			Email:    "user@kmed.com",
			Password: "12345678",
			Name:     "KMed",
			Role:     "user",
		}

		createdUser, err := svc.createUser(newUser)
		assert.Nil(t, err, "Expected no error when creating new user")

		foundUser, err := client.User.Query().
			Where(user.EmailEQ(newUser.Email)).
			First(context.Background())

		if assert.NoError(t, err) {
			assert.Equal(t, User{}.FromEntUser(createdUser), User{}.FromEntUser(foundUser))
			// ensure that the password is hashed
			assert.Nil(
				t,
				bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(newUser.Password)),
			)
		}
	})

	t.Run("creating an existing user returns error", func(t *testing.T) {
		client := newTestClient(t)
		svc := newService(client)
		userDetails := CreateRequest{
			Email:    "user@kmed.com",
			Password: "12345678",
			Name:     "KMed",
			Role:     "user",
		}
		// Add mock user to db
		_, err := client.User.Create().
			SetName(userDetails.Name).
			SetEmail(userDetails.Email).
			SetRole(userDetails.Role).
			SetPassword(userDetails.Password).
			Save(context.Background())
		assert.Nil(t, err, "Expected to create mock user")
		_, err = client.User.Query().
			Where(user.EmailEQ(userDetails.Email)).
			First(context.Background())
		assert.Nil(t, err, "Expected to find the mock user")

		createdUser, err := svc.createUser(userDetails) // Create same user again
		assert.Nil(t, createdUser, "Expected duplicate user to not be created")
		assert.Error(t, err, "Expected error to be returned for duplicate user")
	})

	t.Run("creating user with invalid email returns error", func(t *testing.T) {
		client := newTestClient(t)
		svc := newService(client)

		_, err := svc.createUser(
			CreateRequest{Email: "invalid email", Name: "One", Password: "12345678"},
		)
		assert.Error(t, err, "Expected error to be returned for invalid email")
		assert.ErrorAs(t, err, &CreateError{}, "Expected error to be a ValidationError")
		assert.ErrorContains(t, err, "email", "Expected the error message to contain word email")
	})

	t.Run("creating user with invalid Password returns error", func(t *testing.T) {
		client := newTestClient(t)
		svc := newService(client)

		_, err := svc.createUser(
			CreateRequest{Email: "test@kmed.com", Name: "One", Password: "123", Role: "user"},
		)
		assert.Error(t, err, "Expected error to be returned for invalid password")
		assert.ErrorAs(t, err, &CreateError{}, "Expected error to be a ValidationError")
	})

	t.Run("creating user with invalid Role returns error", func(t *testing.T) {
		client := newTestClient(t)
		svc := newService(client)

		_, err := svc.createUser(
			CreateRequest{
				Email:    "test@kmed.com",
				Name:     "One",
				Password: "123",
				Role:     "non-existent",
			},
		)
		assert.Error(t, err, "Expected error to be returned for invalid role")
		assert.ErrorAs(t, err, &CreateError{}, "Expected error to be a ValidationError")
	})
}

func createMockUser(t *testing.T) (*ent.Client, *ent.User) {
	t.Helper()
	client := newTestClient(t)

	// Create dummy user
	mockUser, err := client.User.Create().
		SetName("Mocky").
		SetEmail("test@kmed.com").
		SetRole("user").
		SetPassword("12345678").
		Save(context.Background())
	if err != nil {
		t.Fatalf("failed to create mock user: %v", err)
	}

	return client, mockUser
}

func TestSqliteStore_GetUser(t *testing.T) {
	client, mockUser := createMockUser(t)
	svc := newService(client)

	t.Run("getting a user", func(t *testing.T) {
		gottenUser, err := svc.getUser(mockUser.ID)
		if assert.NoError(t, err, "Expected no error when getting an existing user") {
			assert.Equal(
				t,
				User{}.FromEntUser(gottenUser),
				User{}.FromEntUser(mockUser),
				"Expected mock user to be returned",
			)
		}
	})

	t.Run("getting a nonexistent user", func(t *testing.T) {
		gottenUser, err := svc.getUser(uuid.Must(uuid.NewV7()))
		assert.Nil(t, gottenUser, "Expected user to be nil when getting nonexistent user")
		assert.ErrorAs(t, err, &GetError{}, "Expected error to be a GetError")
	})
}

func TestSqliteStore_DeleteUser(t *testing.T) {
	client, mockUser := createMockUser(t)
	svc := newService(client)

	t.Run("deleting a user", func(t *testing.T) {
		gottenUser, err := svc.getUser(mockUser.ID)
		if assert.NoError(t, err, "Expected no error when getting an existing user") {
			assert.Equal(
				t,
				User{}.FromEntUser(gottenUser),
				User{}.FromEntUser(mockUser),
				"Expected mock user to be returned",
			)
		}

		err = svc.deleteUser(DeleteRequest{ID: mockUser.ID})
		assert.NoError(t, err, "Expected no error when deleting a user")

		gottenUser, err = svc.getUser(mockUser.ID)
		if assert.Nil(t, gottenUser, "Expected user to be nil after being deleted") {
			assert.Error(t, err, "Expected error to be returned when getting deleted user")
		}
	})

	t.Run("deleting a nonexistent user", func(t *testing.T) {
		userID := uuid.Must(uuid.NewV7())
		gottenUser, _ := svc.getUser(userID)
		assert.Nil(t, gottenUser, "Expected user to be nil when getting nonexistent user")

		err := svc.deleteUser(DeleteRequest{ID: userID})
		assert.Error(t, err, "Expected error to be returned when deleting nonexistent user")
		assert.True(t, ent.IsNotFound(err), "Expected error to be an ent.NotFoundError")
	})
}

func TestSqliteStore_UpdateUser(t *testing.T) {
	client, mockUser := createMockUser(t)
	svc := newService(client)

	t.Run("updating a user", func(t *testing.T) {
		gottenUser, _ := svc.getUser(mockUser.ID)
		assert.NotNil(t, gottenUser, "Expected user to be returned")

		update := UpdateRequest{NewName: ptr("UpdatedName"), NewEmail: ptr("UpdatedEmail@kmed.com")}
		updatedUser, err := svc.updateUser(
			context.WithValue(
				context.Background(),
				auth.UserContextKey,
				auth.UserClaims{Email: gottenUser.Email, ID: gottenUser.ID},
			),
			update,
		)

		if assert.NoError(t, err, "Expected no error when updating user") {
			assert.Equal(t, *update.NewName, updatedUser.Name, "Expected user name to be updated")
			assert.Equal(
				t,
				*update.NewEmail,
				updatedUser.Email,
				"Expected user email to be updated",
			)
			assert.Less(
				t,
				gottenUser.UpdatedAt,
				updatedUser.UpdatedAt,
				"Expected user updated time to be updated",
			)
		}
	})

	t.Run("updating a nonexistent user", func(t *testing.T) {
		userID := uuid.Must(uuid.NewV7())
		_, err := svc.getUser(userID)
		assert.ErrorAs(
			t,
			err,
			&GetError{},
			"Expected error to be returned when getting nonexistent user",
		)

		update := UpdateRequest{NewName: ptr("NewName")}
		updatedUser, err := svc.updateUser(context.WithValue(
			context.Background(),
			auth.UserContextKey,
			auth.UserClaims{ID: userID},
		), update)

		if assert.Nil(t, updatedUser, "Expected user to be nil when updating nonexistent user") {
			assert.Error(t, err, "Expected error to be returned when updating nonexistent user")
			assert.ErrorAs(t, err, &UpdateError{}, "Expected error to be an UpdateError")
		}
	})
}
