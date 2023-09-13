use std::collections::{HashMap, hash_map::RandomState};

use bollard::{Docker, service::{ContainerSummary, PortTypeEnum, ContainerInspectResponse, ContainerStateStatusEnum, ContainerConfig, NetworkSettings, MountPoint, EndpointSettings, PortBinding}, secret::ContainerState as BollardContainerState, container::{Stats, CPUStats, MemoryStats, BlkioStats, NetworkStats, LogOutput}};
use chrono::{DateTime, Utc};

use crate::common::{Port, ContainerProtocol, Volume, Network, ContainerStatistics, ContainerLogLine};

use super::common::{Container, ContainerState};

fn timestamp_to_datetime(ts: Option<String>) -> DateTime<Utc> {
    match ts {
        Some(ts) => DateTime::<Utc>::from(DateTime::parse_from_rfc3339(&ts).unwrap_or_default()),
        None => DateTime::<Utc>::default(),
    }
}

fn read_container_image(config: Option<ContainerConfig>) -> String {
    match config {
        Some(config) => config.image.unwrap_or("unknown".to_string()),
        None => "unknown".to_string(),
    }
}

fn handle_optional_string(name: Option<String>) -> String {
    match name {
        Some(name) => name,
        None => String::from(""),
    }
}

fn handle_optional_date(name: Option<String>) -> Option<DateTime<Utc>> {
    match name {
        Some(name) => Some(DateTime::<Utc>::from(DateTime::parse_from_rfc3339(&name).unwrap_or_default())),
        None => None,
    }
}

fn parse_port_and_protocol(port_string: String) -> (u16, ContainerProtocol) {
    let mut parts = port_string.split("/");
    let port = parts.next().unwrap().parse::<u16>().unwrap();
    let protocol = match parts.next() {
        Some(protocol) => match protocol {
            "tcp" => ContainerProtocol::TCP,
            "udp" => ContainerProtocol::UDP,
            _ => ContainerProtocol::TCP,
        },
        None => ContainerProtocol::TCP,
    };
    (port, protocol)
}

pub fn docker_state_to_state(state: Option<BollardContainerState>) -> ContainerState {
    match state {
        Some(state) => { 
            match state.status {
                Some(status) => match status {
                    ContainerStateStatusEnum::RUNNING => ContainerState::Running,
                    ContainerStateStatusEnum::CREATED => ContainerState::Created,
                    ContainerStateStatusEnum::RESTARTING => ContainerState::Restarting,
                    ContainerStateStatusEnum::REMOVING => ContainerState::Removing,
                    ContainerStateStatusEnum::PAUSED => ContainerState::Paused,
                    ContainerStateStatusEnum::EXITED => ContainerState::Stopped,
                    ContainerStateStatusEnum::DEAD => ContainerState::Dead,
                    ContainerStateStatusEnum::EMPTY => ContainerState::Empty,
                    (_) => ContainerState::Unknown(String::from("unknown")),
                }
                None => ContainerState::Unknown(String::from("unknown")),
            }
        },
        None => ContainerState::Stopped,
    }
}

fn ports_from_network(network: Option<NetworkSettings>) -> Vec<Port> {
    let mut ports = Vec::new();
    match network {
        Some(network) => {
            match network.ports {
                Some(network_ports) => {
                    for (container_port_config, host_ports) in network_ports {
                        let port = Port::from_port_mappings(container_port_config, host_ports);
                        ports.push(port);
                    }
                },
                None => {},
            }
            ports
        },
        None => ports,
    }
}

fn mounts_to_volumes(mounts: Option<Vec<MountPoint>>) -> Vec<Volume> {
    let mut volumes = Vec::new();
    match mounts {
        Some(mounts) => {
            for mount in mounts {
                let source = match mount.source {
                    Some(source) => Some(source),
                    None => None,
                };
                let destination = match mount.destination {
                    Some(destination) => Some(destination),
                    None => None,
                };
                let mode = match mount.mode {
                    Some(mode) => Some(mode),
                    None => None,
                };
                let volume = Volume::new(source, destination, mode);
                volumes.push(volume);
            }
        },
        None => {},
    }
    volumes
}

fn networks_from_network_map(networks_settings: Option<NetworkSettings>) -> Vec<Network> {
    let mut result: Vec<Network> = Vec::new();
    match networks_settings {
        Some(networks) => {
            match networks.networks {
                Some(networks) => {
                    for (network_name, network) in networks {
                        let network = Network::from_docker_network(network_name, network);
                        result.push(network);
                    };
                },
                None => {},
            }
        },
        None => {},
    }
    result
}

impl Container {
    pub fn new_from_docker(container: ContainerInspectResponse, statistics: Option<ContainerStatistics>) -> Container {
        let id = handle_optional_string(container.id);
        let name = handle_optional_string(container.name);
        let image = read_container_image(container.config.clone()); // <- fix later
        let state = docker_state_to_state(container.state.clone()); // <- fix later
        let created = timestamp_to_datetime(container.created);
        let started = handle_optional_date(container.state.clone().unwrap_or_default().started_at); // <- fix later
        let finished = handle_optional_date(container.state.unwrap_or_default().finished_at);
        let ports = ports_from_network(container.network_settings.clone()); // <- fix later
        let volumes = mounts_to_volumes(container.mounts);
        let networks = networks_from_network_map(container.network_settings);
        let command = container.path.to_owned();
        let environment: Option<HashMap<String, String>> = match container.config.clone() { // <- fix later
            Some(config) => match config.env {
                Some(env) => {
                    let mut result: HashMap<String, String> = HashMap::new();
                    for e in env {
                        let mut parts = e.split("=");
                        let key = parts.next().unwrap().to_owned();
                        let value = parts.next().unwrap_or("").to_owned();
                        result.insert(key, value);
                    }
                    Some(result)
                },
                None => None,
            },
            None => None,
        };
        let labels: Option<HashMap<String, String>> = match container.config {
            Some(config) => match config.labels {
                Some(labels) => Some(labels),
                None => None,
            },
            None => None,
        };

        Container {
            id,
            name,
            image,
            state,
            created,
            started,
            finished,
            command,
            ports,
            volumes,
            networks,
            labels,
            environment,
            statistics,
            type_: super::common::ContainerType::Docker,
        }
    }
}

