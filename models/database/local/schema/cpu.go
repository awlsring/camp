package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type CPU struct {
	ent.Schema
}

func (CPU) Fields() []ent.Field {
	return []ent.Field{
		field.Int("cores"),
		field.Enum("architecture").Values("x86", "arm", "Unknown").Default("x86"),
		field.String("model").Optional(),
		field.String("vendor").Optional(),
	}
}

func (CPU) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("machine", Machine.Type).
			Ref("cpu").
			Unique().
			Required(),
	}
}
