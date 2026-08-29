package consultation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"kmed/api/ent"
	"kmed/api/ent/schema"
	entUser "kmed/api/ent/user"

	"github.com/stretchr/testify/assert"
)

func TestPOSTConsultation(t *testing.T) {
	svc := &StubService{}
	handler := &Handler{svc: svc, auth: authMiddleware}
	t.Run("returns accepted on POST", func(t *testing.T) {
		consultation := CreateRequest{
			Vitals: schema.Vitals{
				Temperature:      36.8,
				BloodPressure:    "120/80",
				Pulse:            72,
				OxygenSaturation: 99,
				RespiratoryRate:  17,
				Weight:           58,
			},
			ClinicalNotes: schema.ClinicalNotes{
				Complaint:           "Hurting intestines",
				History:             "Ulcers",
				ExaminationFindings: "Nothing",
			},
			Diagnosis: schema.Diagnosis{
				Primary:        []schema.ICDCode{},
				Differential:   []schema.ICDCode{},
				ManagementPlan: "Drink painkillers",
			},
		}
		request := newPostConsultationRequest(consultation)
		response := httptest.NewRecorder()
		handler.Router().ServeHTTP(response, request)

		assert.Equal(t, http.StatusAccepted, response.Code)
		if len(svc.createCalls) != 1 {
			t.Errorf(
				"consultation service should have created 1 consultation, got %d",
				len(svc.createCalls),
			)
		} else {
			assert.Equal(t, svc.createCalls[0].Vitals, consultation.Vitals)
		}
	})

	vitals := schema.Vitals{
		Temperature:      36.8,
		BloodPressure:    "120/80",
		Pulse:            72,
		OxygenSaturation: 99,
		RespiratoryRate:  17,
		Weight:           58,
	}
	clinicalNotes := schema.ClinicalNotes{
		Complaint:           "Hurting intestines",
		History:             "Ulcers",
		ExaminationFindings: "Nothing",
	}
	diagnosis := schema.Diagnosis{
		Primary:        []schema.ICDCode{},
		Differential:   []schema.ICDCode{},
		ManagementPlan: "Drink painkillers",
	}
	tests := []struct {
		name         string
		consultation CreateRequest
		expectedCode int
	}{
		{
			name:         "empty request",
			consultation: CreateRequest{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "missing vitals",
			consultation: CreateRequest{
				ClinicalNotes: clinicalNotes,
				Diagnosis:     diagnosis,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "missing clinical notes",
			consultation: CreateRequest{
				Vitals:    vitals,
				Diagnosis: diagnosis,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "missing diagnosis",
			consultation: CreateRequest{
				Vitals:        vitals,
				ClinicalNotes: clinicalNotes,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := newPostConsultationRequest(test.consultation)
			response := httptest.NewRecorder()
			handler.Router().ServeHTTP(response, request)
			assert.Equal(t, test.expectedCode, response.Code)
		})
	}
}

func newPostConsultationRequest(payload CreateRequest) *http.Request {
	data, _ := json.Marshal(payload)
	request, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(string(data)))
	request.Header.Set("Content-Type", "application/json")
	return request
}

type StubService struct {
	consultations map[string]ent.Consultation
	createCalls   []CreateRequest
}

func (s *StubService) createConsultation(consultation CreateRequest) (*ent.Consultation, error) {
	s.createCalls = append(s.createCalls, consultation)
	return &ent.Consultation{}, nil
}

func (s *StubService) getConsultation(id uuid.UUID) (*ent.Consultation, error) {
	return &ent.Consultation{}, nil
}

func (s *StubService) getPatientConsultations(id uuid.UUID) ([]*ent.Consultation, error) {
	return []*ent.Consultation{}, nil
}

func (s *StubService) getTodaysConsultationsForDoctor(id uuid.UUID) ([]*ent.Consultation, error) {
	return []*ent.Consultation{}, nil
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
