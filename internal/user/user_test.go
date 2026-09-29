package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"kmed/api/ent"
	entUser "kmed/api/ent/user"
	"kmed/api/internal/httpx"
)

const invalidUserID = "0198f7c3-12ab-8xyz-9abc-1234567890ab"

func TestGETUser(t *testing.T) {
	store := StubUserStore{map[string]ent.User{
		"0198f7c3-12ab-7def-8abc-1234567890ab": {Name: "One"},
		"0198f7c3-45cd-7123-9ef0-fedcba987654": {Name: "Two"},
	}, nil, nil}
	userService := &Handler{&store, authMiddleware}

	t.Run("returns user with id 0198f7c3-12ab-7def-8abc-1234567890ab", func(t *testing.T) {
		request := newGetUserRequest("0198f7c3-12ab-7def-8abc-1234567890ab")
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)

		data, _ := json.Marshal(httpx.Response{User: &ent.User{Name: "One"}})
		assertResponseBody(t, strings.TrimSpace(response.Body.String()), string(data))
		assertStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("returns user with id 0198f7c3-45cd-7123-9ef0-fedcba987654", func(t *testing.T) {
		request := newGetUserRequest("0198f7c3-45cd-7123-9ef0-fedcba987654")
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)

		data, _ := json.Marshal(httpx.Response{User: &ent.User{Name: "Two"}})
		assertResponseBody(t, strings.TrimSpace(response.Body.String()), string(data))
		assertStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 for missing user", func(t *testing.T) {
		request := newGetUserRequest("0198f7c3-89ef-7a56-b123-0fedcba98765")
		response := httptest.NewRecorder()
		userService.serveHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
	})
}

func TestPOSTUser(t *testing.T) {
	store := StubUserStore{map[string]ent.User{}, nil, nil}
	userService := &Handler{&store, authMiddleware}

	t.Run("returns accepted on POST", func(t *testing.T) {
		user := CreateRequest{
			Name:     "NewUser",
			Email:    "user@test.com",
			Password: "password",
			Role:     "user",
		}
		request := newPostUserRequest(user)
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.createCalls) != 1 {
			t.Errorf("userService should have created 1 user, got %d", len(store.createCalls))
		}

		if store.createCalls[0].Email != user.Email {
			t.Errorf(
				"userService should have created user with email %s, got %s",
				user.Email,
				store.createCalls[0].Email,
			)
		}
	})

	t.Run("returns 400 if email is invalid", func(t *testing.T) {
		user := CreateRequest{Name: "NewUser1", Email: "invalidEmail", Password: "password"}
		request := newPostUserRequest(user)
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
		if len(store.createCalls) != 1 {
			t.Errorf("CreateUser method on the store should not be hit if data is invalid")
		}
	})

	t.Run("returns 400 for invalid body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"id"}`))
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})
}

func TestDELETEUser(t *testing.T) {
	userID := "0198f7c3-12ab-7def-8abc-1234567890ab"
	store := StubUserStore{map[string]ent.User{
		userID: {Name: "One"},
	}, nil, nil}
	userService := &Handler{&store, authMiddleware}

	t.Run("returns accepted on DELETE", func(t *testing.T) {
		request := newDeleteRequest(userID)
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.deleteCalls) != 1 {
			t.Errorf("userService should have deleted 1 user, got %d", len(store.deleteCalls))
		}

		if store.deleteCalls[0] != userID {
			t.Errorf(
				"userService should have deleted user with id %s, got %s",
				userID,
				store.deleteCalls[0],
			)
		}
	})

	t.Run("returns 400 if ID is invalid", func(t *testing.T) {
		request := newDeleteRequest(invalidUserID)
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns 400 for invalid body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/", strings.NewReader(`{"id"}`))
		response := httptest.NewRecorder()

		userService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})
}

func newGetUserRequest(id string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", id), nil)
	return request
}

func newPostUserRequest(payload CreateRequest) *http.Request {
	data, _ := json.Marshal(payload)
	request, _ := http.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(string(data)),
	)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func newDeleteRequest(id string) *http.Request {
	request, _ := http.NewRequest(
		http.MethodDelete,
		"/",
		strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, id)),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", id))
	return request
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status code is wrong, got %d want %d", got, want)
	}
}

type StubUserStore struct {
	users       map[string]ent.User
	createCalls []CreateRequest
	deleteCalls []string
}

func (s *StubUserStore) getUser(id uuid.UUID) (*ent.User, error) {
	if user, ok := s.users[id.String()]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *StubUserStore) getUsers() ([]*ent.User, error) {
	return []*ent.User{}, nil
}

func (s *StubUserStore) getUserByEmail(email string) (*ent.User, error) {
	return &ent.User{}, nil
}

func (s *StubUserStore) createUser(user CreateRequest) (*ent.User, error) {
	s.createCalls = append(s.createCalls, user)
	return &ent.User{}, nil
}

func (s *StubUserStore) updateUser(ctx context.Context, user UpdateRequest) (*ent.User, error) {
	return &ent.User{}, nil
}

func (s *StubUserStore) updatePassword(
	ctx context.Context,
	updateReq PasswordUpdateRequest,
) (*ent.User, error) {
	return &ent.User{}, nil
}

func (s *StubUserStore) deleteUser(user DeleteRequest) error {
	s.deleteCalls = append(s.deleteCalls, user.ID.String())
	return nil
}

type AuthMiddleware struct{}

func (a AuthMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (a AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (a AuthMiddleware) RequireLabTech(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (a AuthMiddleware) RequireDoctor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (a AuthMiddleware) RequireAdminOrDoctor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (a AuthMiddleware) RequireRole(next http.Handler, role []entUser.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (a AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

var authMiddleware = AuthMiddleware{}
