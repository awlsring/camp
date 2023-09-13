package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type IpAddress struct {
	ent.Schema
}

func (IpAddress) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("version").Values("v6", "v4", "Unknown").Default("Unknown"),
		field.String("address"),
	}
}

func (IpAddress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("networkInterface", NetworkInterface.Type).
			Ref("ipAddresses").
			Unique(),
	}
}
