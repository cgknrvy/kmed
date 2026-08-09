package schema

import (
	"fmt"
	"net/mail"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Patient struct {
	ent.Schema
}

// Fields of the Patient.
func (Patient) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(func() uuid.UUID { return uuid.Must(uuid.NewV7()) }),
		field.String("name"),
		field.String("phone_number").Optional().Validate(func(s string) error {
			if (len(s) > 0 && len(s) < 10) || len(s) > 13 {
				return fmt.Errorf("invalid phone number: %s", s)
			}
			return nil
		}),
		field.Enum("gender").Values("male", "female", "unspecified").Default("unspecified"),
		field.Enum("marital_status").
			Values("married", "single", "unspecified").
			Default("unspecified"),
		field.Time("dob").Optional(),
		field.String("email").Optional().Unique().Validate(func(s string) error {
			if _, err := mail.ParseAddress(s); err != nil {
				return EmailValidationError{invalidEmail: s}
			}
			return nil
		}),
	}
}

// Edges of the Patient.
func (Patient) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("consultations", Consultation.Type).
			StorageKey(edge.Column("patient_id")).
			Annotations(entsql.OnDelete(entsql.Cascade)), // Delete associated consultations on patient delete
	}
}

func (Patient) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
