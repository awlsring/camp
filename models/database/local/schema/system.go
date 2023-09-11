package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type System struct {
	ent.Schema
}

func (System) Fields() []ent.Field {
	return []ent.Field{
		field.String("family"),
		field.String("kernelVersion"),
		field.String("os"),
		field.String("osVersion"),
		field.String("osPretty"),
		field.String("hostname"),
	}
}

func (System) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("machine", Machine.Type).
			Ref("system").
			Unique().
			Required(),
	}
}
