use sysinfo::{System, SystemExt};

pub struct MemoryObject {
    total: u64,
    used: u64,
    available: u64,
}

impl MemoryObject {
    pub fn total(&self) -> &u64 {
        &self.total
    }

    pub fn used(&self) -> &u64 {
        &self.used
    }

    pub fn available(&self) -> &u64 {
        &self.available
    }
}

pub struct Memory {
    memory: MemoryObject,
    swap: MemoryObject,
}

impl Memory {
    pub fn new(system: &System) -> Memory {
        Memory {
            memory: MemoryObject {
                total: system.total_memory(),
                used: system.used_memory(),
                available: system.free_memory(),
            },
            swap: MemoryObject {
                total: system.total_swap(),
                used: system.used_swap(),
                available: system.free_swap(),
            },
        }
    }

    pub fn memory(&self) -> &MemoryObject {
        &self.memory
    }

    pub fn swap(&self) -> &MemoryObject {
        &self.swap
    }

    pub fn update(&mut self, system: &System) {
        self.memory.used = system.used_memory();
        self.memory.available = system.available_memory();
        self.swap.used = system.used_swap();
        self.swap.available = system.free_swap();
    }
}