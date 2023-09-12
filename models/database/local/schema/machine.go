package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Machine struct {
	ent.Schema
}

func (Machine) Fields() []ent.Field {
	return []ent.Field{
		field.String("identifier").Unique().Immutable(),
		field.Enum("state").Values("Running", "Starting", "Stopping", "Stopped", "Restarting", "Unreachable", "Unknown").Default("Running"),
		field.Enum("class").Values("BareMetal", "VirtualMachine", "Hypervisor", "Unknown").Default("Unknown"),
		field.Strings("tags"),
		field.Time("lastHeartbeat").Default(time.Now),
		field.Time("registeredAt").Default(time.Now),
		field.Time("updatedAt").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Machine) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("system", System.Type).Unique(),
		edge.To("cpu", CPU.Type).Unique(),
		edge.To("disks", Disk.Type).Unique(),
		edge.To("memory", Memory.Type).Unique(),
		edge.To("networkInterfaces", NetworkInterface.Type).Unique(),
		edge.To("volumes", Volume.Type).Unique(),
	}
}