impl ContainerStatistics {
    pub fn new_from_docker(stats: Stats) -> ContainerStatistics {
        let cpu_utilization = get_cpu_utilization(stats.cpu_stats);
        let (memory_utilization, memory_usage, memory_limit) = get_mem_stats(stats.memory_stats);
        let (network_rx_bytes, network_tx_bytes) = get_network_stats(stats.networks);
        let (block_read_bytes, block_write_bytes) = get_block_stats(stats.blkio_stats);

        ContainerStatistics {
            cpu_utilization,
            memory_utilization,
            memory_usage,
            memory_limit,
            network_rx_bytes,
            network_tx_bytes,
            block_read_bytes,
            block_write_bytes,
        }
    }
}

fn get_cpu_utilization(cpu_stats: CPUStats) -> f64 {
    let total_usage = cpu_stats.cpu_usage.total_usage;
    let system_cpu_usage = cpu_stats.system_cpu_usage.unwrap_or_default();
    let online_cpus = cpu_stats.online_cpus.unwrap_or_default();

    let cpu_usage_percentage = (total_usage as f64 / system_cpu_usage as f64) * 100.0 * online_cpus as f64;
    cpu_usage_percentage
}

fn get_mem_stats(memory_stats: MemoryStats) -> (f64, u64, u64) {
    let mem_usage = memory_stats.usage.unwrap_or_default();
    let mem_limit = memory_stats.limit.unwrap_or_default();

    let mem_usage_percentage = (mem_usage as f64 / mem_limit as f64) * 100.0;
    
    (mem_usage_percentage, mem_usage, mem_limit)
}

fn get_block_stats(block_stats: BlkioStats) -> (u64, u64) {
    let io_service_bytes_recursive = block_stats.io_service_bytes_recursive.unwrap_or_default();
    let mut read_bytes = 0;
    let mut write_bytes = 0;
    for io in io_service_bytes_recursive {
        if io.op == "read" {
            read_bytes = io.value;
        } else if io.op == "write" {
            write_bytes = io.value;
        }
    }

    (read_bytes, write_bytes)
}

fn get_network_stats(network_stats: Option<HashMap<String, NetworkStats, RandomState>>) -> (u64, u64) {
    match network_stats {
        Some(network_stats) => {
            let mut rx_bytes = 0;
            let mut tx_bytes = 0;

            for (_, stats) in network_stats {
                rx_bytes += stats.rx_bytes;
                tx_bytes += stats.tx_bytes;
            }
        
            return (rx_bytes, tx_bytes)
        },
        None => return (0 as u64, 0 as u64)
    };
}

impl ContainerLogLine {
    pub fn new_from_docker(line: LogOutput) -> ContainerLogLine {
        let l = match line {
            LogOutput::StdOut { message } => message,
            LogOutput::StdErr { message } => message,
            LogOutput::Console { message } => message,
            LogOutput::StdIn { message } => message,
        };

        let line_string = String::from_utf8(l.to_vec()).unwrap_or_default();
        
        let (time, line) = ContainerLogLine::parse_docker_line(line_string);
        ContainerLogLine { timestamp: time, line: line }
    }

    fn parse_docker_line(line: String) -> (DateTime<Utc>, String) {
        let mut parts = line.splitn(2, " ");
        let timestamp = parts.next().unwrap_or_default();
        let line = parts.next().unwrap_or("");

        let time = DateTime::<Utc>::from(DateTime::parse_from_rfc3339(&timestamp).unwrap_or_default());

        (time, line.to_owned())
    }
}

impl Network {
    pub fn from_docker_network(name: String, endpoint: EndpointSettings) -> Network {
        let endpoint_id = endpoint.endpoint_id.unwrap_or("unknown".to_owned());
        let network_id = endpoint.network_id.unwrap_or("unknown".to_owned());

        Network {
            name,
            endpoint_id,
            network_id,
        }
    }
}

impl Port {
    pub fn from_port_mappings(container_port_string: String, port_bindings: Option<Vec<PortBinding>>) -> Port {
        let (container_port, protocol) = parse_port_and_protocol(container_port_string);
        let mut host_addresses = Vec::new();
        let mut host_port: Option<u16> = None;
        match port_bindings {
            Some(host_ports) => {
                for host_port_binding in host_ports {
                    match host_port_binding.host_ip {
                        Some(host_ip) => host_addresses.push(host_ip),
                        None => {},
                    };
                    match host_port_binding.host_port {
                        Some(hp) => {
                            let hp_int = hp.parse::<u16>().unwrap_or(0);
                            if hp_int > 0 {
                                host_port = Some(hp_int);
                            }
                        }
                        None => {},
                    };
                }
            },
            None => {},
        };
        Port {
            host_addresses,
            protocol,
            host: host_port,
            container: container_port,
        }
    }
}