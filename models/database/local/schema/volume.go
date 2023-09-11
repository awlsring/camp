package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Volume struct {
	ent.Schema
}

func (Volume) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("mountPoint"),
		field.Int64("total"),
		field.String("fileSystem"),
	}
}

func (Volume) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("machine", Machine.Type).
			Ref("volumes").
			Unique().
			Required(),
	}
}
