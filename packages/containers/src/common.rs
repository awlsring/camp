use std::collections::HashMap;
use chrono::{DateTime, Utc};

#[derive(Clone, Debug)]
pub enum ContainerState {
    Running,
    Created,
    Restarting,
    Removing,
    Paused,
    Stopped,
    Dead,
    Empty,
    Unknown(String),
}

#[derive(Clone, Debug)]
pub enum ContainerProtocol {
    TCP,
    UDP,
    SCTP,
    Unknown(String),
}

#[derive(Clone, Debug)]
pub enum ContainerType {
    Docker,
    Podman,
}

#[derive(Clone, Debug)]
pub struct Port {
    pub(crate) host_addresses: Vec<String>,
    pub(crate) container: u16,
    pub(crate) host: Option<u16>,
    pub(crate) protocol: ContainerProtocol,
}

impl Port {
    pub fn new(host_addresses: Vec<String>, container: u16, host: Option<u16>, protocol: ContainerProtocol) -> Port {
        Port {
            host_addresses,
            container,
            host,
            protocol,
        }
    }

    pub fn host_addresses(&self) -> &Vec<String> {
        &self.host_addresses
    }

    pub fn container(&self) -> u16 {
        self.container
    }

    pub fn host(&self) -> Option<u16> {
        self.host
    }

    pub fn protocol(&self) -> &ContainerProtocol {
        &self.protocol
    }
}

#[derive(Clone, Debug)]
pub struct Volume {
    pub(crate) source: Option<String>,
    pub(crate) destination: Option<String>,
    pub(crate) mode: Option<String>,
}

impl Volume {
    pub fn new(source: Option<String>, destination: Option<String>, mode: Option<String>) -> Volume {
        Volume {
            source,
            destination,
            mode,
        }
    }

    pub fn source(&self) -> Option<&str> {
        match &self.source {
            Some(source) => Some(source),
            None => None,
        }
    }

    pub fn destination(&self) -> Option<&str> {
        match &self.destination {
            Some(destination) => Some(destination),
            None => None,
        }
    }

    pub fn mode(&self) -> Option<&str> {
        match &self.mode {
            Some(mode) => Some(mode),
            None => None,
        }
    }
}

#[derive(Clone, Debug)]
pub struct Network {
    pub(crate) name: String,
    pub(crate) network_id: String,
    pub(crate) endpoint_id: String,
}

impl Network {
    pub fn new(name: String, network_id: String, endpoint_id: String) -> Network {
        Network {
            name,
            network_id,
            endpoint_id,
        }
    }

    pub fn name(&self) -> &str {
        &self.name
    }

    pub fn network_id(&self) -> &str {
        &self.network_id
    }

    pub fn endpoint_id(&self) -> &str {
        &self.endpoint_id
    }
}


#[derive(Clone, Debug)]
pub struct Container {
    pub(crate) id: String,
    pub(crate) name: String,
    pub(crate) image: String,
    pub(crate) created: DateTime<Utc>,
    pub(crate) started: Option<DateTime<Utc>>,
    pub(crate) finished: Option<DateTime<Utc>>,
    pub(crate) environment: Option<HashMap<String, String>>,
    pub(crate) command: Option<String>,
    pub(crate) state: ContainerState,
    pub(crate) ports: Vec<Port>,
    pub(crate) volumes: Vec<Volume>,
    pub(crate) networks: Vec<Network>,
    pub(crate) labels: Option<HashMap<String, String>>,
    pub(crate) statistics: Option<ContainerStatistics>,
    pub(crate) type_: ContainerType,
}

impl Container {
    pub fn id (&self) -> &str {
        &self.id
    }

    pub fn name (&self) -> &str {
        &self.name
    }

    pub fn image (&self) -> &str {
        &self.image
    }

    pub fn created (&self) -> &DateTime<Utc> {
        &self.created
    }

    pub fn started (&self) -> Option<&DateTime<Utc>> {
        match &self.started {
            Some(started) => Some(started),
            None => None,
        }
    }

    pub fn finished (&self) -> Option<&DateTime<Utc>> {
        match &self.finished {
            Some(finished) => Some(finished),
            None => None,
        }
    }

    pub fn environment (&self) -> Option<&HashMap<String, String>> {
        match &self.environment {
            Some(environment) => Some(environment),
            None => None,
        }
    }

    pub fn command (&self) -> Option<&str> {
        match &self.command {
            Some(command) => Some(command),
            None => None,
        }
    }

    pub fn state (&self) -> &ContainerState {
        &self.state
    }

    pub fn ports (&self) -> &Vec<Port> {
        &self.ports
    }

    pub fn volumes (&self) -> &Vec<Volume> {
        &self.volumes
    }

    pub fn networks (&self) -> &Vec<Network> {
        &self.networks
    }

    pub fn labels (&self) -> Option<&HashMap<String, String>> {
        match &self.labels {
            Some(labels) => Some(labels),
            None => None,
        }
    }

    pub fn statistics (&self) -> Option<&ContainerStatistics> {
        match &self.statistics {
            Some(statistics) => Some(statistics),
            None => None,
        }
    }

    pub fn type_ (&self) -> &ContainerType {
        &self.type_
    }

}

#[derive(Clone, Debug)]
pub struct ContainerStatistics {
    pub(crate) cpu_utilization: f64,
    pub(crate) memory_utilization: f64,
    pub(crate) memory_usage: u64,
    pub(crate) memory_limit: u64,
    pub(crate) network_rx_bytes: u64,
    pub(crate) network_tx_bytes: u64,
    pub(crate) block_read_bytes: u64,
    pub(crate) block_write_bytes: u64,
}

impl ContainerStatistics {
    pub fn cpu_utilization (&self) -> f64 {
        self.cpu_utilization
    }

    pub fn memory_utilization (&self) -> f64 {
        self.memory_utilization
    }

    pub fn memory_usage (&self) -> u64 {
        self.memory_usage
    }

    pub fn memory_limit (&self) -> u64 {
        self.memory_limit
    }

    pub fn network_rx_bytes (&self) -> u64 {
        self.network_rx_bytes
    }

    pub fn network_tx_bytes (&self) -> u64 {
        self.network_tx_bytes
    }

    pub fn block_read_bytes (&self) -> u64 {
        self.block_read_bytes
    }

    pub fn block_write_bytes (&self) -> u64 {
        self.block_write_bytes
    }
}

#[derive(Clone, Debug)]
pub struct ContainerLogLine {
    pub(crate) timestamp: DateTime<Utc>,
    pub(crate) line: String,
}

impl ContainerLogLine {
    pub fn timestamp (&self) -> &DateTime<Utc> {
        &self.timestamp
    }

    pub fn line (&self) -> &str {
        &self.line
    }
}