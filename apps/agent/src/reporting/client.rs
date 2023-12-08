use aws_smithy_client::{
    erase::{DynConnector, DynMiddleware},
    SdkError,
};
use aws_smithy_http::operation::Request;
use camp_local_rs::{
    config::AuthApiKey,
    types::{
        CpuArchitecture, DiskInterface, DiskType, IpAddressSummary, IpAddressVersion, MachineClass,
        MachineCpuSummary, MachineDiskSummary, MachineMemorySummary,
        MachineNetworkInterfaceSummary, MachineStatus, MachineSystemSummary, MachineVolumeSummary,
        ReportedMachineSummary, ReportedPowerCapabilitiesSummary,
    },
    Builder, Client, Config,
};
use http::{
    uri::{Authority, Scheme},
    Uri,
};
use log::{debug, error};
use std::str::FromStr;

use crate::stats::system_description::{
    CpuDescription, DiskDescription, HostDescription, MemoryDescription, NetworkDescription,
    SystemDescription, VolumeDescription,
};

fn rewrite_base_url(scheme: Scheme, authority: Authority) -> impl Fn(Request) -> Request + Clone {
    move |mut req| {
        let http_req = req.http_mut();
        let mut uri_parts = http_req.uri().clone().into_parts();
        uri_parts.authority = Some(authority.clone());
        uri_parts.scheme = Some(scheme.clone());
        *http_req.uri_mut() = Uri::from_parts(uri_parts).expect("failed to create uri from parts");
        req
    }
}

fn create_reporting_client(
    endpoint: &str,
    api_key: &str,
) -> Client<DynConnector, DynMiddleware<DynConnector>> {
    let authority = Authority::from_str(endpoint).expect("failed to parse authority");
    let raw_client = Builder::new()
        .rustls_connector(Default::default())
        .middleware_fn(rewrite_base_url(Scheme::HTTP, authority))
        .build_dyn();
    let config = Config::builder().api_key(AuthApiKey::from(api_key)).build();
    Client::with_config(raw_client, config)
}

pub(crate) struct ReportingClient {
    client: Client<DynConnector, DynMiddleware<DynConnector>>,
}

impl ReportingClient {
    pub fn new(endpoint: &str, api_key: &str) -> ReportingClient {
        ReportingClient {
            client: create_reporting_client(endpoint, api_key),
        }
    }

    pub async fn heartbeat(&self, id: &str) {
        let resp = self
            .client
            .heartbeat()
            .set_internal_identifier(Some(id.to_string()))
            .send()
            .await;

        match resp {
            Ok(_) => debug!("Sent heartbeat to Local server.",),
            Err(e) => error!("Failed to send heartbeat to Local server: {}", e),
        }
    }

    pub async fn check_registration(&self, id: &str) -> Result<bool, String> {
        let resp = self
            .client
            .describe_machine()
            .set_identifier(Some(id.to_string()))
            .send()
            .await;

        match resp {
            Ok(_) => Result::Ok(true),
            Err(e) => {
                let service_err = e.into_service_error();
                if service_err.is_resource_not_found_error() {
                    return Result::Ok(false);
                }
                Result::Err("Failed to check registration".to_string())
            }
        }
    }

    pub async fn register(&self, description: SystemDescription) {
        let id = description.machine_id.to_owned();
        let class = match description.class.to_lowercase().as_str() {
            "virtualmachine" => MachineClass::VirtualMachine,
            "baremetal" => MachineClass::BareMetal,
            "hypervisor" => MachineClass::Hypervisor,
            _ => MachineClass::UnknownValue,
        };

        let summary = system_decription_to_summary(description);
        let resp = self
            .client
            .register()
            .internal_identifier(id)
            .set_class(Some(class))
            .power_capabilities(ReportedPowerCapabilitiesSummary::builder().build())
            .set_system_summary(Some(summary))
            .set_callback_endpoint(Some(String::from("http://localhost:7032")))
            .set_callback_key(Some("a".to_string()))
            .send()
            .await;
        match resp {
            Ok(_) => debug!("Sent registration to Local server.",),
            Err(e) => error!("Failed to send registration to Local server: {}", e),
        }
    }

    pub async fn report_system_change(&self, description: SystemDescription) {
        let summary = system_decription_to_summary(description);
        let resp = self
            .client
            .report_system_change()
            .set_summary(Some(summary))
            .send()
            .await;
        match resp {
            Ok(_) => debug!("Sent system change to Local server.",),
            Err(e) => error!("Failed to send system change to Local server: {}", e),
        }
    }

    pub async fn report_stop(&self, id: &str) {
        let resp = self
            .client
            .report_status_change()
            .set_internal_identifier(Some(id.to_string()))
            .set_status(Some(MachineStatus::Stopped))
            .send()
            .await;
        match resp {
            Ok(_) => debug!("Sent shutdown to Local server.",),
            Err(e) => error!("Failed to send shutdown to Local server: {}", e),
        }
    }

