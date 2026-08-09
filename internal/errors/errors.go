package errors

import (
	"errors"
	"fmt"
	"strings"

	"kmed/api/ent"
)

type Model string

const (
	PatientModel      Model = "patient"
	UserModel         Model = "user"
	ConsultationModel Model = "consultation"
)

// ConsultationError produces a suited error for the consultation model.
func ConsultationError(err error) error {
	return matchError(err, ConsultationModel)
}

// PatientError produces a suited error for the patient model.
func PatientError(err error) error {
	return matchError(err, PatientModel)
}

// UserError produces a suited error for the user model.
func UserError(err error) error {
	return matchError(err, UserModel)
}

func matchError(err error, model Model) error {
	switch {
	case ent.IsNotFound(err):
		return &NotFoundError{msg: fmt.Sprintf("%s not found", model)}
	case ent.IsValidationError(err):
		var e *ent.ValidationError
		errors.As(err, &e)
		badRequestErr := &BadRequestError{model: model, kind: Validation}
		badRequestErr.fromValidationError(e)
		return badRequestErr
	case ent.IsConstraintError(err):
		var e *ent.ConstraintError
		errors.As(err, &e)
		badRequestErr := &BadRequestError{model: model, kind: Constraint}
		badRequestErr.fromConstraintError(e)
		return badRequestErr
	default:
		return &InternalServerError{msg: err.Error()}
	}
}

type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string { return e.msg }

// Is implements the method for comparison of NotFoundError
// used when calling the errors.Is function
func (e *NotFoundError) Is(target error) bool {
	var t *NotFoundError
	ok := errors.As(target, &t)
	return ok && t.msg == e.msg
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	var e *NotFoundError
	return errors.As(err, &e)
}

type BadRequestErrorKind string

const (
	Constraint BadRequestErrorKind = "constraint"
	Validation BadRequestErrorKind = "validation"
)

type BadRequestError struct {
	model Model
	field string
	kind  BadRequestErrorKind
}

func (e *BadRequestError) Error() string {
	switch e.kind {
	case Constraint:
		return fmt.Sprintf("%s with given %s already exists", e.model, e.field)
	case Validation:
		return fmt.Sprintf("Invalid %s for the %s", e.field, e.model)

	}
	return fmt.Sprintf("Cannot create %s with given data", e.model)
}

// fromConstraintError updates the BadRequestError from the given [ent.ConstraintError]
func (e *BadRequestError) fromConstraintError(err *ent.ConstraintError) {
	// Get the model.field pair from the error
	modelAndField := strings.Split(err.Error(), ":")[3]
	// Get the field name
	field := strings.Split(modelAndField, ".")[1]

	e.field = field
	e.kind = Constraint
}

// fromValidationError updates the BadRequestError with data from [ent.ValidationError]
func (e *BadRequestError) fromValidationError(err *ent.ValidationError) {
	e.field = err.Name
	e.kind = Validation
}

// Is implements the method for comparison of BadRequestError
// used when calling the [errors.Is] function
func (e *BadRequestError) Is(target error) bool {
	var t *BadRequestError
	ok := errors.As(target, &t)
	return ok && t.field == e.field
}

func IsBadRequest(err error) bool {
	if err == nil {
		return false
	}
	var e *BadRequestError
	return errors.As(err, &e)
}

type InternalServerError struct {
	msg string
}

func (e *InternalServerError) Error() string { return e.msg }

// Is implements the method for comparison of InternalServerError
// used when calling the [errors.Is] function
func (e *InternalServerError) Is(target error) bool {
	var t *InternalServerError
	ok := errors.As(target, &t)
	return ok && t.msg == e.msg
}

func IsInternalServerError(err error) bool {
	if err == nil {
		return false
	}
	var e *InternalServerError
	return errors.As(err, &e)
}
