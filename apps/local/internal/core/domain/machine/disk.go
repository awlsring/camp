package machine

type Disk struct {
	Device     DiskIdentifier
	Model      *string
	Vendor     *string
	Interface  DiskInterface
	Type       DiskClass
	Serial     *string
	SectorSize int
	Size       int64
	SizeRaw    *int64
}

func NewDisk(device DiskIdentifier, model *string, vendor *string, interface_ DiskInterface, type_ DiskClass, serial *string, sectorSize int, size int64, sizeRaw *int64) *Disk {
	return &Disk{
		Device:     device,
		Model:      model,
		Vendor:     vendor,
		Interface:  interface_,
		Type:       type_,
		Serial:     serial,
		SectorSize: sectorSize,
		Size:       size,
		SizeRaw:    sizeRaw,
	}
}
