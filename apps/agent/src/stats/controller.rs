use std::collections::HashMap;

use containers::{Container, Containers};
use hw_info::{load_disks, Disk};
use sysinfo::System as Sys;
use sysinfo::SystemExt;

use super::cpu::Cpu;
use super::disk::Storage;
use super::memory::Memory;
use super::network::AddressKind;
use super::network::Network;
use super::system::System;
use super::system_description::CpuDescription;
use super::system_description::DiskDescription;
use super::system_description::HostDescription;
use super::system_description::IpAddressDescription;
use super::system_description::MemoryDescription;
use super::system_description::NetworkDescription;
use super::system_description::SystemDescription;
use super::system_description::VolumeDescription;

pub struct SystemController {
    system_controller: Sys,
    container_controller: Option<Containers>,
    system: System,
    memory: Memory,
    cpu: Cpu,
    network: Network,
    storage: Storage,
    disks: HashMap<String, Disk>,
    containers: HashMap<String, Container>,
}

impl SystemController {
    pub fn new() -> SystemController {
        let mut sys = Sys::new_all();
        sys.refresh_all();
        let container_controller = Containers::new();

        let system = System::new(&sys);
        let memory = Memory::new(&sys);
        let cpu = Cpu::new(&sys);
        let network = Network::new(&sys);
        let storage = Storage::new(&sys);
        let mut disks = HashMap::<String, Disk>::new();
        for disk in load_disks() {
            disks.insert(disk.get_device().to_string(), disk);
        }
        let containers = HashMap::<String, Container>::new();

        SystemController {
            system_controller: sys,
            container_controller,
            system,
            memory,
            cpu,
            network,
            storage,
            disks,
            containers,
        }
    }

    pub fn system(&self) -> &System {
        &self.system
    }

    pub fn cpu(&self) -> &Cpu {
        &self.cpu
    }

    pub fn memory(&self) -> &Memory {
        &self.memory
    }

    pub fn network(&self) -> &Network {
        &self.network
    }

    pub fn storage(&self) -> &Storage {
        &self.storage
    }

    pub fn disks(&self) -> &HashMap<String, Disk> {
        &self.disks
    }

    pub(crate) fn get_system_description(&self) -> SystemDescription {
        let mut disks = vec![];
        for disk in self.disks.values() {
            disks.push(DiskDescription {
                device: disk.get_device().to_string(),
                size: disk.get_size_actual().to_owned(),
                model: disk.get_model().to_string(),
                vendor: disk.get_vendor().to_string(),
                interface_type: match disk.get_interface() {
                    hw_info::DiskInterface::SATA => "SATA".to_string(),
                    hw_info::DiskInterface::SCSI => "PATA".to_string(),
                    hw_info::DiskInterface::PCI_E => "SCSI".to_string(),
                    _ => "Unknown".to_string(),
                },
                serial_number: disk.get_serial().to_string(),
                size_raw: disk.get_size_raw().to_owned(),
                sector_size: disk.get_sector_size().to_owned(),
                disk_type: match disk.get_kind() {
                    hw_info::DiskKind::HDD => "HDD".to_string(),
                    hw_info::DiskKind::SSD => "SSD".to_string(),
                    hw_info::DiskKind::NVME => "SSD".to_string(),
                    _ => "Unknown".to_string(),
                },
            });
        }

        let mut network = vec![];
        for net in self.network.network_interfaces() {
            let mut ip_addresses = vec![];
            for ip in net.addresses() {
                ip_addresses.push(IpAddressDescription {
                    address: ip.address().to_string(),
                    version: match ip.version() {
                        AddressKind::V4 => "v4".to_string(),
                        AddressKind::V6 => "v6".to_string(),
                        _ => "Unknown".to_string(),
                    },
                });
            }
            network.push(NetworkDescription {
                name: net.name().to_string(),
                mac_address: net.mac().to_owned(),
                is_virtual: *net.is_virtual(),
                vendor: net.vendor().to_owned(),
                mtu: net.mtu().to_owned(),
                duplex: net.duplex().to_owned(),
                speed: net.speed().to_owned(),
                ip_addresses,
            });
        }

        let mut volumes = vec![];
        for vol in self.storage.volumes() {
            volumes.push(VolumeDescription {
                name: vol.name().to_string(),
                mount_point: vol.mount_point().to_string(),
                file_system: vol.file_system().to_string(),
                total: vol.total_space().to_owned() as i64,
            });
        }

        SystemDescription {
            machine_id: self.system.machine_id().to_string(),
            class: "BareMetal".to_string(), //TODO: determine from ctl, allow override in config
            host: HostDescription {
                family: self.system.family().to_string(),
                kernel_version: self.system.kernel_version().to_string(),
                os_pretty: self.system.os_pretty().to_string(),
                os_version: self.system.os_version().to_string(),
                os: self.system.os().to_string(),
                hostname: self.system.hostname().to_string(),
            },
            cpu: CpuDescription {
                cores: self.cpu.core_count() as i32,
                architecture: self.cpu.architecture().to_string(),
                vendor: self.cpu.vendor().to_string(),
                model: self.cpu.brand().to_string(),
            },
            memory: MemoryDescription {
                total: self.memory.memory().total().to_owned() as i64,
            },
            disks,
            network,
            volumes,
        }
    }

    pub async fn refresh(&mut self) -> bool {
        self.refresh_system().await;
        self.refresh_memory().await;
        self.refresh_cpu().await;
        self.refresh_network().await;
        self.refresh_storage().await;
        self.refresh_containers().await;
        false
    }

    async fn refresh_system(&mut self) {
        self.system_controller.refresh_system();
        self.system.update_up_time(&self.system_controller)
    }

    async fn refresh_memory(&mut self) {
        self.system_controller.refresh_memory();
        self.memory.update(&self.system_controller);
    }

    async fn refresh_cpu(&mut self) {
        self.system_controller.refresh_cpu();
        self.cpu.update(&self.system_controller);
    }

    async fn refresh_network(&mut self) {
        self.system_controller.refresh_networks_list();
        self.system_controller.refresh_networks();
        self.network.update(&self.system_controller);
    }

    async fn refresh_storage(&mut self) {
        self.system_controller.refresh_disks_list();
        self.system_controller.refresh_disks();
        self.storage.update(&self.system_controller);
    }

    async fn refresh_containers(&mut self) {}
}
