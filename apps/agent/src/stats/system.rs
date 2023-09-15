use std::io::Read;

use sysinfo::System as Sys;
use sysinfo::SystemExt;

use super::util::handle_optional_string;

pub struct System {
    machine_id: String,
    family: String,
    kernel_version: String,
    os_pretty: String,
    os_version: String,
    os: String,
    hostname: String,
    boot_time: u64,
    up_time: u64,
}

fn get_machine_id() -> String {
    let mut machine_id = String::new();
    if let Ok(mut file) = std::fs::File::open("/etc/machine-id") {
        if file.read_to_string(&mut machine_id).is_ok() {
            machine_id = machine_id.trim().to_string();
        }
    }
    machine_id
}

impl System {
    pub fn new(system: &Sys) -> System {
        let family = handle_optional_string(system.name());
        let kernel_version = handle_optional_string(system.kernel_version());
        let os_pretty = handle_optional_string(system.long_os_version());
        let os_version = handle_optional_string(system.os_version());
        let os = handle_optional_string(Some(system.distribution_id()));
        let hostname = handle_optional_string(system.host_name());
        let boot_time = system.boot_time();
        let up_time = system.uptime();

        System {
            machine_id: get_machine_id(),
            family,
            kernel_version,
            os_pretty,
            os_version,
            os,
            hostname,
            boot_time,
            up_time,
        }
    }

    pub fn machine_id(&self) -> &String {
        &self.machine_id
    }

    pub fn family(&self) -> &String {
        &self.family
    }

    pub fn kernel_version(&self) -> &String {
        &self.kernel_version
    }

    pub fn os_pretty(&self) -> &String {
        &self.os_pretty
    }

    pub fn os_version(&self) -> &String {
        &self.os_version
    }

    pub fn os(&self) -> &String {
        &self.os
    }

    pub fn hostname(&self) -> &String {
        &self.hostname
    }

    pub fn boot_time(&self) -> &u64 {
        &self.boot_time
    }

    pub fn up_time(&self) -> &u64 {
        &self.up_time
    }

    pub fn update_up_time(&mut self, system: &Sys) {
        self.up_time = system.uptime();
    }
}
