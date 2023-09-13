use std::sync::Arc;

use aws_smithy_http_server::Extension;
use camp_agent_server::{
    error, input::GetOverviewInput, model::OverviewSummary, output::GetOverviewOutput,
};

use crate::server::http::State;

use super::conversion::{
    cpu_to_summary, disks_to_summaries, memory_to_summary, network_interfaces_to_summaries,
    system_to_summary, volumes_to_summaries,
};

pub async fn get_overview(
    _input: GetOverviewInput,
    state: Extension<Arc<State>>,
) -> Result<GetOverviewOutput, error::GetOverviewError> {
    let ctl = state.controller.lock().await;
    let network = ctl.network();
    let cpu = ctl.cpu();
    let storage = ctl.storage();
    let disks = ctl.disks();
    let mem = ctl.memory();
    let sys = ctl.system();

    let network_interfaces = network_interfaces_to_summaries(network.network_interfaces());
    let cpu = cpu_to_summary(cpu);
    let memory = memory_to_summary(mem);
    let system = system_to_summary(sys);
    let volumes = volumes_to_summaries(storage.volumes());
    let disks = disks_to_summaries(disks);
    let addresses = vec![];

    let sum = OverviewSummary {
        network_interfaces,
        cpu,
        memory,
        system,
        volumes,
        disks,
        addresses,
    };

    let output = GetOverviewOutput { summary: sum };

    Ok(output)
}
