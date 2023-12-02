package machine_repository

import (
	"context"
	"strings"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/tag"
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"

	_ "github.com/lib/pq"

	"github.com/rs/zerolog/log"
)

var _ repository.Machine = &MachineRepo{}

type MachineRepo struct {
	database database.Database
	tagRepo  repository.Tag
}

func New(db database.Database, tagRepo repository.Tag) (repository.Machine, error) {
	r := &MachineRepo{
		database: db,
		tagRepo:  tagRepo,
	}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *MachineRepo) init() error {
	return r.initTables()
}

func (r *MachineRepo) initTables() error {
	createMachinesTable := `
		CREATE TABLE IF NOT EXISTS machines (
			identifier VARCHAR(64) PRIMARY KEY,
			class VARCHAR(255) NOT NULL,
			last_heartbeat TIMESTAMP NOT NULL,
			registered_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			status VARCHAR(255) NOT NULL
		);`

	createSystemsTables := `
		CREATE TABLE IF NOT EXISTS system_models (
			id SERIAL PRIMARY KEY,
			family VARCHAR(255),
			kernel_version VARCHAR(255),
			os VARCHAR(255),
			os_version VARCHAR(255),
			os_pretty VARCHAR(255),
			hostname VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createCpusTable := `
		CREATE TABLE IF NOT EXISTS cpu_models (
			id SERIAL PRIMARY KEY,
			cores INT NOT NULL,
			architecture VARCHAR(255) NOT NULL,
			model VARCHAR(255),
			vendor VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createMemoryTable := `
		CREATE TABLE IF NOT EXISTS memory_models (
			id SERIAL PRIMARY KEY,
			total BIGINT NOT NULL,
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createDisksTable := `
		CREATE TABLE IF NOT EXISTS disk_models (
			id SERIAL PRIMARY KEY,
			device VARCHAR(255) NOT NULL,
			model VARCHAR(255),
			vendor VARCHAR(255),
			interface VARCHAR(255) NOT NULL,
			type VARCHAR(255) NOT NULL,
			serial VARCHAR(255),
			sector_size INT NOT NULL,
			size BIGINT NOT NULL,
			size_raw BIGINT,
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createNetworkInterfacesTable := `
		CREATE TABLE IF NOT EXISTS network_interface_models (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			virtual BOOLEAN NOT NULL,
			mac_address VARCHAR(255),
			vendor VARCHAR(255),
			mtu BIGINT,
			speed BIGINT,
			duplex VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createVolumesTable := `
		CREATE TABLE IF NOT EXISTS volume_models (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			mount_point VARCHAR(255) NOT NULL,
			total BIGINT NOT NULL,
			file_system VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createAddressesTable := `
		CREATE TABLE IF NOT EXISTS address_models (
			id SERIAL PRIMARY KEY,
			address VARCHAR(255) NOT NULL,
			version VARCHAR(255) NOT NULL,
			nic_id INT NOT NULL,
			FOREIGN KEY (nic_id) REFERENCES network_interface_models(id)
		);
	`

	tableQueries := []string{
		createMachinesTable,
		createSystemsTables,
		createCpusTable,
		createMemoryTable,
		createDisksTable,
		createNetworkInterfacesTable,
		createVolumesTable,
		createAddressesTable,
	}

	ctx := context.Background()
	for _, query := range tableQueries {
		_, err := r.database.ExecContext(ctx, query)
		if err != nil {
			log.Error().Err(err).Msgf("Error creating table with query %s", query)
			return err
		}
	}
	return nil
}

func (r *MachineRepo) Close() error {
	return r.database.Close()
}

func (r *MachineRepo) Get(ctx context.Context, id machine.Identifier) (*machine.Machine, error) {
	var mdb MachineSql
	err := r.database.GetContext(ctx, &mdb, "SELECT * FROM machines WHERE identifier = $1", id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, repository.ErrMachineNotFound
		}
		return nil, err
	}
	err = r.enrichMachineEntry(ctx, &mdb)
	if err != nil {
		return nil, err
	}

	mod, err := mdb.ToModel()
	if err != nil {
		return nil, err
	}
	tags, err := r.tagRepo.ListForResource(ctx, id.String())
	if err != nil {
		return nil, err
	}
	mod.Tags = tags

	return mod, nil
}

func (r *MachineRepo) List(ctx context.Context, filters *repository.ListMachinesFilters) ([]*machine.Machine, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Invoke MachineRepo.GetMachines")

	log.Debug().Msg("Listing machines from database")
	var machinesModels []*MachineSql
	err := r.database.SelectContext(ctx, &machinesModels, "SELECT * FROM machines")
	if err != nil {
		log.Error().Err(err).Msg("Failed to list machines")
		return nil, err
	}

	log.Debug().Msgf("Found %d machines", len(machinesModels))
	models := []*machine.Machine{}
	for _, m := range machinesModels {
		log.Debug().Msgf("Enriching machine: %+v", m.Identifier)
		err = r.enrichMachineEntry(ctx, m)
		if err != nil {
			return nil, err
		}

		log.Debug().Msgf("Converting machine to domain: %+v", m.Identifier)
		mod, err := m.ToModel()
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert machine to domain")
			return nil, err
		}
		tags, err := r.tagRepo.ListForResource(ctx, m.Identifier)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list tags for machine")
			return nil, err
		}
		mod.Tags = tags

		log.Debug().Msgf("Appending machine to list: %+v", m.Identifier)
		models = append(models, mod)
	}

	log.Debug().Msgf("Returning %d machines", len(models))
	return models, nil
}

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

