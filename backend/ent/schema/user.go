package schema

import (
	"fmt"
	"net/mail"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type EmailValidationError struct {
	invalidEmail string
}

func (e EmailValidationError) Error() string {
	return fmt.Sprintf("invalid user email: %s", e.invalidEmail)
}

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(func() uuid.UUID { return uuid.Must(uuid.NewV7()) }),
		field.String("name"),
		field.String("email").NotEmpty().Unique().Validate(func(s string) error {
			if _, err := mail.ParseAddress(s); err != nil {
				return EmailValidationError{invalidEmail: s}
			}
			return nil
		}),
		field.String("password").MinLen(8).Sensitive(),
		field.Enum("role").Values("admin", "doctor", "lab-tech", "user"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("consultations", Consultation.Type).
			StorageKey(edge.Column("doctor_id")),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
