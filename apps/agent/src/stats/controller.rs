use std::collections::HashMap;

use containers::{Container, Containers};
use hw_info::{load_disks, Disk};
use sysinfo::System as Sys;
use sysinfo::SystemExt;

use super::cpu::Cpu;
use super::disk::Storage;
use super::memory::Memory;
use super::network::Network;
use super::system::System;

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
            container_controller: container_controller,
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

    pub fn containers(&self) -> &HashMap<String, Container> {
        &self.containers
    }

    pub async fn refresh(&mut self) {
        self.refresh_system().await;
        self.refresh_memory().await;
        self.refresh_cpu().await;
        self.refresh_network().await;
        self.refresh_storage().await;
        self.refresh_containers().await;
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

    async fn refresh_containers(&mut self) {
        // if let Some(ref mut container_controller) = self.container_controller {
        //     let containers = container_controller.list_containers().await;
        //     for container in containers {
        //         self.containers.insert(container.id().to_string(), container);
        //     }
        // }
    }
}
