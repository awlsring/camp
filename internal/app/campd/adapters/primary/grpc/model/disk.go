package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	grpcmodel "github.com/awlsring/camp/pkg/grpc_model"
)

func NewDiskIdentifier(n string) *campd.DiskIdentifier {
	return &campd.DiskIdentifier{
		Value: n,
	}
}

func NewPartitionIdentifier(n string) *campd.PartitionIdentifier {
	return &campd.PartitionIdentifier{
		Value: n,
	}
}

func DiskTypeFromDomain(t storage.DriveType) campd.DiskType {
	switch t {
	case storage.DriveTypeFDD:
		return campd.DiskType_FDD
	case storage.DriveTypeHDD:
		return campd.DiskType_HDD
	case storage.DriveTypeSSD:
		return campd.DiskType_SSD
	case storage.DriveTypeODD:
		return campd.DiskType_ODD
	case storage.DriveTypeVirtual:
		return campd.DiskType_VIRTUAL
	default:
		return campd.DiskType_DISKTYPE_UNKNOWN
	}
}

func StorageControllerFromDomain(t storage.StorageController) campd.StorageController {
	switch t {
	case storage.StorageControllerIde:
		return campd.StorageController_IDE
	case storage.StorageControllerLoop:
		return campd.StorageController_LOOP
	case storage.StorageControllerScsi:
		return campd.StorageController_SCSI
	case storage.StorageControllerMmc:
		return campd.StorageController_MMC
	case storage.StorageControllerNvme:
		return campd.StorageController_NVME
	case storage.StorageControllerVirtio:
		return campd.StorageController_VIRTIO
	default:
		return campd.StorageController_STORAGECONTROLLER_UNKNOWN
	}
}

func DiskFromDomain(in *storage.Disk) *campd.DiskSummary {
	partitionIdentifiers := make([]*campd.PartitionIdentifier, len(in.Partitions))
	for i, p := range in.Partitions {
		partitionIdentifiers[i] = NewPartitionIdentifier(p.Name)
	}
	return &campd.DiskSummary{
		Name:              NewDiskIdentifier(in.Name),
		Size:              int64(in.Size),
		Partitions:        partitionIdentifiers,
		Type:              DiskTypeFromDomain(in.DriveType),
		StorageController: StorageControllerFromDomain(in.StorageController),
		Removable:         in.Removable,
		Serial:            grpcmodel.NewStringValue(in.Serial),
		Model:             grpcmodel.NewStringValue(in.Model),
		Vendor:            grpcmodel.NewStringValue(in.Vendor),
		Wwn:               grpcmodel.NewStringValue(in.WWN),
	}
}
