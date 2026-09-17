package patient

import (
	"net/http"

	"kmed/api/internal/auth"
	"kmed/api/internal/errors"
	"kmed/api/internal/httpx"

	"github.com/justinas/alice"
)

// Handler is the handler for the Patients.
type Handler struct {
	svc  Service
	auth auth.Middleware
}

func NewHandler(service Service, authMiddleware auth.Middleware) *Handler {
	return &Handler{svc: service, auth: authMiddleware}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	authChain := alice.New(h.auth.Authenticate, h.auth.RequireAdminOrDoctor)

	mux.Handle("GET /{id}", authChain.Then(http.HandlerFunc(h.getPatient)))
	mux.Handle("GET /", authChain.Then(http.HandlerFunc(h.getPatients)))
	mux.Handle("GET /search", authChain.Then(http.HandlerFunc(h.searchPatientsByName)))
	mux.Handle("GET /count", authChain.Then(http.HandlerFunc(h.getPatientsTotalCount)))
	mux.Handle("GET /count/today", authChain.Then(http.HandlerFunc(h.getPatientsCreatedTodayCount)))
	mux.Handle("POST /", authChain.Then(http.HandlerFunc(h.createPatient)))
	mux.Handle("DELETE /", authChain.Then(http.HandlerFunc(h.deletePatient)))
	mux.Handle("PUT /", authChain.Then(http.HandlerFunc(h.updatePatient)))

	return mux
}

func (h *Handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router().ServeHTTP(w, r)
}

// getPatient returns a single patient with the given id
func (h *Handler) getPatient(w http.ResponseWriter, r *http.Request) {
	id, ok := httpx.GetIDFromPath(w, r)
	if !ok {
		return
	}

	patient, err := h.svc.getPatient(id)

	switch {
	case err == nil:
		httpx.JSON(w, http.StatusOK, httpx.Response{Patient: patient})
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

// getPatients returns all the patients in the database.
// TODO: Add a filter to enable only getting specific patients
func (h *Handler) getPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := h.svc.getPatients()
	switch err {
	case nil:
		httpx.JSON(w, http.StatusOK,
			httpx.Response{
				Patients: patients,
			})
		return
	default:
		httpx.JSONError(w, http.StatusInternalServerError,
			httpx.ErrorResponse{
				Message: "patient not found",
			})
	}
}

func (h *Handler) searchPatientsByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		httpx.JSON(w, http.StatusOK, httpx.Response{SearchResults: []*httpx.SearchResult{}})
		return
	}

	results, err := h.svc.searchPatientsByName(name)
	if err != nil {
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
		return
	}
	httpx.JSON(w, http.StatusOK, httpx.Response{SearchResults: results})
}

// getPatientsTotalCount returns the total count of patients
func (h *Handler) getPatientsTotalCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.svc.getPatientsTotalCount()
	if err != nil {
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
		return
	}
	httpx.JSON(w, http.StatusOK, httpx.Response{Count: ptr(count)})
}

// getPatientsCreatedTodayCount returns total count of patients created in the last 24hrs
func (h *Handler) getPatientsCreatedTodayCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.svc.getPatientsCreatedTodayCount()
	if err != nil {
		httpx.JSONError(
			w,
			http.StatusInternalServerError,
			httpx.ErrorResponse{Message: err.Error()},
		)
		return
	}
	httpx.JSON(w, http.StatusOK, httpx.Response{Count: ptr(count)})
}

// createPatient creates a new patient with the given data.
// The patient data in the request must be able to be parsed
// into the type [CreateRequest].
func (h *Handler) createPatient(w http.ResponseWriter, r *http.Request) {
	req, ok := httpx.Parse[CreateRequest](w, r)
	if !ok {
		return
	}

	createdPatient, err := h.svc.createPatient(req)
	switch {
	case err == nil:
		httpx.JSON(w, http.StatusAccepted, httpx.Response{
			Patient: createdPatient,
		})
		return
	case errors.IsBadRequest(err):
		httpx.JSONError(w, http.StatusBadRequest,
			httpx.ErrorResponse{
				Message: err.Error(),
			})
		return
	default:
		httpx.JSONError(w, http.StatusInternalServerError,
			httpx.ErrorResponse{
				Message: "failed to create patient",
				Errors:  map[string]string{"create": err.Error()},
			})
	}
}

// deletePatient deletes a single patient with the given id.
// Id is provided in the body following the format of [DeleteRequest].
func (h *Handler) deletePatient(w http.ResponseWriter, r *http.Request) {
	req, ok := httpx.Parse[DeleteRequest](w, r)
	if !ok {
		return
	}

	err := h.svc.deletePatient(req.ID)
	switch {
	case err == nil:
		httpx.JSON(w, http.StatusAccepted, httpx.Response{
			Message: ptr("deleted patient"),
		})
		return
	case errors.IsNotFound(err):
		httpx.JSONError(w, http.StatusNotFound,
			httpx.ErrorResponse{
				Message: err.Error(),
			})
		return
	default:
		httpx.JSONError(w, http.StatusInternalServerError,
			httpx.ErrorResponse{
				Message: "failed to delete patient",
				Errors:  map[string]string{"delete": err.Error()},
			})
	}
}

func (h *Handler) updatePatient(w http.ResponseWriter, r *http.Request) {
	httpx.JSONError(w, http.StatusNotImplemented, httpx.ErrorResponse{
		Message: "not implemented",
	})
}

func ptr[T any](v T) *T {
	return &v
}
