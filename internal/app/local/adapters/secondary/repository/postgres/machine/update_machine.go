package machine_repository

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
)

func (r *Repository) Update(ctx context.Context, m *machine.Machine) error {
	mid, err := r.Get(ctx, m.Identifier)
	if err != nil {
		return err
	}

	if m.Host != nil {
		err = r.updateHostModel(ctx, mid.Identifier, m.Host)
		if err != nil {
			return err
		}
	}

	if m.Cpu != nil {
		err = r.updateCpuModel(ctx, mid.Identifier, m.Cpu)
		if err != nil {
			return err
		}
	}

	if m.Memory != nil {
		err = r.updateMemoryModel(ctx, mid.Identifier, m.Memory)
		if err != nil {
			return err
		}
	}

	for _, d := range m.Disks {
		err = r.updateDiskModel(ctx, mid.Identifier, d)
		if err != nil {
			return err
		}
	}

	for _, ni := range m.NetworkInterfaces {
		err = r.updateNetworkInterfaceModel(ctx, mid.Identifier, ni)
		if err != nil {
			return err
		}
	}

	for _, v := range m.Volumes {
		err = r.updateVolumeModel(ctx, mid.Identifier, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) updateHostModel(ctx context.Context, mid machine.Identifier, m *host.Host) error {
	_, err := r.database.ExecContext(ctx, "UPDATE host_models SET family = $1, kernel_version = $2, os = $3, os_version = $4, os_platform = $5, hostname = $6 WHERE machine_id = $7", m.OS.Family, m.OS.Kernel, m.OS, m.OS.Version, m.OS.Platform, m.Hostname, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) updateCpuModel(ctx context.Context, mid machine.Identifier, m *cpu.CPU) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "UPDATE cpu_models SET total_cores = $1, total_threads = $2, architecture = $3, model = $4, vendor = $5 WHERE machine_id = $6 RETURNING id", m.TotalCores, m.TotalThreads, m.Architecture, m.Model, m.Vendor, mid)

	for _, p := range m.Processors {
		err = r.updateProcessorModel(ctx, id, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) updateProcessorModel(ctx context.Context, cpuId int64, m *cpu.Processor) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "UPDATE processor_models SET identifier = $1, core_count = $2, thread_count = $3, model = $4, vendor = $5 WHERE cpu_id = $6 RETURNING id", m.Id, m.CoreCount, m.ThreadCount, m.Model, m.Vendor, cpuId)

	for _, c := range m.Cores {
		err = r.updateCoreModel(ctx, id, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) updateCoreModel(ctx context.Context, pId int64, m *cpu.Core) error {
	_, err := r.database.ExecContext(ctx, "UPDATE core_models SET identifier = $1, threads = $2 WHERE processor_id = $3", m.Id, m.Threads, pId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) updateMemoryModel(ctx context.Context, mid machine.Identifier, m *memory.Memory) error {
	_, err := r.database.ExecContext(ctx, "UPDATE memory_models SET total = $1 WHERE machine_id = $2", m.Total, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) updateDiskModel(ctx context.Context, mid machine.Identifier, m *storage.Disk) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "UPDATE disk_models SET name = $1, size = $2, drive_type = $3, storage_controller = $4, removable = $5, vendor = $6, model = $7, serial = $8, wwn = $9 WHERE machine_id = $10 RETURNING id", m.Name, m.Size, m.DriveType, m.StorageController, m.Removable, m.Vendor, m.Model, m.Serial, m.WWN, mid)
	for _, p := range m.Partitions {
		err = r.updatePartitionModel(ctx, id, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) updatePartitionModel(ctx context.Context, id int64, m *storage.Partition) error {
	_, err := r.database.ExecContext(ctx, "UPDATE partition_models SET name = $1, size = $2, file_system = $3, mount_point = $4 WHERE disk_id = $5", m.Name, m.Size, m.FileSystem, m.MountPoint, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) updateNetworkInterfaceModel(ctx context.Context, mid machine.Identifier, m *network.Nic) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "UPDATE network_interface_models SET name = $1, virtual = $2, mac_address = $3, vendor = $4, speed = $5, duplex = $6, pci_address = $7 WHERE machine_id = $8 RETURNING id", m.Name, m.Virtual, m.MacAddress, m.Vendor, m.Speed, m.Duplex, m.PCIAddress, mid)
	if err != nil {
		return err
	}
	for _, a := range m.IpAddresses {
		err = r.updateAddressModel(ctx, id, a)
	}
	return nil
}

func (r *Repository) updateAddressModel(ctx context.Context, nid int64, m *network.IpAddress) error {
	_, err := r.database.ExecContext(ctx, "UPDATE address_models SET address = $1, version = $2 WHERE nic_id = $3", m.Address, m.Version.String(), nid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) updateVolumeModel(ctx context.Context, mid machine.Identifier, m *storage.Volume) error {
	_, err := r.database.ExecContext(ctx, "UPDATE volume_models SET name = $1, mount_point = $2, total = $3, file_system = $4 WHERE machine_id = $5", m.Name, m.MountPoint, m.Total, m.FileSystem, mid)
	if err != nil {
		return err
	}
	return nil
}
