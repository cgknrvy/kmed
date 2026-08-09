package patient

import (
	"context"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"kmed/api/ent"
	"kmed/api/ent/patient"
)

func TestMatchError(t *testing.T) {
	database := newTestDatabase(t)
	var e1 *ent.NotFoundError
	errors.As(notFoundError(t, database.Client), &e1)

	var e2 *ent.ConstraintError
	errors.As(constraintError(t, database.Client), &e2)

	var e3 *ent.ValidationError
	errors.As(validationError(t, database.Client), &e3)

	e4 := errors.New("test error")

	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "valid random error",
			err:      e4,
			expected: &InternalServerError{msg: e4.Error()},
		},
		{
			name:     "ent.NotFoundError",
			err:      e1,
			expected: &NotFoundError{msg: "patient not found"},
		},
		{
			name: "ent.ConstraintError",
			err:  e2,
			expected: &BadRequestError{
				field: "email",
				kind:  Constraint,
			},
		},
		{
			name: "ent.ValidationError",
			err:  e3,
			expected: &BadRequestError{
				field: e3.Name,
				kind:  Validation,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := matchError(test.err)
			assert.ErrorIs(t, err, test.expected)
		})
	}
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
				field: "email",
				kind:  Validation,
			},
			expected: true,
		},
		{
			name: "Valid BadRequestError from a constraint error",
			err: &BadRequestError{
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

func notFoundError(t *testing.T, database *ent.Client) error {
	t.Helper()

	_, err := database.Patient.Query().
		Where(patient.EmailEQ("invalidemail@maily.com")).
		First(context.Background())

	if ent.IsNotFound(err) {
		return err
	}

	t.Fatal("expected error to be a not found error")
	return nil
}

func constraintError(t *testing.T, database *ent.Client) error {
	t.Helper()
	_, err := database.Patient.Create().
		SetName("one").
		SetEmail("email@mail.com").
		Save(context.Background())
	assert.Nil(t, err, "expected patient to be created")
	// Repeat same request with same value for email which must be unique
	_, err = database.Patient.Create().
		SetName("two").
		SetEmail("email@mail.com").
		Save(context.Background())

	if ent.IsConstraintError(err) {
		return err
	}
	t.Fatal("expected error to be a constraint error")
	return nil
}

func validationError(t *testing.T, database *ent.Client) error {
	t.Helper()
	// Create user with an invalid email
	_, err := database.Patient.Create().
		SetName("three").
		SetEmail("invalidmail.com").
		Save(context.Background())

	if ent.IsValidationError(err) {
		return err
	}
	t.Fatal("expected error to be a validation error")
	return nil
}
