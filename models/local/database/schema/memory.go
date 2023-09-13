package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Memory struct {
	ent.Schema
}

func (Memory) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("total"),
	}
}

func (Memory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("machine", Machine.Type).
			Ref("memory").
			Unique().
			Required(),
	}
}
