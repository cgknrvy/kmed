package patient

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"kmed/api/ent"
	"kmed/api/ent/patient"
	"kmed/api/internal/database"
	"kmed/api/internal/date"
	kmederrors "kmed/api/internal/errors"
)

func TestService_GetPatient(t *testing.T) {
	database, mockPatient := createMockPatient(t)
	svc := NewService(database)

	t.Run("getting existing patient", func(t *testing.T) {
		gottenPatient, err := svc.getPatient(mockPatient.ID)
		if assert.NoError(t, err, "expected no error when getting an existing patient") {
			assert.EqualValues(
				t,
				Patient{}.fromEntPatient(mockPatient),
				Patient{}.fromEntPatient(gottenPatient),
				"expected both patients to be equal",
			)
		}
	})

	t.Run("getting not existing patient", func(t *testing.T) {
		gottenPatient, err := svc.getPatient(uuid.Must(uuid.NewV7()))
		assert.Nil(t, gottenPatient)
		ok := kmederrors.IsNotFound(err)
		assert.True(t, ok, "expected error to be a NotFound")
	})
}

func TestService_CreatePatient(t *testing.T) {
	now := date.Today()

	t.Run("creating a valid patient", func(t *testing.T) {
		database := newTestDatabase(t)
		svc := NewService(database)
		newPatient := CreateRequest{
			Name:          "Patient001",
			Email:         "patient@patient.com",
			PhoneNumber:   "0712345678",
			Gender:        "male",
			MaritalStatus: "married",
			DateOfBirth:   &now,
		}

		createdPatient, err := svc.createPatient(newPatient)
		assert.Nil(t, err)
		assert.NotNil(t, createdPatient)

		gottenPatient, err := database.Client.Patient.Query().
			Where(patient.EmailEQ("patient@patient.com")).
			Only(context.Background())
		assert.Nil(t, err)
		assert.NotNil(t, gottenPatient)

		assert.EqualValues(
			t,
			Patient{}.fromEntPatient(createdPatient),
			Patient{}.fromEntPatient(gottenPatient),
			"expected both patients to be equal",
		)
	})

	tests := []struct {
		name          string
		field         string
		createRequest CreateRequest
	}{
		{
			name:  "invalid email",
			field: "email",
			createRequest: CreateRequest{
				Name:          "Patient001",
				Email:         "invalidemail.com",
				DateOfBirth:   &now,
				Gender:        "male",
				MaritalStatus: "married",
			},
		},
		{
			name:  "invalid phone",
			field: "phone_number",
			createRequest: CreateRequest{
				Name:          "Patient001",
				PhoneNumber:   "07123",
				DateOfBirth:   &now,
				Gender:        "male",
				MaritalStatus: "married",
			},
		},
		{
			name:  "invalid gender",
			field: "gender",
			createRequest: CreateRequest{
				Name:          "Patient001",
				DateOfBirth:   &now,
				Gender:        "invalid",
				MaritalStatus: "married",
			},
		},
		{
			name:  "invalid marital status",
			field: "marital_status",
			createRequest: CreateRequest{
				Name:          "Patient001",
				DateOfBirth:   &now,
				Gender:        "male",
				MaritalStatus: "unmarried",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			database := newTestDatabase(t)
			svc := NewService(database)

			createdPatient, err := svc.createPatient(test.createRequest)
			assert.Nil(t, createdPatient)
			assert.NotNil(t, err)

			ok := kmederrors.IsBadRequest(err)
			assert.True(t, ok, "expected error to be a BadRequest")
		})
	}

	t.Run("creating an already existing patient", func(t *testing.T) {
		database, mockPatient := createMockPatient(t)
		svc := NewService(database)
		now := date.Today()

		tests := []struct {
			name          string
			field         string
			createRequest CreateRequest
		}{
			{
				name:  "already existing email",
				field: "email",
				createRequest: CreateRequest{
					Name:          "Patient001",
					Email:         mockPatient.Email,
					DateOfBirth:   &now,
					Gender:        "male",
					MaritalStatus: "married",
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				createdPatient, err := svc.createPatient(test.createRequest)
				assert.Nil(t, createdPatient)
				assert.NotNil(t, err)

				ok := kmederrors.IsBadRequest(err)
				assert.True(t, ok, "expected error to be a BadRequest")
			})
		}
	})
}

func TestService_DeletePatient(t *testing.T) {
	t.Run("delete an existing patient", func(t *testing.T) {
		database, mockPatient := createMockPatient(t)
		svc := NewService(database)

		err := svc.deletePatient(mockPatient.ID)
		assert.Nil(t, err)
		gottenPatient, err := database.Client.Patient.Get(context.Background(), mockPatient.ID)

		assert.Nil(t, gottenPatient)
		assert.NotNil(t, err)

		ok := kmederrors.IsNotFound(kmederrors.PatientError(err))
		assert.True(t, ok, "expected error to be a NotFound")
	})

	t.Run("delete a non existing patient", func(t *testing.T) {
		database := newTestDatabase(t)
		svc := NewService(database)

		err := svc.deletePatient(uuid.Must(uuid.NewV7()))
		assert.NotNil(t, err)

		ok := kmederrors.IsNotFound(err)
		assert.True(t, ok, "expected error to be a NotFound")
	})
}

func newTestDatabase(t *testing.T) *database.Database {
	t.Helper()

	database, err := database.OpenDatabase("file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	err = database.Migrate(context.Background())
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	t.Cleanup(func() {
		_ = database.Client.Close()
	})

	return database
}

func createMockPatient(t *testing.T) (*database.Database, *ent.Patient) {
	t.Helper()
	database := newTestDatabase(t)

	// Create dummy user
	mockUser, err := database.Client.Patient.Create().
		SetName("Patient0").
		SetEmail("test@patient.com").
		SetPhoneNumber("0712345678").
		SetGender(patient.GenderMale).
		SetMaritalStatus(patient.MaritalStatusSingle).
		SetDob(date.Today().Time).
		Save(context.Background())
	if err != nil {
		t.Fatalf("failed to create mock patient: %v", err)
	}

	return database, mockUser
}

type Patient struct {
	ID            uuid.UUID
	Name          string
	Email         string
	PhoneNumber   string
	Gender        patient.Gender
	MaritalStatus patient.MaritalStatus
	DateOfBirth   string
}

func (p Patient) fromEntPatient(patient *ent.Patient) Patient {
	return Patient{
		ID:            patient.ID,
		Name:          patient.Name,
		Email:         patient.Email,
		PhoneNumber:   patient.PhoneNumber,
		Gender:        patient.Gender,
		MaritalStatus: patient.MaritalStatus,
		DateOfBirth:   patient.Dob.Format(time.DateOnly),
	}
}
