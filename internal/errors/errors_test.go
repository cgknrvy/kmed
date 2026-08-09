package errors

import (
	"context"
	"errors"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"kmed/api/ent"
	"kmed/api/ent/patient"
)

func TestMatchError(t *testing.T) {
	client := newTestClient(t)
	var e1 *ent.NotFoundError
	errors.As(notFoundError(t, client), &e1)

	var e2 *ent.ConstraintError
	errors.As(constraintError(t, client), &e2)

	var e3 *ent.ValidationError
	errors.As(validationError(t, client), &e3)

	e4 := errors.New("test error")

	tests := []struct {
		name     string
		err      error
		model    Model
		expected error
	}{
		{
			name:     "valid random error",
			err:      e4,
			model:    PatientModel,
			expected: &InternalServerError{msg: e4.Error()},
		},
		{
			name:     "ent.NotFoundError",
			err:      e1,
			model:    PatientModel,
			expected: &NotFoundError{msg: fmt.Sprintf("%s not found", PatientModel)},
		},
		{
			name:  "ent.ConstraintError",
			err:   e2,
			model: UserModel,
			expected: &BadRequestError{
				model: UserModel,
				field: "email",
				kind:  Constraint,
			},
		},
		{
			name:  "ent.ValidationError",
			err:   e3,
			model: ConsultationModel,
			expected: &BadRequestError{
				model: ConsultationModel,
				field: e3.Name,
				kind:  Validation,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := matchError(test.err, test.model)
			assert.ErrorIs(t, err, test.expected)
		})
	}
}

func TestConsultationError(t *testing.T) {
	client := newTestClient(t)

	var e *ent.ConstraintError
	errors.As(constraintError(t, client), &e)

	t.Run("ConsultationError returns an error with ConsultationModel", func(t *testing.T) {
		var err *BadRequestError
		errors.As(ConsultationError(e), &err)
		assert.Equal(t, ConsultationModel, err.model)
	})
}

func TestPatientError(t *testing.T) {
	client := newTestClient(t)

	var e *ent.ConstraintError
	errors.As(constraintError(t, client), &e)

	t.Run("PatientError returns an error with PatientModel", func(t *testing.T) {
		var err *BadRequestError
		errors.As(PatientError(e), &err)
		assert.Equal(t, PatientModel, err.model)
	})
}

func TestUserError(t *testing.T) {
	client := newTestClient(t)

	var e *ent.ConstraintError
	errors.As(constraintError(t, client), &e)

	t.Run("UserError returns an error with UserModel", func(t *testing.T) {
		var err *BadRequestError
		errors.As(UserError(e), &err)
		assert.Equal(t, UserModel, err.model)
	})
}

func TestIsBadRequest(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "Valid BadRequestError from a validation error",
			err: &BadRequestError{
				model: UserModel,
				field: "email",
				kind:  Validation,
			},
			expected: true,
		},
		{
			name: "Valid BadRequestError from a constraint error",
			err: &BadRequestError{
				model: UserModel,
				field: "email",
				kind:  Constraint,
			},
			expected: true,
		},
		{
			name:     "Invalid BadRequestError",
			err:      errors.New("not bad request"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := IsBadRequest(test.err)
			assert.Equal(t, test.expected, ok)
		})
	}
}

func TestIsInternalServerError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{{
		name:     "Valid InternalServerErrorError",
		err:      &InternalServerError{msg: "test error"},
		expected: true,
	}, {
		name:     "Invalid InternalServerErrorError",
		err:      errors.New("test error"),
		expected: false,
	}, {
		name:     "nil error",
		err:      nil,
		expected: false,
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := IsInternalServerError(test.err)
			assert.Equal(t, test.expected, ok)
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{{
		name:     "Valid NotFoundError",
		err:      &NotFoundError{msg: "test error"},
		expected: true,
	}, {
		name:     "Invalid NotFoundError",
		err:      errors.New("test error"),
		expected: false,
	}, {
		name:     "nil error",
		err:      nil,
		expected: false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok := IsNotFound(test.err)
			assert.Equal(t, test.expected, ok)
		})
	}
}

func notFoundError(t *testing.T, client *ent.Client) error {
	t.Helper()

	_, err := client.Patient.Query().
		Where(patient.EmailEQ("invalidemail@maily.com")).
		First(context.Background())

	if ent.IsNotFound(err) {
		return err
	}

	t.Fatal("expected error to be a not found error")
	return nil
}

func constraintError(t *testing.T, client *ent.Client) error {
	t.Helper()
	_, err := client.Patient.Create().
		SetName("one").
		SetEmail("email@mail.com").
		Save(context.Background())
	assert.Nil(t, err, "expected patient to be created")
	// Repeat same request with same value for email which must be unique
	_, err = client.Patient.Create().
		SetName("two").
		SetEmail("email@mail.com").
		Save(context.Background())

	if ent.IsConstraintError(err) {
		return err
	}
	t.Fatal("expected error to be a constraint error")
	return nil
}

func validationError(t *testing.T, client *ent.Client) error {
	t.Helper()
	// Create user with an invalid email
	_, err := client.Patient.Create().
		SetName("three").
		SetEmail("invalidmail.com").
		Save(context.Background())

	if ent.IsValidationError(err) {
		return err
	}
	t.Fatal("expected error to be a validation error")
	return nil
}

func newTestClient(t *testing.T) *ent.Client {
	t.Helper()

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	return client
}