func (r *MachineRepo) Add(ctx context.Context, m *machine.Machine) error {
	err := r.createMachineEntry(ctx, m)
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

func (r *MachineRepo) AddTags(ctx context.Context, id machine.Identifier, tags []*tag.Tag) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("Adding tags to machine %s", id.String())

	for _, t := range tags {
		err := r.tagRepo.AddToResource(ctx, t, id.String(), tag.ResourceTypeMachine)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func (r *MachineRepo) UpdateHeartbeat(ctx context.Context, id machine.Identifier) error {
	_, err := r.database.ExecContext(ctx, "UPDATE machines SET last_heartbeat = NOW() WHERE identifier = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r MachineRepo) UpdateStatus(ctx context.Context, id machine.Identifier, status machine.MachineStatus) error {
	_, err := r.database.ExecContext(ctx, "UPDATE machines SET status = $1 WHERE identifier = $2", status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) Delete(ctx context.Context, id machine.Identifier) error {
	_, err := r.database.ExecContext(ctx, "DELETE FROM machines WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) createMachineEntry(ctx context.Context, m *machine.Machine) error {
	now := time.Now()
	_, err := r.database.ExecContext(ctx, "INSERT INTO machines (identifier, class, last_heartbeat, registered_at, updated_at, status) VALUES ($1, $2, $3, $4, $5, $6)", m.Identifier, m.Class, now, now, now, machine.MachineStatusRunning.String())
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
	_, err := r.database.ExecContext(ctx, "INSERT INTO cpu_models (cores, architecture, model, vendor, machine_id) VALUES ($1, $2, $3, $4, $5)", m.Cores, m.Architecture, m.Model, m.Vendor, mid)
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
	_, err := r.database.ExecContext(ctx, "INSERT INTO disk_models (device, model, vendor, interface, type, serial, sector_size, size, size_raw, machine_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", m.Device, m.Model, m.Vendor, m.Interface, m.Type, m.Serial, m.SectorSize, m.Size, m.SizeRaw, mid)
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
	_, err := r.database.ExecContext(ctx, "INSERT INTO address_models (address, version, nic_id) VALUES ($1, $2, $3)", m.Address, m.Version, nid)
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

func (r *MachineRepo) updateSystemModel(ctx context.Context, mid machine.Identifier, m *machine.System) error {
	_, err := r.database.ExecContext(ctx, "UPDATE system_models SET family = $1, kernel_version = $2, os = $3, os_version = $4, os_pretty = $5, hostname = $6 WHERE machine_id = $7", m.Os.Family, m.Os.Kernel, m.Os, m.Os.Version, m.Os.PrettyName, m.Hostname, mid)
	if err != nil {
		return err
	}
	return nil
}

func (r *MachineRepo) updateCpuModel(ctx context.Context, mid machine.Identifier, m *machine.Cpu) error {
	_, err := r.database.ExecContext(ctx, "UPDATE cpu_models SET cores = $1, architecture = $2, model = $3, vendor = $4 WHERE machine_id = $5", m.Cores, m.Architecture, m.Model, m.Vendor, mid)
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
	_, err := r.database.ExecContext(ctx, "UPDATE disk_models SET device = $1, model = $2, vendor = $3, interface = $4, type = $5, serial = $6, sector_size = $7, size = $8, size_raw = $9 WHERE machine_id = $10", m.Device, m.Model, m.Vendor, m.Interface, m.Type, m.Serial, m.SectorSize, m.Size, m.SizeRaw, mid)
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
	_, err := r.database.ExecContext(ctx, "UPDATE address_models SET address = $1, version = $2 WHERE nic_id = $3", m.Address, m.Version, nid)
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
