package machine_repository

import "context"

func (r *MachineRepo) enrichMachineEntry(ctx context.Context, m *MachineSql) error {
	system, err := r.getMachineSystem(ctx, m.Identifier)
	if err != nil {
		return err
	}
	m.System = system

	cpu, err := r.getMachineCpu(ctx, m.Identifier)
	if err != nil {
		return err
	}
	m.Cpu = cpu

	memory, err := r.getMachineMemory(ctx, m.Identifier)
	if err != nil {
		return err
	}
	m.Memory = memory

	disks, err := r.getMachineDisks(ctx, m.Identifier)
	if err != nil {
		return err
	}
	m.Disks = disks

	networkInterfaces, err := r.getMachineNetworkInterfaces(ctx, m.Identifier)
	if err != nil {
		return err
	}
	m.NetworkInterfaces = networkInterfaces

	volumes, err := r.getMachineVolumes(ctx, m.Identifier)
	if err != nil {
		return err
	}
	m.Volumes = volumes

	for _, ni := range networkInterfaces {
		for _, a := range ni.Addresses {
			m.Addresses = append(m.Addresses, a)
		}
	}

	return nil
}

func (r *MachineRepo) getMachineSystem(ctx context.Context, id string) (*SystemModelSql, error) {
	var m SystemModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM system_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MachineRepo) getMachineCpu(ctx context.Context, id string) (*CpuModelSql, error) {
	var m CpuModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM cpu_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MachineRepo) getMachineMemory(ctx context.Context, id string) (*MemoryModelSql, error) {
	var m MemoryModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM memory_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MachineRepo) getMachineDisks(ctx context.Context, id string) ([]*DiskModelSql, error) {
	var m []*DiskModelSql
	err := r.database.SelectContext(ctx, &m, "SELECT * FROM disk_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MachineRepo) getNicsAddresses(ctx context.Context, nicId int64) ([]*IpAddressModelSql, error) {
	var m []*IpAddressModelSql
	err := r.database.SelectContext(ctx, &m, "SELECT * FROM address_models WHERE nic_id = $1", nicId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MachineRepo) getMachineNetworkInterfaces(ctx context.Context, id string) ([]*NetworkInterfaceModelSql, error) {
	var m []*NetworkInterfaceModelSql
	err := r.database.SelectContext(ctx, &m, "SELECT * FROM network_interface_models WHERE machine_id = $1", id)
	for _, ni := range m {
		addresses, err := r.getNicsAddresses(ctx, ni.Id)
		if err != nil {
			return nil, err
		}
		ni.Addresses = addresses
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MachineRepo) getMachineVolumes(ctx context.Context, id string) ([]*VolumeModelSql, error) {
	var m []*VolumeModelSql
	err := r.database.SelectContext(ctx, &m, "SELECT * FROM volume_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return m, nil
}
