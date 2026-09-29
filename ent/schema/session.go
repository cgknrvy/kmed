package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Immutable().
			Default(func() uuid.UUID { return uuid.Must(uuid.NewV7()) }),
		field.UUID("user_id", uuid.UUID{}),
		field.Text("refresh_token_hash").Unique().NotEmpty().Sensitive(),
		field.Time("expires_at"),
		field.Time("revoked_at").Optional().Nillable(),
	}
}

func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Field("user_id").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the RefreshToken.
func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}

func (Session) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
