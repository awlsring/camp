use std::sync::Arc;

use aws_smithy_http_server::Extension;
use camp_agent_server::{output::GetCpuOutput, model::CpuSummary, model::CoreUtilization, input::GetCpuInput, error};

use crate::{server::http::State, stats::cpu::Cpu};


pub async fn get_cpu(_input: GetCpuInput, state: Extension<Arc<State>>) -> Result<GetCpuOutput, error::GetCpuError> {
    let ctl = state.controller.lock().await;
    let cpu = ctl.cpu();

    let sum = cpu_to_summary(cpu);

    let output = GetCpuOutput {
        summary: sum
    };

    Ok(output)
}

pub fn cpu_to_summary(cpu: &Cpu) -> CpuSummary {
    let cores = cpu.core_count();
    let architecture = cpu.architecture();
    let model = cpu.brand();
    let vendor = cpu.vendor();

    let mut utils = Vec::new();
    for core in cpu.cores() {
        let name = core.name().to_string();
        let usage = *core.usage();
        let frequency =  *core.frequency()as f32;
        let util = CoreUtilization {
            name,
            usage,
            frequency,
        };

        utils.push(util);
    }

    CpuSummary { 
        utilization: utils,
        cores: cores as i32,
        architecture: architecture.to_string(),
        model: model.to_string(),
        vendor: vendor.to_string(),
    }
}