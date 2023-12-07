package machine_repository

import (
	"context"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (r *MachineRepo) Add(ctx context.Context, m *machine.Machine) error {
	err := r.createMachineEntry(ctx, m)
	if err != nil {
		return err
	}

	err = r.createPowerCapabilityModel(ctx, m.Identifier, m.PowerCapabilities)
	if err != nil {
		return err
	}

	if m.System != nil {
		err = r.createSystemModel(ctx, m.Identifier, m.System)
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

func (r *MachineRepo) createMachineEntry(ctx context.Context, m *machine.Machine) error {
	now := time.Now().UTC()
	_, err := r.database.ExecContext(ctx, "INSERT INTO machines (identifier, endpoint, key, class, last_heartbeat, registered_at, updated_at, status) VALUES ($1, $2, $3, $4, $5, $6)", m.Identifier.String(), m.AgentEndpoint.String(), m.AgentApiKey, m.Class.String(), now, now, now, machine.MachineStatusRunning.String())
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) createPowerCapabilityModel(ctx context.Context, mid machine.Identifier, cap machine.PowerCapabilities) error {
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

func (r *MachineRepo) createSystemModel(ctx context.Context, mid machine.Identifier, m *machine.System) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO system_models (family, kernel_version, os, os_version, os_pretty, hostname, machine_id) VALUES ($1, $2, $3, $4, $5, $6, $7)", m.Os.Family, m.Os.Kernel, m.Os, m.Os.Version, m.Os.PrettyName, m.Hostname, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) createCpuModel(ctx context.Context, mid machine.Identifier, m *machine.Cpu) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO cpu_models (cores, architecture, model, vendor, machine_id) VALUES ($1, $2, $3, $4, $5)", m.Cores, m.Architecture.String(), m.Model, m.Vendor, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) createMemoryModel(ctx context.Context, mid machine.Identifier, m *machine.Memory) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO memory_models (total, machine_id) VALUES ($1, $2)", m.Total, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) createDiskModel(ctx context.Context, mid machine.Identifier, m *machine.Disk) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO disk_models (device, model, vendor, interface, type, serial, sector_size, size, size_raw, machine_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", m.Device.String(), m.Model, m.Vendor, m.Interface.String(), m.Type.String(), m.Serial, m.SectorSize, m.Size, m.SizeRaw, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) createNetworkInterfaceModel(ctx context.Context, mid machine.Identifier, m *machine.NetworkInterface) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "INSERT INTO network_interface_models (name, virtual, mac_address, vendor, mtu, speed, duplex, machine_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", m.Name, m.Virtual, m.MacAddress, m.Vendor, m.Mtu, m.Speed, m.Duplex, mid)
	if err != nil {
		return err
	}

	for _, a := range m.IpAddresses {
		err = r.createAddressModel(ctx, id, a)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MachineRepo) createAddressModel(ctx context.Context, nid int64, m *machine.IpAddress) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO address_models (address, version, nic_id) VALUES ($1, $2, $3)", m.Address, m.Version.String(), nid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) createVolumeModel(ctx context.Context, mid machine.Identifier, m *machine.Volume) error {
	_, err := r.database.ExecContext(ctx, "INSERT INTO volume_models (name, mount_point, total, file_system, machine_id) VALUES ($1, $2, $3, $4, $5)", m.Name, m.MountPoint, m.Total, m.FileSystem, mid)
	if err != nil {
		return err
	}
	return nil
}
