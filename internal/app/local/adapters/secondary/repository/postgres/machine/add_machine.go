package machine_repository

import (
	"context"
	"time"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/network"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/storage"
	"github.com/awlsring/camp/internal/pkg/domain/tag"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *Repository) Add(ctx context.Context, m *machine.Machine) error {
	err := r.createMachineEntry(ctx, m)
	if err != nil {
		return err
	}

	err = r.createPowerStateModel(ctx, m.Identifier, power.StatusCodeRunning)
	if err != nil {
		return err
	}

	err = r.createPowerCapabilityModel(ctx, m.Identifier, m.PowerCapabilities)
	if err != nil {
		return err
	}

	if m.Host != nil {
		err = r.createHostModel(ctx, m.Identifier, m.Host)
		if err != nil {
			return err
		}
	}

	if m.Cpu != nil {
		err = r.createCpuModel(ctx, m.Identifier, m.Cpu)
		if err != nil {
			return err
		}
	}

	if m.Memory != nil {
		err = r.createMemoryModel(ctx, m.Identifier, m.Memory)
		if err != nil {
			return err
		}
	}

	for _, d := range m.Disks {
		err = r.createDiskModel(ctx, m.Identifier, d)
		if err != nil {
			return err
		}
	}

	for _, ni := range m.NetworkInterfaces {
		err = r.createNetworkInterfaceModel(ctx, m.Identifier, ni)
		if err != nil {
			return err
		}
	}

	for _, v := range m.Volumes {
		err = r.createVolumeModel(ctx, m.Identifier, v)
		if err != nil {
			return err
		}
	}

	for _, t := range m.Tags {
		err = r.tagRepo.AddToResource(ctx, t, m.Identifier.String(), tag.ResourceTypeMachine)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) createMachineEntry(ctx context.Context, m *machine.Machine) error {
	now := time.Now().UTC()
	_, err := r.database.ExecContext(ctx, "INSERT INTO machines (identifier, endpoint, key, class, last_heartbeat, registered_at, updated_at, status) VALUES ($1, $2, $3, $4, $5, $6)", m.Identifier.String(), m.AgentEndpoint.String(), m.AgentApiKey, m.Class.String(), now, now, now, power.StatusCodeRunning.String())
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) createPowerStateModel(ctx context.Context, mid machine.Identifier, state power.StatusCode) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO power_states (state, updated_at, machine_id) VALUES ($1,  NOW() AT TIME ZONE 'UTC', $2)", state.String(), mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) createPowerCapabilityModel(ctx context.Context, mid machine.Identifier, cap machine.PowerCapabilities) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Creating power capability model for machine %s", mid)

	rebootEnabled := cap.Reboot.Enabled

	powerOffEnabled := cap.PowerOff.Enabled

	wakeOnLanEnabled := cap.WakeOnLan.Enabled

	var wakeOnLanMac string
	if cap.WakeOnLan.MacAddress != nil {
		wakeOnLanMac = cap.WakeOnLan.MacAddress.String()
	}

	_, err := r.database.ExecContext(ctx, "INSERT INTO power_capabilities (reboot_enabled, power_off_enabled, wake_on_lan_enabled, wake_on_lan_mac, machine_id) VALUES ($1, $2, $3, $4, $5)", rebootEnabled, powerOffEnabled, wakeOnLanEnabled, wakeOnLanMac, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) createHostModel(ctx context.Context, mid machine.Identifier, m *host.Host) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO host_models (os_family, kernel, os_name, os_version, os_platform, hostname, host_id, machine_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", m.OS.Family, m.OS.Kernel, m.OS.Name, m.OS.Version, m.OS.Platform, m.Hostname, m.HostId, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) createCpuModel(ctx context.Context, mid machine.Identifier, m *cpu.CPU) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "INSERT INTO cpu_models (total_cores, total_threads, architecture, model, vendor, machine_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", m.TotalCores, m.TotalThreads, m.Architecture, m.Model, m.Vendor, mid)
	if err != nil {
		return err
	}

	for _, p := range m.Processors {
		err = r.createProcessorModel(ctx, id, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) createProcessorModel(ctx context.Context, id int64, m *cpu.Processor) error {
	var pid int64
	err := r.database.GetContext(ctx, &pid, "INSERT INTO processor_models (identifier, core_count, thread_count, model, vendor, cpu_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", m.Id, m.CoreCount, m.ThreadCount, m.Model, m.Vendor, id)
	if err != nil {
		return err
	}

	for _, c := range m.Cores {
		err = r.createCoreModel(ctx, pid, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) createCoreModel(ctx context.Context, pid int64, m *cpu.Core) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO core_models (identifier, threads, processor_id) VALUES ($1, $2, $3)", m.Id, m.Threads, pid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) createMemoryModel(ctx context.Context, mid machine.Identifier, m *memory.Memory) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO memory_models (total, machine_id) VALUES ($1, $2)", m.Total, mid)
	if err != nil {
		return err
	}
	return nil
}

// start here
func (r *Repository) createDiskModel(ctx context.Context, mid machine.Identifier, m *storage.Disk) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "INSERT INTO disk_models (name, size, drive_type, storage_controller, removable, vendor, model, serial, wwn, machine_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", m.Name, m.Size, m.DriveType.String(), m.StorageController.String(), m.Removable, m.Vendor, m.Model, m.Serial, m.WWN, mid)
	if err != nil {
		return err
	}

	for _, p := range m.Partitions {
		err = r.createPartitionModel(ctx, id, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) createPartitionModel(ctx context.Context, did int64, m *storage.Partition) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO partition_models (name, size, readonly, label, type, file_system, uuid, mount_point, disk_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.Name, m.Size, m.Readonly, m.Label, m.Type, m.FileSystem, m.UUID, m.MountPoint, did)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createNetworkInterfaceModel(ctx context.Context, mid machine.Identifier, m *network.Nic) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "INSERT INTO network_interface_models (name, virtual, mac_address, speed, duplex, pci_address, vendor, machine_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", m.Name, m.Virtual, m.MacAddress.String(), m.Speed, m.Duplex, m.PCIAddress, m.Vendor, mid)

	for _, a := range m.IpAddresses {
		err = r.createAddressModel(ctx, id, a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) createAddressModel(ctx context.Context, nid int64, m *network.IpAddress) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO address_models (address, version, nic_id) VALUES ($1, $2, $3)", m.Address, m.Version.String(), nid)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) createVolumeModel(ctx context.Context, mid machine.Identifier, m *storage.Volume) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO volume_models (name, mount_point, total, file_system, machine_id) VALUES ($1, $2, $3, $4, $5)", m.Name, m.MountPoint, m.Total, m.FileSystem, mid)
	if err != nil {
		return err
	}
	return nil
}
