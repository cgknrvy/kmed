package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Consultation holds the schema definition for the Consultation entity
type Consultation struct {
	ent.Schema
}

// Vitals holds the measurements taken during consultation
type Vitals struct {
	Temperature      float64 `json:"temperature,omitempty"     validate:"omitempty,number"`
	BloodPressure    string  `json:"bloodPressure,omitempty"   validate:"omitempty"`
	Pulse            int     `json:"pulse,omitempty"           validate:"omitempty,number"`
	OxygenSaturation float64 `json:"oxygenSat,omitempty"       validate:"omitempty,number"`
	RespiratoryRate  int     `json:"respiratoryRate,omitempty" validate:"omitempty,number"`
	Weight           float64 `json:"weight,omitempty"          validate:"omitempty,number"`
}

// ClinicalNotes holds the patient's complaint and its history
type ClinicalNotes struct {
	Complaint           string `json:"complaint"                     validate:"required"`
	History             string `json:"history"                       validate:"required"`
	ExaminationFindings string `json:"examinationFindings,omitempty" validate:"omitempty"`
}

// Diagnosis holds the final clinician's conclusion on what might be the patient's ailment
type Diagnosis struct {
	Primary        string `json:"primary"                validate:"required"`
	Differential   string `json:"differential,omitempty" validate:"omitempty"`
	Severity       string `json:"severity"               validate:"required,oneof=low medium high"`
	ManagementPlan string `json:"managementPlan"         validate:"required"`
}

// Fields of the Consultation
func (Consultation) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(func() uuid.UUID { return uuid.Must(uuid.NewV7()) }),
		field.JSON("vitals", &Vitals{}),
		field.JSON("clinical_notes", &ClinicalNotes{}),
		field.JSON("diagnosis", &Diagnosis{}),
	}
}

func (Consultation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

func (Consultation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("patient", Patient.Type).
			Ref("consultations").
			Unique().   // Ensure only one patient per consultation
			Required(), // Ensures that there is a patient when creating a consultation
		edge.From("doctor", User.Type).
			Ref("consultations").
			Unique().   // Ensure only one doctor per consultation
			Required(), // Ensure the doctor is present at consultation creation
	}
}
