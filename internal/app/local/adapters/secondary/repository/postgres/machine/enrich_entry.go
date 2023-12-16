package machine_repository

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
)

var (
	ErrPowerRelationNotFound = "pq: relation \"power_capabilities\" does not exist"
)

func (r *Repository) enrichMachineEntry(ctx context.Context, m *MachineSql) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Enriching machine %s", m.Identifier)

	log.Debug().Msgf("Getting machine %s power state", m.Identifier)
	powerState, err := r.getMachineState(ctx, m.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get machine %s power state", m.Identifier)
		return err
	}
	m.Status = powerState

	log.Debug().Msgf("Getting machine %s power capabilities", m.Identifier)
	powerCapabilities, err := r.getMachinePowerCapabilities(ctx, m.Identifier)
	if err != nil {
		if err.Error() == ErrPowerRelationNotFound {
			log.Warn().Msgf("Machine %s has no power capabilities, setting defaults", m.Identifier)
			powerCapabilities = &PowerCapabilityModelSql{
				Id:        0,
				MachineId: m.Identifier,
				Reboot:    true,
				PowerOff:  true,
				WakeOnLan: true,
			}
		} else {
			log.Error().Err(err).Msgf("Failed to get machine %s power capabilities", m.Identifier)
			return err
		}
	}
	log.Debug().Msgf("Machine %s power capabilities retrieved", m.Identifier)
	m.PowerCapabilities = powerCapabilities

	log.Debug().Msgf("Getting machine %s system", m.Identifier)
	host, err := r.getMachineHost(ctx, m.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get machine %s system", m.Identifier)
		return err
	}
	log.Debug().Msgf("Machine %s system retrieved", m.Identifier)
	m.Host = host

	log.Debug().Msgf("Getting machine %s cpu", m.Identifier)
	cpu, err := r.getMachineCpu(ctx, m.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get machine %s cpu", m.Identifier)
		return err
	}
	log.Debug().Msgf("Machine %s cpu retrieved", m.Identifier)
	m.Cpu = cpu

	log.Debug().Msgf("Getting machine %s memory", m.Identifier)
	memory, err := r.getMachineMemory(ctx, m.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get machine %s memory", m.Identifier)
		return err
	}
	log.Debug().Msgf("Machine %s memory retrieved", m.Identifier)
	m.Memory = memory

	log.Debug().Msgf("Getting machine %s disks", m.Identifier)
	disks, err := r.getMachineDisks(ctx, m.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get machine %s disks", m.Identifier)
		return err
	}
	log.Debug().Msgf("Machine %s disks retrieved", m.Identifier)
	m.Disks = disks

	log.Debug().Msgf("Getting machine %s network interfaces", m.Identifier)
	networkInterfaces, err := r.getMachineNetworkInterfaces(ctx, m.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get machine %s network interfaces", m.Identifier)
		return err
	}
	log.Debug().Msgf("Machine %s network interfaces retrieved", m.Identifier)
	m.NetworkInterfaces = networkInterfaces

	log.Debug().Msgf("Getting machine %s volumes", m.Identifier)
	volumes, err := r.getMachineVolumes(ctx, m.Identifier)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get machine %s volumes", m.Identifier)
		return err
	}
	log.Debug().Msgf("Machine %s volumes retrieved", m.Identifier)
	m.Volumes = volumes

	log.Debug().Msg("Appending addresses to nics")
	for _, ni := range networkInterfaces {
		for _, a := range ni.Addresses {
			m.Addresses = append(m.Addresses, a)
		}
	}

	log.Debug().Msgf("Machine %s enriched", m.Identifier)
	return nil
}

func (r *Repository) getMachineState(ctx context.Context, id string) (*StatusModelSql, error) {
	var m StatusModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM power_state WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) getMachinePowerCapabilities(ctx context.Context, id string) (*PowerCapabilityModelSql, error) {
	var m PowerCapabilityModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM power_capabilities WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) getMachineHost(ctx context.Context, id string) (*HostModelSql, error) {
	var m HostModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM host_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) getMachineCpu(ctx context.Context, id string) (*CpuModelSql, error) {
	var m CpuModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM cpu_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) getMachineMemory(ctx context.Context, id string) (*MemoryModelSql, error) {
	var m MemoryModelSql
	err := r.database.GetContext(ctx, &m, "SELECT * FROM memory_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) getMachineDisks(ctx context.Context, id string) ([]*DiskModelSql, error) {
	var m []*DiskModelSql
	err := r.database.SelectContext(ctx, &m, "SELECT * FROM disk_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *Repository) getNicsAddresses(ctx context.Context, nicId int64) ([]*IpAddressModelSql, error) {
	var m []*IpAddressModelSql
	err := r.database.SelectContext(ctx, &m, "SELECT * FROM address_models WHERE nic_id = $1", nicId)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *Repository) getMachineNetworkInterfaces(ctx context.Context, id string) ([]*NetworkInterfaceModelSql, error) {
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

func (r *Repository) getMachineVolumes(ctx context.Context, id string) ([]*VolumeModelSql, error) {
	var m []*VolumeModelSql
	err := r.database.SelectContext(ctx, &m, "SELECT * FROM volume_models WHERE machine_id = $1", id)
	if err != nil {
		return nil, err
	}
	return m, nil
}
