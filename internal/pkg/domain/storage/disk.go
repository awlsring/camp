package storage

import "github.com/awlsring/camp/internal/pkg/values"

type Disk struct {
	Name              string
	Size              uint64
	DriveType         DriveType
	StorageController StorageController
	Removable         bool
	Partitions        []*Partition
	Vendor            *string
	Model             *string
	Serial            *string
	WWN               *string
}

func NewDisk(name string, size uint64, driveType DriveType, storageController StorageController, removable bool, vendor string, model string, serial string, wwn string) *Disk {
	return &Disk{
		Name:              name,
		Size:              size,
		DriveType:         driveType,
		StorageController: storageController,
		Removable:         removable,
		Vendor:            values.ParseOptional(vendor),
		Model:             values.ParseOptional(model),
		Serial:            values.ParseOptional(serial),
		WWN:               values.ParseOptional(wwn),
	}
}

func (d *Disk) AddPartitions(partitions ...*Partition) {
	d.Partitions = append(d.Partitions, partitions...)
}
