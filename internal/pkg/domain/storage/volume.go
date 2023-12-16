package storage

import "github.com/awlsring/camp/internal/pkg/values"

type Volume struct {
	Name       string
	MountPoint string
	Total      uint64
	FileSystem *string
}

func NewVolume(name string, mountPoint string, total uint64, fileSystem string) Volume {
	return Volume{
		Name:       name,
		MountPoint: mountPoint,
		Total:      total,
		FileSystem: values.ParseOptional(fileSystem),
	}
}