    pub async fn report_start(&self, id: &str) {
        let resp = self
            .client
            .report_status_change()
            .set_internal_identifier(Some(id.to_string()))
            .set_status(Some(MachineStatus::Running))
            .send()
            .await;
        match resp {
            Ok(_) => debug!("Sent startup to Local server.",),
            Err(e) => error!("Failed to send startup to Local server: {}", e),
        }
    }
}

fn get_status_code_from_err<T>(err: SdkError<T>) -> u16 {
    let resp = err.raw_response();
    let r = match resp {
        Some(resp) => resp,
        None => return 0,
    };
    r.http().status().as_u16()
}

fn system_decription_to_summary(description: SystemDescription) -> ReportedMachineSummary {
    let mut disks = vec![];
    for disk in description.disks.iter() {
        disks.push(disk_description_to_disk_summary(disk.to_owned()));
    }

    let mut addrs = vec![];
    let mut nics = vec![];
    for net in description.network.iter() {
        let n = network_description_to_network_summary(net.to_owned());
        nics.push(n.clone());
        for address in n.addresses().iter() {
            for addr in address.iter() {
                addrs.push(addr.to_owned());
            }
        }
    }

    let mut volumes = vec![];
    for vol in description.volumes.iter() {
        volumes.push(volume_description_to_volume_summary(vol.to_owned()));
    }

    ReportedMachineSummary::builder()
        // .internal_identifier(description.machine_id)
        // .class(class)
        .system(host_description_to_system_summary(
            description.host.to_owned(),
        ))
        .cpu(cpu_description_to_cpu_summary(description.cpu.to_owned()))
        .memory(memory_description_to_memory_summary(
            description.memory.to_owned(),
        ))
        .set_disks(Some(disks))
        .set_volumes(Some(volumes))
        .set_network_interfaces(Some(nics))
        .set_addresses(Some(addrs))
        .build()
}

fn host_description_to_system_summary(description: HostDescription) -> MachineSystemSummary {
    MachineSystemSummary::builder()
        .family(description.family)
        .kernel_version(description.kernel_version)
        .os_pretty(description.os_pretty)
        .os_version(description.os_version)
        .os(description.os)
        .hostname(description.hostname)
        .build()
}

fn cpu_description_to_cpu_summary(description: CpuDescription) -> MachineCpuSummary {
    let arch = match description.architecture.to_lowercase().as_str() {
        "x86" => CpuArchitecture::X86,
        "x86_64" => CpuArchitecture::X86,
        "arm" => CpuArchitecture::Arm,
        _ => CpuArchitecture::UnknownValue,
    };
    MachineCpuSummary::builder()
        .cores(description.cores)
        .architecture(arch)
        .vendor(description.vendor)
        .model(description.model)
        .build()
}

fn memory_description_to_memory_summary(description: MemoryDescription) -> MachineMemorySummary {
    MachineMemorySummary::builder()
        .total(description.total)
        .build()
}

fn disk_description_to_disk_summary(description: DiskDescription) -> MachineDiskSummary {
    let interface = match description.interface_type.to_lowercase().as_str() {
        "sata" => DiskInterface::Sata,
        "scsi" => DiskInterface::Scsi,
        _ => DiskInterface::UnknownValue,
    };

    let disk_type = match description.disk_type.to_lowercase().as_str() {
        "ssd" => DiskType::Sdd,
        "hdd" => DiskType::Hdd,
        _ => DiskType::UnknownValue,
    };

    MachineDiskSummary::builder()
        .device(description.device)
        .model(description.model)
        .vendor(description.vendor)
        .interface(interface)
        .serial(description.serial_number)
        .size(description.size)
        .size_raw(description.size_raw)
        .sector_size(description.sector_size as i32)
        .r#type(disk_type)
        .build()
}

fn network_description_to_network_summary(
    description: NetworkDescription,
) -> MachineNetworkInterfaceSummary {
    let mut addresses = vec![];
    for addr in description.ip_addresses {
        let v = match addr.version.to_lowercase().as_str() {
            "v4" => IpAddressVersion::V4,
            "v6" => IpAddressVersion::V6,
            _ => IpAddressVersion::UnknownValue,
        };
        addresses.push(
            IpAddressSummary::builder()
                .address(addr.address)
                .version(v)
                .build(),
        );
    }
    let speed = description.speed.map(|s| s as i32);

    let mtu = description.mtu.map(|m| m as i32);

    MachineNetworkInterfaceSummary::builder()
        .name(description.name)
        .r#virtual(description.is_virtual)
        .set_mac_address(description.mac_address)
        .set_vendor(description.vendor)
        .set_duplex(description.duplex)
        .set_addresses(Some(addresses))
        .set_speed(speed)
        .set_mtu(mtu)
        .build()
}

fn volume_description_to_volume_summary(description: VolumeDescription) -> MachineVolumeSummary {
    MachineVolumeSummary::builder()
        .name(description.name)
        .mount_point(description.mount_point)
        .file_system(description.file_system)
        .total(description.total)
        .build()
}
