package consultation

import (
	"context"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/google/uuid"

	"kmed/api/ent"
	entConsultation "kmed/api/ent/consultation"
	"kmed/api/ent/patient"
	"kmed/api/ent/schema"
	"kmed/api/ent/user"
	"kmed/api/internal/errors"
	"kmed/api/internal/httpx"
)

type Service interface {
	createConsultation(consultation CreateRequest) (*ent.Consultation, error)
	getConsultation(id uuid.UUID) (*ent.Consultation, error)
	getPatientConsultations(id uuid.UUID) ([]*ent.Consultation, error)
	getTodaysConsultationsForDoctor(id uuid.UUID) ([]*ent.Consultation, error)
}

type service struct {
	client *ent.Client
}

func newService(client *ent.Client) Service {
	return &service{client}
}

// createConsultation creates a new consultation with the given [CreateRequest] data.
// The patient and doctor are retrieved using the [CreateRequest.PatientID] and [CreateRequest.DoctorID]
// respectively. If either is not found, then an error is returned and is expected to be [errors.NotFoundError].
func (s *service) createConsultation(consultation CreateRequest) (*ent.Consultation, error) {
	// consultation is already validated by the [httpx.Parse] function on [Handler.createConsultation]

	patient, err := s.client.Patient.Get(context.Background(), consultation.PatientID)
	if err != nil {
		return nil, errors.PatientError(err)
	}
	doctor, err := s.client.User.Get(context.Background(), consultation.DoctorID)
	if err != nil {
		return nil, errors.UserError(err)
	}

	c, err := s.client.Consultation.Create().
		SetVitals(&consultation.Vitals).
		SetClinicalNotes(&consultation.ClinicalNotes).
		SetDiagnosis(&consultation.Diagnosis).
		SetPatient(patient).
		SetDoctor(doctor).
		Save(context.Background())
	if err != nil {
		return nil, errors.ConsultationError(err)
	}
	return c, nil
}

// getConsultation returns consultation with the given id.
// An error is returned if no consultation is found.
func (s *service) getConsultation(id uuid.UUID) (*ent.Consultation, error) {
	c, err := s.client.Consultation.Get(context.Background(), id)
	if err != nil {
		return nil, errors.ConsultationError(err)
	}

	return c, nil
}

// getPatientConsultations returns consultations for the patient with the passed id.
// Returns nil if there is no patient with the given id
func (s *service) getPatientConsultations(id uuid.UUID) ([]*ent.Consultation, error) {
	consultations, err := s.client.Consultation.Query().
		Where(entConsultation.HasPatientWith(patient.IDEQ(id))).
		All(context.Background())
	if err != nil {
		return nil, errors.ConsultationError(err)
	}
	return consultations, nil
}

func (s *service) getTodaysConsultationsForDoctor(id uuid.UUID) ([]*ent.Consultation, error) {
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	consultations, err := s.client.Consultation.Query().WithPatient().
		Where(entConsultation.HasDoctorWith(user.IDEQ(id)), entConsultation.UpdatedAtGT(yesterday)).
		Order(entConsultation.ByUpdatedAt(entsql.OrderDesc())).
		All(context.Background())
	if err != nil {
		return nil, errors.ConsultationError(err)
	}
	return consultations, nil
}

type CreateRequest struct {
	Vitals        schema.Vitals        `json:"vitals"`
	ClinicalNotes schema.ClinicalNotes `json:"clinicalNotes"`
	Diagnosis     schema.Diagnosis     `json:"diagnosis"`
	PatientID     uuid.UUID            `json:"patientID"     validate:"uuid"`
	DoctorID      uuid.UUID            `json:"doctorID"      validate:"uuid"`
}

func (r CreateRequest) Validate() error {
	return httpx.Validator.Struct(r)
}
