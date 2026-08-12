package patient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kmed/api/ent"
	entUser "kmed/api/ent/user"
	"kmed/api/internal/date"
	"kmed/api/internal/errors"
	"kmed/api/internal/httpx"

	"github.com/google/uuid"
)

func TestGETPatient(t *testing.T) {
	svc := &StubService{
		map[string]ent.Patient{
			"0198f7c3-12ab-7def-8abc-1234567890ab": {Name: "One"},
			"0198f7c3-45cd-7123-9ef0-fedcba987654": {Name: "Two"},
		}, nil, nil,
	}
	patientService := &Handler{svc, authMiddleware}

	t.Run("GET patient", func(t *testing.T) {
		request := newGetPatientRequest("0198f7c3-12ab-7def-8abc-1234567890ab")
		response := httptest.NewRecorder()

		patientService.serveHTTP(response, request)

		data, _ := json.Marshal(httpx.Response{Patient: &ent.Patient{Name: "One"}})
		assertResponseBody(t, strings.TrimSpace(response.Body.String()), string(data))
		assertStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("GET non-existent patient", func(t *testing.T) {
		request := newGetPatientRequest("0198f7c3-12ab-7def-8abc-1234567890ac")
		response := httptest.NewRecorder()
		patientService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusNotFound)
	})
}

func TestPOSTPatient(t *testing.T) {
	svc := &StubService{map[string]ent.Patient{}, []CreateRequest{}, nil}
	patientService := &Handler{svc, authMiddleware}
	dob := date.Today()

	t.Run("returns accepted on POST", func(t *testing.T) {
		patient := CreateRequest{
			Name:          "Patient",
			DateOfBirth:   &dob,
			Gender:        "male",
			MaritalStatus: "married",
			PhoneNumber:   "+254712345678",
			Email:         "patient@kmed.com",
		}
		request := newPostPatientRequest(patient)
		response := httptest.NewRecorder()
		patientService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(svc.createCalls) != 1 {
			t.Errorf("patientService should have created 1 patient, got %d", len(svc.createCalls))
		}

		if svc.createCalls[0].Name != patient.Name {
			t.Errorf(
				"patientService should have created patient with name %s, got %s",
				patient.Name,
				svc.createCalls[0].Name,
			)
		}
	})

	t.Run("returns 400 for missing name", func(t *testing.T) {
		patient := CreateRequest{DateOfBirth: &dob} // Patient missing Name that is required
		assertInvalidBodyPostPatientRequest(t, patientService, svc, patient)
	})

	t.Run("returns 400 if invalid email is provided", func(t *testing.T) {
		patient := CreateRequest{Name: "Patient", DateOfBirth: &dob, Email: "invalidEmail"}
		assertInvalidBodyPostPatientRequest(t, patientService, svc, patient)
	})

	t.Run("returns 400 if invalid phone number is provided", func(t *testing.T) {
		patient := CreateRequest{Name: "Patient", PhoneNumber: "+254712345678901"}
		assertInvalidBodyPostPatientRequest(t, patientService, svc, patient)
	})

	t.Run("returns 400 if invalid gender is provided", func(t *testing.T) {
		patient := CreateRequest{Name: "Patient", Gender: "malee"}
		assertInvalidBodyPostPatientRequest(t, patientService, svc, patient)
	})

	t.Run("returns 400 if invalid marital status is provided", func(t *testing.T) {
		patient := CreateRequest{Name: "Patient", MaritalStatus: "mal"}
		assertInvalidBodyPostPatientRequest(t, patientService, svc, patient)
	})
}

func TestDELETEPatient(t *testing.T) {
	patientID := "0198f7c3-12ab-7def-8abc-1234567890ab"
	svc := &StubService{
		map[string]ent.Patient{
			patientID: {Name: "Patient One"},
		}, nil, nil,
	}
	patientService := &Handler{svc, authMiddleware}

	t.Run("returns accepted on DELETE", func(t *testing.T) {
		request := newDeleteRequest(patientID)
		response := httptest.NewRecorder()
		patientService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(svc.deleteCalls) != 1 {
			t.Errorf("patientService should have deleted 1 patient, got %d", len(svc.deleteCalls))
		}

		if svc.deleteCalls[0] != patientID {
			t.Errorf(
				"patientService should have deleted patient with id %s, got %s",
				patientID,
				svc.deleteCalls[0],
			)
		}
	})

	t.Run("returns accepted if given ID does not match a patient", func(t *testing.T) {
		request := newDeleteRequest("0198f7c3-89ef-7a56-b123-0fedcba98765")
		response := httptest.NewRecorder()
		patientService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(svc.deleteCalls) != 2 {
			t.Errorf("patientService should have deleted 1 patient, got %d", len(svc.deleteCalls))
		}
	})

	t.Run("returns 400 if given ID is not a valid UUID", func(t *testing.T) {
		request := newDeleteRequest("0198f7c3-12ab-8xyz-9abc-1234567890ab")
		response := httptest.NewRecorder()
		patientService.serveHTTP(response, request)
		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})
}

func newGetPatientRequest(id string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", id), nil)
	return request
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func newDeleteRequest(patientID string) *http.Request {
	request, _ := http.NewRequest(
		http.MethodDelete,
		"/",
		strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, patientID)),
	)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response status code is wrong, got %d want %d", got, want)
	}
}

func newPostPatientRequest(payload CreateRequest) *http.Request {
	data, _ := json.Marshal(payload)
	request, _ := http.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(string(data)),
	)
	request.Header.Set("Content-Type", "application/json")
	return request
}

func assertInvalidBodyPostPatientRequest(
	t *testing.T,
	patientService *Handler,
	service *StubService,
	patient CreateRequest,
) {
	t.Helper()
	request := newPostPatientRequest(patient)
	response := httptest.NewRecorder()
	patientService.serveHTTP(response, request)
	assertStatusCode(t, response.Code, http.StatusBadRequest)
	if len(service.createCalls) != 1 {
		t.Errorf("CreatePatient store method should not be hit if data is invalid")
	}
}

type StubService struct {
	patients    map[string]ent.Patient
	createCalls []CreateRequest
	deleteCalls []string
}

func (s *StubService) getPatients() ([]*ent.Patient, error) {
	// TODO implement me
	panic("implement me")
}

func (s *StubService) getPatient(id uuid.UUID) (*ent.Patient, error) {
	if patient, ok := s.patients[id.String()]; ok {
		return &patient, nil
	}
	return nil, &errors.NotFoundError{}
}

func (s *StubService) createPatient(patient CreateRequest) (*ent.Patient, error) {
	s.createCalls = append(s.createCalls, patient)
	return &ent.Patient{}, nil
}

func (s *StubService) deletePatient(id uuid.UUID) error {
	s.deleteCalls = append(s.deleteCalls, id.String())
	return nil
}

func (s *StubService) searchPatientsByName(name string) ([]*httpx.SearchResult, error) {
	return nil, nil
}

func (s *StubService) getPatientsTotalCount() (int, error) {
	return 0, nil
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
