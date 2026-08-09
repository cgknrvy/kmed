package patient

import (
	"errors"
	"fmt"
	"strings"

	"kmed/api/ent"
)

func matchError(err error) error {
	switch {
	case ent.IsNotFound(err):
		return &NotFoundError{msg: "patient not found"}
	case ent.IsValidationError(err):
		var e *ent.ValidationError
		errors.As(err, &e)
		return (&BadRequestError{}).fromValidationError(e)
	case ent.IsConstraintError(err):
		var e *ent.ConstraintError
		errors.As(err, &e)
		return (&BadRequestError{}).fromConstraintError(e)
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
	field string
	kind  BadRequestErrorKind
}

func (e *BadRequestError) Error() string {
	switch e.kind {
	case Constraint:
		return fmt.Sprintf("Patient with given %s already exists", e.field)
	case Validation:
		return fmt.Sprintf("Invalid %s for the patient", e.field)

	}
	return "Cannot create patient with given data"
}

// fromConstraintError creates a new BadRequestError from the given [ent.ConstraintError]
func (e *BadRequestError) fromConstraintError(err *ent.ConstraintError) *BadRequestError {
	// Get the model.field pair from the error
	modelAndField := strings.Split(err.Error(), ":")[3]
	// Get the field name
	field := strings.Split(modelAndField, ".")[1]

	return &BadRequestError{field: field, kind: Constraint}
}

// fromValidationError creates a new BadRequestError from the given [ent.ValidationError]
func (e *BadRequestError) fromValidationError(err *ent.ValidationError) *BadRequestError {
	return &BadRequestError{
		field: err.Name,
		kind:  Validation,
	}
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
