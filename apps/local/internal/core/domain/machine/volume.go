package machine

type Volume struct {
	Name       VolumeIdentifier
	MountPoint VolumeMountPoint
	Total      int64
	FileSystem *string
}

func NewVolume(name VolumeIdentifier, mountPoint VolumeMountPoint, total int64, fileSystem *string) Volume {
	return Volume{
		Name:       name,
		MountPoint: mountPoint,
		Total:      total,
		FileSystem: fileSystem,
	}
}
