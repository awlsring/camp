use std::{sync::Arc, str::FromStr, collections::HashMap};

use aws_smithy_http_server::Extension;
use camp_agent_server::{output::GetDiskOutput, output::ListDisksOutput, model::{DiskSummary, DiskType, DiskInterface as SmithyDiskInterface}, input::GetDiskInput, input::ListDisksInput, error};
use hw_info::{Disk, DiskInterface, DiskKind};

use crate::server::http::State;


pub async fn get_disk(input: GetDiskInput, state: Extension<Arc<State>>) -> Result<GetDiskOutput, error::GetDiskError> {
    let ctl = state.controller.lock().await;
    let disks = ctl.disks();

    let dev = input.name();

    let disk = match disks.get(dev) {
        Some(d) => {
            let summary = disk_to_summary(d);
            let output = GetDiskOutput { summary };
            return Ok(output)
        },
        None => return Err(error::GetDiskError::ResourceNotFoundException(error::ResourceNotFoundException { message: format!("Disk {} not found", input.name()) })),
    };
}

pub async fn list_disks(_input: ListDisksInput, state: Extension<Arc<State>>) -> Result<ListDisksOutput, error::ListDisksError> {
    let ctl = state.controller.lock().await;
    let disks = ctl.disks();
    let sums = disks_to_summaries(disks);
    let output = ListDisksOutput { summaries: sums };
    Ok(output)
}

pub fn disks_to_summaries(disks: &HashMap<String, Disk>) -> Vec<DiskSummary> {
    let mut summaries = Vec::new();
    for (_, disk) in disks {
        let sum = disk_to_summary(disk);
        summaries.push(sum);
    }

    summaries
}

pub fn disk_to_summary(disk: &Disk) -> DiskSummary {
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
        DiskKind::SSD => DiskType::Ssd,
        DiskKind::NVME => DiskType::Nvme,
        DiskKind::Unknown(_) => todo!(),
    };
    let sector_size = *disk.get_sector_size() as i32;
    let size_raw = *disk.get_size_raw() as i64;
    let size = disk.get_size_actual();

    DiskSummary {
        device,
        model,
        serial,
        vendor,
        interface: i,
        r#type: t,
        sector_size,
        size_raw,
        size_actual: *size,
    }
}