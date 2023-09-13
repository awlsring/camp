use std::{sync::Arc, str::FromStr};

use aws_smithy_http_server::Extension;
use camp_agent_server::{output::GetVolumeOutput, output::ListVolumesOutput, model::{VolumeSummary, VolumeType}, input::GetVolumeInput, input::ListVolumesInput, error};
use sysinfo::DiskKind;

use crate::{server::http::State, stats::disk::Disk};


pub async fn get_volume(input: GetVolumeInput, state: Extension<Arc<State>>) -> Result<GetVolumeOutput, error::GetVolumeError> {
    let ctl = state.controller.lock().await;
    let volumes = ctl.storage();

    let volume = volumes.get_volume(input.name());

    match volume {
        Some(d) => {
            let sum = volume_to_summary(d);
            let output = GetVolumeOutput { summary: sum };
            Ok(output)
        }
        None => Err(error::GetVolumeError::ResourceNotFoundException(error::ResourceNotFoundException { message: format!("Volume {} not found", input.name()) }))
    }
}

pub async fn list_volumes(_input: ListVolumesInput, state: Extension<Arc<State>>) -> Result<ListVolumesOutput, error::ListVolumesError> {
    let ctl = state.controller.lock().await;
    let volumes = ctl.storage();
    let sums = volumes_to_summaries(volumes.volumes());
    let output = ListVolumesOutput { summaries: sums };
    Ok(output)
}

pub fn volumes_to_summaries(disks: Vec<&Disk>) -> Vec<VolumeSummary> {
    let mut summaries = Vec::new();
    for disk in disks {
        let sum = volume_to_summary(disk);
        summaries.push(sum);
    }

    summaries
}

pub fn volume_to_summary(disk: &Disk) -> VolumeSummary {
    let name = disk.name().to_owned();
    let mount_point = disk.mount_point().to_owned();
    let file_system = disk.file_system().to_owned();
    let total = *disk.total_space() as i64;
    let available = *disk.available_space() as i64;
    let removeable = *disk.is_removable();
    let t = match disk.disk_type() {
        DiskKind::HDD => "HDD",
        DiskKind::SSD => "SSD",
        DiskKind::Unknown(_) => "HDD",
    };
    let kind = VolumeType::from_str(t);

    VolumeSummary {
        name,
        mount_point,
        file_system,
        used_space: total - available,
        total_space: total,
        available_space: available,
        removeable,
        r#type: kind.unwrap(),
    }
}