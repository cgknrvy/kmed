package patient

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"kmed/api/ent"
	"kmed/api/ent/patient"
	"kmed/api/internal/database"
	"kmed/api/internal/date"
	"kmed/api/internal/errors"
	"kmed/api/internal/httpx"
)

type Service interface {
	getPatient(id uuid.UUID) (*ent.Patient, error)
	createPatient(patient CreateRequest) (*ent.Patient, error)
	deletePatient(id uuid.UUID) error
	getPatients() ([]*ent.Patient, error)
	searchPatientsByName(name string) ([]*httpx.SearchResult, error)
	getPatientsTotalCount() (int, error)
	getPatientsCreatedTodayCount() (int, error)
}

type service struct {
	client *ent.Client
	db     *sql.DB
}

func NewService(database *database.Database) Service {
	return &service{client: database.Client, db: database.DB}
}

// getPatient returns patient with the given id if found
func (s service) getPatient(id uuid.UUID) (*ent.Patient, error) {
	p, err := s.client.Patient.Get(context.Background(), id)
	if err != nil {
		return nil, errors.PatientError(err)
	}
	return p, nil
}

// getPatients returns all patients
func (s service) getPatients() ([]*ent.Patient, error) {
	patients, err := s.client.Patient.Query().All(context.Background())
	if err != nil {
		return nil, errors.PatientError(err)
	}
	return patients, nil
}

// searchPatientsByName searches for patients by mathing their stored name
// to the passed name. If any patients are found, then they are returned.
func (s service) searchPatientsByName(name string) ([]*httpx.SearchResult, error) {
	rows, err := s.db.QueryContext(context.Background(), `
			SELECT id, name
			FROM patient_search
			WHERE patient_search MATCH ?
			LIMIT 10
			`,
		name)
	if err != nil {
		return []*httpx.SearchResult{}, err
	}

	var results []*httpx.SearchResult

	for rows.Next() {
		var result httpx.SearchResult
		rows.Scan(&result.ID, &result.Name)
		results = append(results, &result)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return results, nil
}

// getPatientsTotalCount returns the total number of patients
func (s service) getPatientsTotalCount() (int, error) {
	count, err := s.client.Patient.Query().Count(context.Background())
	if err != nil {
		return 0, errors.PatientError(err)
	}

	return count, nil
}

// getPatientsCreatedTodayCount returns number of patients created in the last
// 24hrs.
func (s service) getPatientsCreatedTodayCount() (int, error) {
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	count, err := s.client.Patient.Query().
		Where(patient.CreatedAtGTE(yesterday)).
		Count(context.Background())
	if err != nil {
		return 0, errors.PatientError(err)
	}
	return count, nil
}

// createPatient creates a new patient with the provided patient data
func (s service) createPatient(patient CreateRequest) (*ent.Patient, error) {
	// Request validation is already done by the Parse function

	fmt.Println(patient)
	// Build the create query
	// Add the required fields
	createQuery := s.client.Patient.Create().SetName(patient.Name).
		SetGender(patient.Gender).SetMaritalStatus(patient.MaritalStatus)

	// Add optional fields
	if patient.Email != "" {
		createQuery.SetEmail(patient.Email)
	}
	if patient.PhoneNumber != "" {
		createQuery.SetPhoneNumber(patient.PhoneNumber)
	}
	if patient.DateOfBirth != nil {
		// Age can be calculated from the date of birth.
		createQuery.SetDob(patient.DateOfBirth.Time)
	}
	if patient.KnownAllergies != "" {
		createQuery.SetKnownAllergies(patient.KnownAllergies)
	}
	if patient.PreExistingConditions != "" {
		createQuery.SetPreExistingConditions(patient.PreExistingConditions)
	}

	createdPatient, err := createQuery.Save(context.Background())
	if err != nil {
		return nil, errors.PatientError(err)
	}

	return createdPatient, nil
}

// deletePatient deletes patient mathing the given id
func (s service) deletePatient(id uuid.UUID) error {
	err := s.client.Patient.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return errors.PatientError(err)
	}

	return nil
}

// CreateRequest represents the patient data that is needed when creating a new
// patient.
type CreateRequest struct {
	Name                  string                `json:"name"                    validate:"required"`
	Email                 string                `json:"email,omitempty"         validate:"omitempty,email"`
	PhoneNumber           string                `json:"phone_number,omitempty"  validate:"omitempty,min=10,max=13"`
	Gender                patient.Gender        `json:"gender"                  validate:"omitempty,oneof=male female"`
	MaritalStatus         patient.MaritalStatus `json:"marital_status"          validate:"omitempty,oneof=married single unspecified"`
	DateOfBirth           *date.Date            `json:"dob,omitempty"           validate:"omitempty,omitnil"`
	KnownAllergies        string                `json:"known_allergies"         validate:"omitempty,omitnil"`
	PreExistingConditions string                `json:"pre_existing_conditions" validate:"omitempty,omitnil"`
}

func (r CreateRequest) Validate() error {
	return httpx.Validator.Struct(r)
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}

func (r DeleteRequest) Validate() error {
	return httpx.Validator.Struct(r)
}
