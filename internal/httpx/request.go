package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Validatable interface {
	Validate() error
}

// Parse parses the body of the given request and binds it to a variable of
// type T that must implement the Validatable interface.
// The parsed body is then validated by calling the `Validate` method on the
// type T.
// An error is returned if the body cannot be decoded or if the validation fails.
func Parse[T Validatable](w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest,
			ErrorResponse{
				Message: "error decoding body",
				Errors:  map[string]string{"decode": err.Error()},
			})
		return req, false
	}
	if err := req.Validate(); err != nil {
		JSONError(w, http.StatusBadRequest,
			ErrorResponse{
				Message: "error validating request data",
				Errors:  map[string]string{"validate": err.Error()},
			})
		return req, false
	}

	return req, true
}

// GetIDFromPath returns the id linked to the `id` path value. `id` must be a UUID else
// an error is returned
func GetIDFromPath(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest,
			ErrorResponse{
				Message: fmt.Sprintf("provided id \"%s\" is not a valid UUID", idStr),
				Errors:  map[string]string{"id": err.Error()},
			})
		return id, false
	}
	return id, true
}
