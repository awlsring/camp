package storage

import "github.com/awlsring/camp/internal/pkg/values"

type Partition struct {
	Name       string
	Size       uint64
	Disk       *Disk
	Readonly   bool
	Label      *string
	Type       *string
	FileSystem *string
	UUID       *string
	MountPoint *string
}

func NewPartition(name string, size uint64, disk *Disk, readonly bool, label string, t string, fs string, uuid string, mountPoint string) *Partition {
	return &Partition{
		Name:       name,
		Size:       size,
		Disk:       disk,
		Readonly:   readonly,
		Label:      values.ParseOptional(label),
		Type:       values.ParseOptional(t),
		FileSystem: values.ParseOptional(fs),
		UUID:       values.ParseOptional(uuid),
		MountPoint: values.ParseOptional(mountPoint),
	}
}
