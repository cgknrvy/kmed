package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"kmed/api/ent"
	"kmed/api/ent/user"
	"kmed/api/internal/auth"
	"kmed/api/internal/httpx"
)

type Service interface {
	getUser(id uuid.UUID) (*ent.User, error)
	getUsers() ([]*ent.User, error)
	getUserByEmail(email string) (*ent.User, error)
	createUser(user CreateRequest) (*ent.User, error)
	updateUser(ctx context.Context, user UpdateRequest) (*ent.User, error)
	updatePassword(ctx context.Context, updateReq PasswordUpdateRequest) (*ent.User, error)
	deleteUser(user DeleteRequest) error
}

type service struct {
	client *ent.Client
}

func newService(client *ent.Client) Service {
	return &service{client: client}
}

// getUser returns a user with the given id
func (s *service) getUser(id uuid.UUID) (*ent.User, error) {
	u, err := s.client.User.Get(context.Background(), id)
	if err != nil {
		return nil, GetError{err: err}
	}

	return u, nil
}

// getUsers returns all users in the database
func (s *service) getUsers() ([]*ent.User, error) {
	users, err := s.client.User.Query().All(context.Background())
	if err != nil {
		return nil, GetError{err}
	}
	return users, nil
}

// getUserByEmail returns user with the given email
func (s *service) getUserByEmail(email string) (*ent.User, error) {
	u, err := s.client.User.Query().Where(user.EmailEQ(email)).First(context.Background())
	// TODO: Check for not found error to return more descriptive error
	if err != nil {
		return nil, GetError{err: err}
	}
	return u, nil
}

type CreateRequest struct {
	Name     string    `json:"name"     validate:"required,min=3,max=20"`
	Email    string    `json:"email"    validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=72"`
	Role     user.Role `json:"role"     validate:"required,oneof=admin doctor lab-tech user"`
}

func (r CreateRequest) Validate() error {
	return httpx.Validator.Struct(r)
}

// createUser creates a new user with the given user data
func (s *service) createUser(user CreateRequest) (*ent.User, error) {
	if err := user.Validate(); err != nil {
		return nil, CreateError{err}
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, CreateError{err}
	}

	createdUser, err := s.client.User.Create().
		SetEmail(user.Email).
		SetName(user.Name).
		SetRole(user.Role).
		SetPassword(hashedPassword).
		Save(context.Background())
		// TODO: Check for constraint error to return more descriptive error
	if err != nil {
		return nil, CreateError{err: err}
	}

	return createdUser, nil
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id" validator:"required,uuid"`
}

func (r DeleteRequest) Validate() error {
	return httpx.Validator.Struct(r)
}

// deleteUser deletes user with the given id as in the [DeleteRequest]
func (s *service) deleteUser(user DeleteRequest) error {
	err := s.client.User.DeleteOneID(user.ID).Exec(context.Background())
	if ent.IsNotFound(err) {
		return err
	} else if err != nil {
		return DeleteError{err: err}
	}

	return nil
}

type UpdateRequest struct {
	NewName  *string `json:"new_name,omitempty"  validate:"omitnil,min=3,max=20"`
	NewEmail *string `json:"new_email,omitempty" validate:"omitnil,email"`
}

func (r UpdateRequest) Validate() error {
	if r.NewName == nil && r.NewEmail == nil {
		return fmt.Errorf("must provide either Name or Email or both when updating")
	}
	return httpx.Validator.Struct(r)
}

// updateUser updates a user with the new data provided in the [UpdateRequest].
// id of the user to update is gotten from the passed context `ctx`.
//
// This function does not update the password. To update password use the [service.updatePassword]
// function instead.
func (s *service) updateUser(ctx context.Context, updateReq UpdateRequest) (*ent.User, error) {
	// Validation is done by the httpx.Parse function in the handler

	userClaims, ok := ctx.Value(auth.UserContextKey).(auth.UserClaims)
	if !ok {
		return nil, fmt.Errorf("cannot access user claims from context")
	}

	update := s.client.User.UpdateOneID(userClaims.ID)
	if updateReq.NewName != nil {
		update.SetName(*updateReq.NewName)
	}
	if updateReq.NewEmail != nil {
		update.SetEmail(*updateReq.NewEmail)
	}

	// TODO: Check for constraint error to return more descriptive error
	if updatedUser, err := update.Save(context.Background()); err != nil {
		return nil, UpdateError{err}
	} else {
		return updatedUser, nil
	}
}

type PasswordUpdateRequest struct {
	OldPassword string `json:"old_password" validate:"min=8,max=72"`
	NewPassword string `json:"new_password" validate:"min=8,max=72"`
}

func (r PasswordUpdateRequest) Validate() error {
	return httpx.Validator.Struct(r)
}

// updatePassword updates the password of the given user.
// user id is obtained from the passed context `ctx`.
//
// TO update other fields of user use the [service.updateUser] function instead.
func (s *service) updatePassword(
	ctx context.Context,
	updateReq PasswordUpdateRequest,
) (*ent.User, error) {
	userClaims, ok := ctx.Value(auth.UserContextKey).(auth.UserClaims)
	if !ok {
		return nil, fmt.Errorf("cannot access user claims from context")
	}

	if err := s.validatePassword(userClaims.ID, updateReq.OldPassword); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	hashedPassword, err := hashPassword(updateReq.NewPassword)
	if err != nil {
		return nil, UpdateError{err}
	}

	return s.client.User.UpdateOneID(userClaims.ID).
		SetPassword(hashedPassword).
		Save(context.Background())
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}

// validatePassword validates that the given password is correct for the user
// with the passed id.
func (s *service) validatePassword(id uuid.UUID, password string) error {
	// Get user with given id
	u, err := s.client.User.Get(context.Background(), id)
	if err != nil {
		return err
	}
	// Compare the password hashes
	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	); err != nil {
		return err
	}
	return nil
}

type UpdateError struct {
	err error
}

func (e UpdateError) Error() string {
	msg := splitError(e.err)
	return fmt.Sprintf("users: %s", msg)
}

type DeleteError struct {
	err error
}

func (e DeleteError) Error() string {
	msg := splitError(e.err)
	return fmt.Sprintf("users: %s", msg)
}

type GetError struct {
	err error
}

func (e GetError) Error() string {
	msg := splitError(e.err)
	return fmt.Sprintf("users: %s", msg)
}

type CreateError struct {
	err error
}

func (e CreateError) Error() string {
	msg := splitError(e.err)
	return fmt.Sprintf("users: failed to create user: %s", msg)
}

func splitError(e error) string {
	a := strings.Split(e.Error(), ":")

	if len(a) == 0 {
		return ""
	} else if len(a) == 1 {
		return strings.TrimSpace(a[0])
	}
	s := strings.Join(a[1:], ":")

	return strings.TrimSpace(s)
}
