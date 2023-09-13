use std::sync::Arc;

use aws_smithy_http_server::Extension;
use camp_agent_server::{input::GetSystemInput, output::GetSystemOutput, model::SystemSummary, error};

use crate::{server::http::State, stats::system::System};


pub async fn get_system(_input: GetSystemInput, state: Extension<Arc<State>>) -> Result<GetSystemOutput, error::GetSystemError> {
    let ctl = state.controller.lock().await;
    let sys = ctl.system();

    let sum = system_to_summary(sys);

    let output = GetSystemOutput {
        summary: sum
    };

    Ok(output)
}

pub fn system_to_summary(system: &System) -> SystemSummary {
    let machine_id = system.machine_id().to_owned();
    let fam_opt = system.family().to_owned();
    let kernel = system.kernel_version().to_owned();
    let os = system.os().to_owned();
    let os_version = system.os_version().to_owned();
    let os_pretty = system.os_pretty().to_owned();
    let hostname = system.hostname().to_owned();
    let boot_time = system.boot_time().to_owned();
    let up_time = system.up_time().to_owned();

    SystemSummary {
        machine_id,
        family: fam_opt,
        kernel_version: kernel,
        os,
        os_version,
        os_pretty,
        hostname,
        boot_time: boot_time as i64,
        up_time: up_time as i64,
    }
}