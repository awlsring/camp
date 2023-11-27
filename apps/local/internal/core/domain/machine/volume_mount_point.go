package machine

type VolumeMountPoint string

func (v VolumeMountPoint) String() string {
	return string(v)
}

func MountPointFromString(v string) (VolumeMountPoint, error) {
	return NonNullString[VolumeMountPoint](v)
}
