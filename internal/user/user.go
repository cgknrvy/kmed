package user

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/justinas/alice"

	"kmed/api/ent"
	"kmed/api/internal/auth"
	"kmed/api/internal/httpx"
)

type Handler struct {
	svc  Service
	auth auth.Middleware
}

func NewHandler(client *ent.Client, auth auth.Middleware) *Handler {
	svc := newService(client)
	return &Handler{svc, auth}
}

func (h *Handler) SetService(client *ent.Client) *Handler {
	if h.svc != nil && h.svc.(*service).client != nil {
		h.svc.(*service).client.Close()
	}
	h.svc = newService(client)
	return h
}

func (h *Handler) Router() *http.ServeMux {
	userRouter := http.NewServeMux()

	authChain := alice.New(h.auth.Authenticate)
	adminChain := authChain.Append(h.auth.RequireAdmin, h.auth.RequireFullScope)

	userRouter.Handle("GET /me", authChain.Then(http.HandlerFunc(h.getCurrentUser)))
	userRouter.Handle("GET /{id}", adminChain.Then(http.HandlerFunc(h.getUser)))
	userRouter.Handle("GET /all", adminChain.Then(http.HandlerFunc(h.getUsers)))
	userRouter.Handle("POST /", adminChain.Then(http.HandlerFunc(h.createUser)))
	userRouter.Handle("DELETE /", adminChain.Then(http.HandlerFunc(h.deleteUser)))
	userRouter.Handle("PUT /", authChain.Then(http.HandlerFunc(h.updateUser)))
	userRouter.Handle("PUT /password", authChain.Then(http.HandlerFunc(h.updatePassword)))

	return userRouter
}

// serveHTTP same function call to http.ServeHTTP used for testing the router
func (h *Handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router().ServeHTTP(w, r)
}

func (h *Handler) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	userClaims, _ := r.Context().Value(auth.UserContextKey).(auth.UserClaims)

	user, err := h.svc.getUser(userClaims.ID)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound,
			httpx.ErrorResponse{
				Message: "user not found",
				Errors:  nil,
			})
		return
	}

	httpx.JSON(w, http.StatusOK, httpx.Response{
		User: user,
	})
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.GetIDFromPath(w, r)
	if !ok {
		return
	}

	user, err := h.svc.getUser(id)
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound,
			httpx.ErrorResponse{
				Message: "user not found",
				Errors:  nil,
			})
		return
	}

	httpx.JSON(w, http.StatusOK,
		httpx.Response{
			User: user,
		})
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.getUsers()
	if err != nil {
		httpx.JSONError(w, http.StatusNotFound,
			httpx.ErrorResponse{
				Message: "user not found",
				Errors:  nil,
			})
		return
	}
	httpx.JSON(w, http.StatusOK,
		httpx.Response{Users: users})
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	req, ok := httpx.Parse[CreateRequest](w, r)
	if !ok {
		return
	}

	createdUser, err := h.svc.createUser(req)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError,
			httpx.ErrorResponse{
				Message: "failed to create user",
				Errors:  map[string]string{"create": err.Error()},
			})
		return
	}

	httpx.JSON(w, http.StatusAccepted,
		httpx.Response{
			Message: ptr("created new user"),
			User:    createdUser,
		})
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	req, ok := httpx.Parse[DeleteRequest](w, r)
	if !ok {
		return
	}

	err := h.svc.deleteUser(req)
	if err != nil {
		if ent.IsNotFound(err) {
			httpx.JSONError(w, http.StatusNotFound, httpx.ErrorResponse{
				Message: "user not found",
			})
		} else {
			httpx.JSONError(w, http.StatusInternalServerError, httpx.ErrorResponse{
				Message: "failed to delete user",
				Errors:  map[string]string{"delete": err.Error()},
			})
		}
		return
	}

	httpx.JSON(w, http.StatusAccepted,
		httpx.Response{
			Message: ptr("deleted user"),
		})
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	req, ok := httpx.Parse[UpdateRequest](w, r)
	if !ok {
		return
	}

	updatedUser, err := h.svc.updateUser(r.Context(), req)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, httpx.ErrorResponse{
			Message: "failed to update user",
		})
		return
	}

	httpx.JSON(w, http.StatusAccepted, httpx.Response{User: updatedUser})
}

func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	req, ok := httpx.Parse[PasswordUpdateRequest](w, r)
	if !ok {
		return
	}

	updatedUser, err := h.svc.updatePassword(r.Context(), req)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, httpx.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// TODO: remove all sessions and have the user login again with the new password

	httpx.JSON(w, http.StatusAccepted, httpx.Response{User: updatedUser})
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	Password  string    `json:"-"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u User) FromEntUser(entUser *ent.User) User {
	return User{
		ID:        entUser.ID,
		Name:      entUser.Name,
		Password:  entUser.Password,
		Email:     entUser.Email,
		CreatedAt: entUser.CreatedAt,
		UpdatedAt: entUser.UpdatedAt,
	}
}

func ptr[T any](v T) *T {
	return &v
}
