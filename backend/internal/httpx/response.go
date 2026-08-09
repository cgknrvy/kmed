package httpx

import (
	"encoding/json"
	"log"
	"net/http"

	"kmed/api/ent"

	"github.com/google/uuid"
)

type SearchResult struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`
}

type Response struct {
	Message       *string             `json:"message,omitempty"`
	User          *ent.User           `json:"user,omitempty"`
	AccessToken   *string             `json:"accessToken,omitempty"`
	Patient       *ent.Patient        `json:"patient,omitempty"`
	Patients      []*ent.Patient      `json:"patients,omitempty"`
	SearchResults []*SearchResult     `json:"searchResults,omitempty"`
	Consultation  *ent.Consultation   `json:"consultation,omitempty"`
	Consultations []*ent.Consultation `json:"consultations,omitempty"`
	Count         *int                `json:"count,omitempty"`
}

func JSON(w http.ResponseWriter, code int, response Response) {
	jsonResponse[Response](w, code, response)
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func JSONError(w http.ResponseWriter, code int, response ErrorResponse) {
	jsonResponse[ErrorResponse](w, code, response)
}

func jsonResponse[T Response | ErrorResponse](w http.ResponseWriter, code int, response T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("error encoding response: ", err)
	}
}
