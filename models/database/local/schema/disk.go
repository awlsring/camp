package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Disk struct {
	ent.Schema
}

func (Disk) Fields() []ent.Field {
	return []ent.Field{
		field.String("device"),
		field.String("model").Optional(),
		field.String("vendor").Optional(),
		field.String("serial").Optional(),
		field.Int("sectorSize"),
		field.Int64("size"),
		field.Int64("sizeRaw"),
		field.Enum("interface").Values("SATA", "SCSI", "PCIe", "Unknown").Default("Unknown"),
		field.Enum("diskType").Values("HDD", "SSD", "Unknown").Default("Unknown"),
	}
}

func (Disk) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("machine", Machine.Type).
			Ref("disks").
			Unique().
			Required(),
	}
}
