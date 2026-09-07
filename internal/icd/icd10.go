package icd

import (
	"net/http"
	"strings"

	"kmed/api/internal/auth"
	"kmed/api/internal/httpx"

	"github.com/justinas/alice"
)

type Handler struct {
	svc  Service
	auth auth.Middleware
}

func NewHandler(service Service, authMiddleware auth.Middleware) *Handler {
	return &Handler{svc: service, auth: authMiddleware}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	authChain := alice.New(h.auth.Authenticate)

	mux.Handle("GET /search", authChain.Then(http.HandlerFunc(h.searchICD10Codes)))

	return mux
}

func (h *Handler) searchICD10Codes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	query = strings.TrimSpace(query)
	codes, err := h.svc.searchICD10(query)
	if err != nil {
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
		return
	}

	httpx.JSON(w, http.StatusOK, httpx.Response{ICD10SearchResults: codes})
}
