package storage

import "strings"

type StorageController int

const (
	StorageControllerUnknown StorageController = iota
	StorageControllerIde
	StorageControllerScsi
	StorageControllerNvme
	StorageControllerVirtio
	StorageControllerMmc
	StorageControllerLoop
)

func (s StorageController) String() string {
	switch s {
	case StorageControllerIde:
		return "IDE"
	case StorageControllerScsi:
		return "SCSI"
	case StorageControllerNvme:
		return "NVMe"
	case StorageControllerVirtio:
		return "VirtIO"
	case StorageControllerMmc:
		return "MMC"
	case StorageControllerLoop:
		return "Loop"
	default:
		return "Unknown"
	}
}

func StorageControllerFromString(v string) StorageController {
	switch strings.ToLower(v) {
	case "ide":
		return StorageControllerIde
	case "scsi":
		return StorageControllerScsi
	case "nvme":
		return StorageControllerNvme
	case "virtio":
		return StorageControllerVirtio
	case "mmc":
		return StorageControllerMmc
	case "loop":
		return StorageControllerLoop
	default:
		return StorageControllerUnknown
	}
}
