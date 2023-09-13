package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type NetworkInterface struct {
	ent.Schema
}

func (NetworkInterface) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Bool("virtual"),
		field.String("macAddress").Optional(),
		field.String("vendor").Optional(),
		field.Int("mtu").Optional(),
		field.String("duplex").Optional(),
		field.Int("speed").Optional(),
	}
}

func (NetworkInterface) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("machine", Machine.Type).
			Ref("networkInterfaces").
			Unique().
			Required(),
		edge.To("ipAddresses", IpAddress.Type),
	}
}
