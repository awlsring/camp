use sysinfo::{System, SystemExt, CpuExt};
use std::{env::consts::ARCH, collections::HashMap};

use super::util;

pub struct Core {
    name: String,
    frequency: u64,
    usage: f32,
}

impl Core {
    pub fn name(&self) -> &String {
        &self.name
    }

    pub fn frequency(&self) -> &u64 {
        &self.frequency
    }

    pub fn usage(&self) -> &f32 {
        &self.usage
    }
}

pub struct Cpu {
    cores: HashMap<String, Core>,
    core_count: usize,
    architecture: String,
    vendor: String,
    brand: String,
}

impl Cpu {
    pub fn new(system: &System) -> Cpu {
        let arch = String::from(ARCH);
        
        

        let cpu = &system.cpus()[0];

        let vendor = cpu.vendor_id().to_string();
        let brand = cpu.brand().to_string();

        let mut cores = HashMap::new();

        for cpu in system.cpus() {
            let core = Core {
                name: cpu.name().to_string(),
                frequency: cpu.frequency(),
                usage: cpu.cpu_usage(),
            };
            cores.insert(cpu.name().to_string(), core);
        }

        let core_count = util::handle_optional_usize(system.physical_core_count());

        Cpu {
            core_count,
            cores,
            architecture: arch,
            vendor,
            brand,
        }
    }

    pub fn update(&mut self, system: &System) {
        for cpu in system.cpus() {
            let core = Core {
                name: cpu.name().to_string(),
                frequency: cpu.frequency(),
                usage: cpu.cpu_usage(),
            };
            self.cores.insert(cpu.name().to_string(), core);
        }
    }

    pub fn cores(&self) -> Vec<&Core> {
        let mut cores = Vec::new();

        for core in self.cores.values() {
            cores.push(core);
        }

        cores
    }

    pub fn core_count(&self) -> usize {
        self.core_count
    }

    pub fn architecture(&self) -> &String {
        &self.architecture
    }

    pub fn vendor(&self) -> &String {
        &self.vendor
    }

    pub fn brand(&self) -> &String {
        &self.brand
    }
}



