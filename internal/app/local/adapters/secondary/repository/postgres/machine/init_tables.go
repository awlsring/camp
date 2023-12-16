package machine_repository

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (r *Repository) initTables() error {
	createMachinesTable := `
		CREATE TABLE IF NOT EXISTS machines (
			identifier VARCHAR(64) PRIMARY KEY,
			endpoint VARCHAR(255) NOT NULL,
			key VARCHAR(255) NOT NULL,
			class VARCHAR(255) NOT NULL,
			last_heartbeat TIMESTAMP NOT NULL,
			registered_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);`

	createStateTable := `CREATE TABLE IF NOT EXISTS power_state (
			id SERIAL PRIMARY KEY,
			state VARCHAR(255) NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);`

	createPowerCapabilityTable := `
		CREATE TABLE IF NOT EXISTS power_capabilites (
			id SERIAL PRIMARY KEY,
			reboot_enabled BOOLEAN NOT NULL,
			power_off_enabled BOOLEAN NOT NULL,
			wake_on_lan_enabled BOOLEAN NOT NULL,
			wake_on_lan_mac VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);`

	createHostTables := `
		CREATE TABLE IF NOT EXISTS host_models (
			id SERIAL PRIMARY KEY,
			os_family VARCHAR(255),
			kernel VARCHAR(255),
			os_name VARCHAR(255),
			os_version VARCHAR(255),
			os_platform VARCHAR(255),
			hostname VARCHAR(255),
			host_id VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createCpusTable := `
		CREATE TABLE IF NOT EXISTS cpu_models (
			id SERIAL PRIMARY KEY,
			total_cores INT NOT NULL,
			total_threads INT NOT NULL,
			architecture VARCHAR(255) NOT NULL,
			model VARCHAR(255),
			vendor VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createProcessorTable := `
		CREATE TABLE IF NOT EXISTS processor_models (
			id SERIAL PRIMARY KEY,
			identifier INT NOT NULL,
			core_count INT NOT NULL,
			thread_count INT NOT NULL,
			model VARCHAR(255),
			vendor VARCHAR(255),
			cpu_id INT NOT NULL,
			FOREIGN KEY (cpu_id) REFERENCES cpu_models(id)
		);
	`

	createCoreTable := `
		CREATE TABLE IF NOT EXISTS core_models (
			id SERIAL PRIMARY KEY,
			identifier INT NOT NULL,
			threads INT NOT NULL,
			processor_id INT NOT NULL,
			FOREIGN KEY (processor_id) REFERENCES processor_models(id)
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
			name VARCHAR(255) NOT NULL,
			size BIGINT NOT NULL,
			type VARCHAR(255) NOT NULL,
			storage_controller VARCHAR(255) NOT NULL,
			removable BOOLEAN NOT NULL,
			serial VARCHAR(255),
			model VARCHAR(255),
			vendor VARCHAR(255),
			wwn VARCHAR(255),
			machine_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (machine_id) REFERENCES machines(identifier)
		);
	`

	createPartitionsTable := `
		CREATE TABLE IF NOT EXISTS disk_models (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			size BIGINT NOT NULL,
			label VARCHAR(255),
			readonly BOOLEAN NOT NULL,
			type VARCHAR(255),
			fileSystem VARCHAR(255),
			uuid VARCHAR(255),
			mountPoint VARCHAR(255),
			disk_id VARCHAR(64) NOT NULL,
			FOREIGN KEY (disk_id) REFERENCES machines(id)
		);
	`

	createNetworkInterfacesTable := `
		CREATE TABLE IF NOT EXISTS network_interface_models (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			virtual BOOLEAN NOT NULL,
			mac_address VARCHAR(255),
			vendor VARCHAR(255),
			speed VARCHAR(255),
			duplex VARCHAR(255),
			pci_address VARCHAR(255),
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
		createStateTable,
		createPowerCapabilityTable,
		createHostTables,
		createCpusTable,
		createProcessorTable,
		createCoreTable,
		createMemoryTable,
		createDisksTable,
		createPartitionsTable,
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
