package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"kmed/api/ent"
	"kmed/api/ent/user"
	"kmed/api/internal/httpx"
)

type Service interface {
	getUser(id uuid.UUID) (*ent.User, error)
	getUserByEmail(email string) (*ent.User, error)
	createUser(user CreateRequest) (*ent.User, error)
	updateUser(id uuid.UUID, user *UpdateRequest) (*ent.User, error)
	deleteUser(user DeleteRequest) error
}

type service struct {
	client *ent.Client
}

func newService(client *ent.Client) Service {
	return &service{client: client}
}

func (s *service) getUser(id uuid.UUID) (*ent.User, error) {
	u, err := s.client.User.Get(context.Background(), id)
	if err != nil {
		return nil, GetError{err: err}
	}

	return u, nil
}

func (s *service) getUserByEmail(email string) (*ent.User, error) {
	u, err := s.client.User.Query().Where(user.EmailEQ(email)).First(context.Background())
	//TODO: Check for not found error to return more descriptive error
	if err != nil {
		return nil, GetError{err: err}
	}
	return u, nil
}

type CreateRequest struct {
	Name     string    `json:"name" validate:"required,min=3,max=20"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=72"`
	Role     user.Role `json:"role" validate:"required,oneof=admin doctor lab-tech user"`
}

func (r CreateRequest) Validate() error {
	return httpx.Validator.Struct(r)
}

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

	//TODO: Check for constraint error to return more descriptive error
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
	Name  *string
	Email *string
}

func (r UpdateRequest) Validate() error {
	if r.Name == nil && r.Email == nil {
		return fmt.Errorf("must provide either Name or Email or both when updating")
	}
	return nil
}

func (s *service) updateUser(id uuid.UUID, updateReq *UpdateRequest) (*ent.User, error) {
	if updateReq == nil {
		return nil, nil
	}
	if err := updateReq.Validate(); err != nil {
		return nil, nil // Don't update if both fields are nil
	}

	update := s.client.User.UpdateOneID(id)
	if updateReq.Name != nil {
		update.SetName(*updateReq.Name)
	}
	if updateReq.Email != nil {
		update.SetEmail(*updateReq.Email)
	}

	//TODO: Check for constraint error to return more descriptive error
	if updatedUser, err := update.Save(context.Background()); err != nil {
		return nil, UpdateError{err}
	} else {
		return updatedUser, nil
	}
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
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
