package machine_repository

import (
	"context"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

func (r *MachineRepo) Update(ctx context.Context, m *machine.Machine) error {
	mid, err := r.Get(ctx, m.Identifier)
	if err != nil {
		return err
	}

	if m.System != nil {
		err = r.updateSystemModel(ctx, mid.Identifier, m.System)
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

func (r *MachineRepo) updateSystemModel(ctx context.Context, mid machine.Identifier, m *machine.System) error {
	_, err := r.database.ExecContext(ctx, "UPDATE system_models SET family = $1, kernel_version = $2, os = $3, os_version = $4, os_pretty = $5, hostname = $6 WHERE machine_id = $7", m.Os.Family, m.Os.Kernel, m.Os, m.Os.Version, m.Os.PrettyName, m.Hostname, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) updateCpuModel(ctx context.Context, mid machine.Identifier, m *machine.Cpu) error {
	_, err := r.database.ExecContext(ctx, "UPDATE cpu_models SET cores = $1, architecture = $2, model = $3, vendor = $4 WHERE machine_id = $5", m.Cores, m.Architecture.String(), m.Model, m.Vendor, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) updateMemoryModel(ctx context.Context, mid machine.Identifier, m *machine.Memory) error {
	_, err := r.database.ExecContext(ctx, "UPDATE memory_models SET total = $1 WHERE machine_id = $2", m.Total, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) updateDiskModel(ctx context.Context, mid machine.Identifier, m *machine.Disk) error {
	_, err := r.database.ExecContext(ctx, "UPDATE disk_models SET device = $1, model = $2, vendor = $3, interface = $4, type = $5, serial = $6, sector_size = $7, size = $8, size_raw = $9 WHERE machine_id = $10", m.Device, m.Model, m.Vendor, m.Interface.String(), m.Type.String(), m.Serial, m.SectorSize, m.Size, m.SizeRaw, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) updateNetworkInterfaceModel(ctx context.Context, mid machine.Identifier, m *machine.NetworkInterface) error {
	var id int64
	err := r.database.GetContext(ctx, &id, "UPDATE network_interface_models SET name = $1, virtual = $2, mac_address = $3, vendor = $4, mtu = $5 WHERE machine_id = $6 RETURNING id", m.Name, m.Virtual, m.MacAddress, m.Vendor, m.Mtu, mid)
	if err != nil {
		return err
	}
	for _, a := range m.IpAddresses {
		err = r.updateAddressModel(ctx, id, a)
	}
	return nil
}

func (r *MachineRepo) updateAddressModel(ctx context.Context, nid int64, m *machine.IpAddress) error {
	_, err := r.database.ExecContext(ctx, "UPDATE address_models SET address = $1, version = $2 WHERE nic_id = $3", m.Address, m.Version.String(), nid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) updateVolumeModel(ctx context.Context, mid machine.Identifier, m *machine.Volume) error {
	_, err := r.database.ExecContext(ctx, "UPDATE volume_models SET name = $1, mount_point = $2, total = $3, file_system = $4 WHERE machine_id = $5", m.Name, m.MountPoint, m.Total, m.FileSystem, mid)
	if err != nil {
		return err
	}
	return nil
}
