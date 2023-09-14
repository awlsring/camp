use std::sync::Arc;

use crate::server::http::State;
use aws_smithy_http_server::Extension;
use aws_smithy_types::DateTime;
use camp_agent_server::{
    error, input::GetUptimeInput, model::UptimeSummary, output::GetUptimeOutput,
};

pub async fn get_uptime(
    _: GetUptimeInput,
    state: Extension<Arc<State>>,
) -> Result<GetUptimeOutput, error::GetUptimeError> {
    let ctl = state.controller.lock().await;
    let sys = ctl.system();

    let up_time = sys.up_time().to_owned() as i64;
    let boot_time = DateTime::from_secs(sys.boot_time().to_owned() as i64);

    let output = GetUptimeOutput {
        summary: UptimeSummary { up_time, boot_time },
    };

    Ok(output)
}
