use std::{collections::HashMap, net::IpAddr};

use camp_agent_server::model::{
    CpuArchitecture, DiskInterface as SmithyDiskInterface, DiskType, IpAddressSummary,
    IpAddressVersion, MachineCpuSummary, MachineDiskSummary, MachineMemorySummary,
    MachineNetworkInterfaceSummary, MachineSystemSummary, MachineVolumeSummary,
};
use hw_info::{Disk, DiskInterface, DiskKind};

use crate::stats::{
    cpu::Cpu,
    disk::Disk as StatDisk,
    memory::Memory,
    network::{AddressKind, NetworkInterface},
    system::System,
};

pub fn system_to_summary(system: &System) -> MachineSystemSummary {
    let machine_id = system.machine_id().to_owned();
    let fam_opt = system.family().to_owned();
    let kernel = system.kernel_version().to_owned();
    let os = system.os().to_owned();
    let os_version = Some(system.os_version().to_owned());
    let os_pretty = Some(system.os_pretty().to_owned());
    let hostname = Some(system.hostname().to_owned());

    MachineSystemSummary {
        family: Some(fam_opt),
        kernel_version: Some(kernel),
        os: Some(os),
        os_version,
        os_pretty,
        hostname,
    }
}

pub fn cpu_to_summary(cpu: &Cpu) -> MachineCpuSummary {
    let cores = cpu.core_count();
    let arch = cpu.architecture();
    let model = cpu.brand();
    let vendor = cpu.vendor();

    let architecture = match arch.as_str() {
        "x86" => CpuArchitecture::X86,
        "arm" => CpuArchitecture::Arm,
        _ => CpuArchitecture::Unknown,
    };

    MachineCpuSummary {
        cores: cores as i32,
        architecture,
        model: Some(model.to_string()),
        vendor: Some(vendor.to_string()),
    }
}

pub fn memory_to_summary(mem: &Memory) -> MachineMemorySummary {
    MachineMemorySummary {
        total: *mem.memory().total() as i64,
    }
}

pub fn disks_to_summaries(disks: &HashMap<String, Disk>) -> Vec<MachineDiskSummary> {
    let mut summaries = Vec::new();
    for (_, disk) in disks {
        let sum = disk_to_summary(disk);
        summaries.push(sum);
    }

    summaries
}

pub fn disk_to_summary(disk: &Disk) -> MachineDiskSummary {
    let device = disk.get_device().to_owned();
    let model = disk.get_model().to_owned();
    let serial = disk.get_serial().to_owned();
    let vendor = disk.get_vendor().to_owned();
    let interface = disk.get_interface().to_owned();
    let i = match disk.get_interface() {
        DiskInterface::SATA => SmithyDiskInterface::Sata,
        DiskInterface::SCSI => SmithyDiskInterface::Scsi,
        DiskInterface::PCI_E => SmithyDiskInterface::PciE,
        _ => SmithyDiskInterface::Unknown,
    };
    let kind = disk.get_kind().to_owned();
    let t = match disk.get_kind() {
        DiskKind::HDD => DiskType::Hdd,
        DiskKind::SSD => DiskType::Sdd,
        DiskKind::NVME => DiskType::Sdd,
        DiskKind::Unknown(_) => todo!(),
    };
    let sector_size = *disk.get_sector_size() as i32;
    let size_raw = *disk.get_size_raw() as i64;
    let size = disk.get_size_actual();

    MachineDiskSummary {
        device,
        model: Some(model),
        serial: Some(serial),
        vendor: Some(vendor),
        interface: i,
        r#type: t,
        sector_size: Some(sector_size),
        size_raw: Some(size_raw),
        size: *size,
    }
}

pub fn network_interfaces_to_summaries(
    nics: Vec<&NetworkInterface>,
) -> Vec<MachineNetworkInterfaceSummary> {
    let mut summaries = Vec::new();
    for nic in nics {
        let sum = network_interface_to_summary(nic);
        summaries.push(sum);
    }

    summaries
}

pub fn network_interface_to_summary(iface: &NetworkInterface) -> MachineNetworkInterfaceSummary {
    let name = iface.name().to_string();

    let mut addresses = Vec::new();
    for addr in iface.addresses() {
        let version = addr.version();
        let address = addr.address().to_string();
        let netmask = handle_optional_ip(addr.netmask());
        let broadcast = handle_optional_ip(addr.broadcast());

        let addr = IpAddressSummary {
            version: address_kind_to_smithy(version),
            address,
        };

        addresses.push(addr);
    }

    let mtu: Option<i32> = match iface.mtu() {
        Some(mtu) => Some(*mtu as i32),
        None => None,
    };

    let speed = match iface.speed() {
        Some(speed) => Some(*speed as i32),
        None => None,
    };

    MachineNetworkInterfaceSummary {
        name,
        addresses,
        r#virtual: *iface.is_virtual(),
        mac_address: iface.mac().to_owned(),
        vendor: iface.vendor().to_owned(),
        mtu,
        duplex: iface.duplex().to_owned(),
        speed,
    }
}

pub fn volumes_to_summaries(disks: Vec<&StatDisk>) -> Vec<MachineVolumeSummary> {
    let mut summaries = Vec::new();
    for disk in disks {
        let sum = volume_to_summary(disk);
        summaries.push(sum);
    }

    summaries
}

pub fn volume_to_summary(disk: &StatDisk) -> MachineVolumeSummary {
    let name = disk.name().to_owned();
    let mount_point = disk.mount_point().to_owned();
    let file_system = Some(disk.file_system().to_owned());
    let total = *disk.total_space() as i64;

    MachineVolumeSummary {
        name,
        mount_point,
        file_system,
        total,
    }
}

fn address_kind_to_smithy(kind: &AddressKind) -> IpAddressVersion {
    match kind {
        AddressKind::V4 => IpAddressVersion::V4,
        AddressKind::V6 => IpAddressVersion::V6,
        AddressKind::V6Local => IpAddressVersion::V6,
    }
}

fn handle_optional_ip(ip: &Option<IpAddr>) -> Option<String> {
    ip.as_ref().map(|ip| ip.to_string())
}
