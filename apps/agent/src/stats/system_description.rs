use serde::Serialize;

#[derive(Debug, Clone, Serialize)]
pub(crate) struct HostDescription {
    pub(crate) family: String,
    pub(crate) kernel_version: String,
    pub(crate) os_pretty: String,
    pub(crate) os_version: String,
    pub(crate) os: String,
    pub(crate) hostname: String,
}

#[derive(Debug, Clone, Serialize)]
pub(crate) struct CpuDescription {
    pub(crate) cores: i32,
    pub(crate) architecture: String,
    pub(crate) model: String,
    pub(crate) vendor: String,
}

#[derive(Debug, Clone, Serialize)]
pub(crate) struct MemoryDescription {
    pub(crate) total: i64,
}

#[derive(Debug, Clone, Serialize)]
pub(crate) struct DiskDescription {
    pub(crate) device: String,
    pub(crate) model: String,
    pub(crate) vendor: String,
    pub(crate) interface_type: String,
    pub(crate) serial_number: String,
    pub(crate) size: i64,
    pub(crate) size_raw: i64,
    pub(crate) sector_size: i16,
    pub(crate) disk_type: String,
}

#[derive(Debug, Clone, Serialize)]
pub(crate) struct IpAddressDescription {
    pub(crate) address: String,
    pub(crate) version: String,
}

#[derive(Debug, Clone, Serialize)]
pub(crate) struct NetworkDescription {
    pub(crate) name: String,
    pub(crate) mac_address: Option<String>,
    pub(crate) is_virtual: bool,
    pub(crate) vendor: Option<String>,
    pub(crate) mtu: Option<u16>,
    pub(crate) duplex: Option<String>,
    pub(crate) speed: Option<u16>,
    pub(crate) ip_addresses: Vec<IpAddressDescription>,
}

#[derive(Debug, Clone, Serialize)]
pub(crate) struct VolumeDescription {
    pub(crate) name: String,
    pub(crate) mount_point: String,
    pub(crate) file_system: String,
    pub(crate) total: i64,
}

#[derive(Debug, Clone, Serialize)]
pub(crate) struct SystemDescription {
    pub(crate) machine_id: String,
    pub(crate) class: String,
    pub(crate) host: HostDescription,
    pub(crate) cpu: CpuDescription,
    pub(crate) memory: MemoryDescription,
    pub(crate) disks: Vec<DiskDescription>,
    pub(crate) network: Vec<NetworkDescription>,
    pub(crate) volumes: Vec<VolumeDescription>,
}
