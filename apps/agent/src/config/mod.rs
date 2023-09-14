use log::{debug, warn};
use serde::Deserialize;
use std::{env, fs};
use toml;

#[derive(Debug, Deserialize)]
pub struct Config {
    agent: AgentConfig,
    server: ServerConfig,
    reporting: ReportingConfig,
}

impl Default for Config {
    fn default() -> Self {
        Config {
            agent: AgentConfig {
                interval: 10000,
                cpu: None,
                mem: None,
                disk: None,
                network: None,
            },
            server: ServerConfig {
                port: 7032,
                allowed_keys: vec![String::from("toes")],
                no_auth_operations: vec![String::from("Health")],
            },
            reporting: ReportingConfig {
                api_key: String::from("a"),
                endpoint: String::from("127.0.0.1:8080"),
            },
        }
    }
}
impl Config {
    pub fn take_server(&mut self) -> ServerConfig {
        std::mem::replace(&mut self.server, Default::default())
    }

    pub fn take_agent(&mut self) -> AgentConfig {
        std::mem::replace(&mut self.agent, Default::default())
    }

    pub fn take_reporting(&mut self) -> ReportingConfig {
        std::mem::replace(&mut self.reporting, Default::default())
    }
}

#[derive(Clone, Default, Debug, Deserialize)]
pub struct ServerConfig {
    port: u16,
    #[serde(rename = "allowedKeys")]
    allowed_keys: Vec<String>,
    #[serde(rename = "noAuthOperations")]
    no_auth_operations: Vec<String>,
}

impl ServerConfig {
    pub fn get_server_port(&self) -> u16 {
        self.port
    }

    pub fn allowed_keys(&self) -> &Vec<String> {
        &self.allowed_keys
    }

    pub fn no_auth_operations(&self) -> &Vec<String> {
        &self.no_auth_operations
    }
}

#[derive(Clone, Default, Debug, Deserialize)]
pub struct AgentConfig {
    interval: u64,
    cpu: Option<u64>,
    mem: Option<u64>,
    disk: Option<u64>,
    network: Option<u64>,
}

impl AgentConfig {
    pub fn get_interval(&self) -> u64 {
        self.interval
    }
    pub fn get_cpu_interval(&self) -> Option<u64> {
        self.cpu
    }
    pub fn get_mem_interval(&self) -> Option<u64> {
        self.mem
    }
    pub fn get_disk_interval(&self) -> Option<u64> {
        self.disk
    }
    pub fn get_network_interval(&self) -> Option<u64> {
        self.network
    }
}

#[derive(Clone, Default, Debug, Deserialize)]
pub struct ReportingConfig {
    #[serde(rename = "apiKey")]
    api_key: String,
    endpoint: String,
}

impl ReportingConfig {
    pub fn get_api_key(&self) -> &str {
        &self.api_key
    }
    pub fn get_endpoint(&self) -> &str {
        &self.endpoint
    }
}

pub fn load_config() -> Config {
    let path = env::var("CONFIG_PATH").unwrap_or_else(|_| "config.toml".to_string());
    debug!("Loading config from: {}", path);
    let config = fs::read_to_string(path);
    match config {
        Ok(config) => {
            debug!("Loaded config from file");
            parse_config(config)
        }
        Err(_) => {
            warn!("Failed to load config from file, using default config");
            Config::default()
        }
    }
}

fn parse_config(config: String) -> Config {
    let config = toml::from_str(&config);
    match config {
        Ok(config) => {
            debug!("Parsed config");
            config
        }
        Err(_) => {
            warn!("Failed to parse config, using default config");
            Config::default()
        }
    }
}
