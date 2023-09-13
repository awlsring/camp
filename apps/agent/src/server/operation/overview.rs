use std::sync::Arc;

use aws_smithy_http_server::Extension;
use camp_agent_server::{output::GetOverviewOutput, model::{OverviewSummary, DiskSummary}, input::GetOverviewInput, error};

use crate::server::http::State;

use super::{system::system_to_summary, volume::volumes_to_summaries, memory::memory_to_summary, cpu::cpu_to_summary, network::network_interfaces_to_summaries, disk::disks_to_summaries, container::containers_to_summaries};


pub async fn get_overview(_input: GetOverviewInput, state: Extension<Arc<State>>) -> Result<GetOverviewOutput, error::GetOverviewError> {
    let ctl = state.controller.lock().await;
    let network = ctl.network();
    let cpu = ctl.cpu();
    let storage = ctl.storage();
    let disks = ctl.disks();
    let mem = ctl.memory();
    let sys = ctl.system();
    let conts = ctl.containers();

    let network = network_interfaces_to_summaries(network.network_interfaces());
    let cpu = cpu_to_summary(cpu);
    let memory = memory_to_summary(mem);
    let system = system_to_summary(sys);
    let volumes = volumes_to_summaries(storage.volumes());
    let disks: Vec<DiskSummary> = disks_to_summaries(disks);
    let containers = containers_to_summaries(conts);

    let sum = OverviewSummary {
        network,
        cpu,
        memory,
        system,
        volumes,
        disks,
        containers: match containers.len() {
            0 => None,
            _ => Some(containers)
        }
    };

    let output = GetOverviewOutput {
        summary: sum
    };

    Ok(output)
}