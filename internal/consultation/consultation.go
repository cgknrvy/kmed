package consultation

import (
	"net/http"

	"kmed/api/ent"
	"kmed/api/internal/auth"
	"kmed/api/internal/errors"
	"kmed/api/internal/httpx"
	"kmed/api/internal/icd"

	"github.com/justinas/alice"
)

type Handler struct {
	svc  Service
	auth auth.Middleware
}

func NewHandler(client *ent.Client, authMiddleware auth.Middleware) *Handler {
	svc := newService(client)
	return &Handler{svc: svc, auth: authMiddleware}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	authChain := alice.New(h.auth.Authenticate, h.auth.RequireDoctor)

	mux.Handle("POST /", authChain.Then(http.HandlerFunc(h.createConsultation)))
	mux.Handle("GET /{id}", authChain.Then(http.HandlerFunc(h.getConsultation)))
	mux.Handle("GET /patient/{id}", authChain.Then(http.HandlerFunc(h.getPatientConsultations)))
	mux.Handle(
		"GET /doctor/{id}",
		authChain.Then(http.HandlerFunc(h.getTodaysConsultationsForDoctor)),
	)
	mux.Handle("GET /icd11/search", authChain.Then(http.HandlerFunc(h.searchICD11Codes)))

	return mux
}

func (h *Handler) createConsultation(w http.ResponseWriter, r *http.Request) {
	req, ok := httpx.Parse[CreateRequest](w, r)
	if !ok {
		return
	}

	consultation, err := h.svc.createConsultation(req)
	switch {
	case err == nil:
		httpx.JSON(w, http.StatusAccepted, httpx.Response{Consultation: consultation})
		return
	case errors.IsBadRequest(err) || errors.IsNotFound(err): // NotFoundErr is returned if the doctor or patient IDs are invalid
		httpx.JSONError(w, http.StatusBadRequest, httpx.ErrorResponse{Message: err.Error()})
		return
	default:
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
		return
	}
}

func (h *Handler) getConsultation(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.GetIDFromPath(w, r)
	if !ok {
		return
	}

	consultation, err := h.svc.getConsultation(id)
	switch {
	case err == nil:
		httpx.JSON(w, http.StatusOK, httpx.Response{Consultation: consultation})
		return
	case errors.IsNotFound(err):
		httpx.JSONError(w, http.StatusNotFound, httpx.ErrorResponse{Message: err.Error()})
		return
	default:
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
		return
	}
}

func (h *Handler) getPatientConsultations(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.GetIDFromPath(w, r)
	if !ok {
		return
	}

	consultations, err := h.svc.getPatientConsultations(id)
	switch {
	case err == nil:
		httpx.JSON(w, http.StatusOK, httpx.Response{Consultations: consultations})
		return
	case errors.IsNotFound(err):
		httpx.JSONError(w, http.StatusNotFound, httpx.ErrorResponse{Message: err.Error()})
		return
	default:
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
	}
}

func (h *Handler) getTodaysConsultationsForDoctor(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.GetIDFromPath(w, r)
	if !ok {
		return
	}

	consultations, err := h.svc.getTodaysConsultationsForDoctor(id)
	switch {
	case err == nil:
		httpx.JSON(w, http.StatusOK, httpx.Response{Consultations: consultations})
		return
	case errors.IsNotFound(err):
		httpx.JSONError(w, http.StatusNotFound, httpx.ErrorResponse{Message: err.Error()})
		return
	default:
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
	}
}

func (h *Handler) searchICD11Codes(w http.ResponseWriter, r *http.Request) {
	diseaseName := r.URL.Query().Get("q")
	codes, err := icd.SearchICD11(diseaseName)
	if err != nil {
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
		return
	}

	var vcodes []httpx.DestinationEntity
	for _, code := range codes {
		vcodes = append(
			vcodes,
			httpx.DestinationEntity{
				Code:  code.Code,
				Title: code.Title,
				Score: code.Score,
			},
		)
	}

	httpx.JSON(
		w,
		http.StatusOK,
		httpx.Response{ICD11SearchResults: vcodes},
	)
}

func (h *Handler) deleteConsultation(w http.ResponseWriter, r *http.Request) {
}
