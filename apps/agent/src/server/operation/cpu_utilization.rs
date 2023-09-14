use std::sync::Arc;

use crate::server::http::State;
use aws_smithy_http_server::Extension;
use camp_agent_server::{error, input::GetCpuUtilizationInput, output::GetCpuUtilizationOutput};

pub async fn get_cpu_utilization(
    _: GetCpuUtilizationInput,
    _: Extension<Arc<State>>,
) -> Result<GetCpuUtilizationOutput, error::GetCpuUtilizationError> {
    Err(error::GetCpuUtilizationError::UnauthorizedException(
        error::UnauthorizedException {
            message: "Unauthorized".to_string(),
        },
    ))
}
