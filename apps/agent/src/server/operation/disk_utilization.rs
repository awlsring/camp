use std::sync::Arc;

use crate::server::http::State;
use aws_smithy_http_server::Extension;
use camp_agent_server::{error, input::GetDiskUtilizationInput, output::GetDiskUtilizationOutput};

pub async fn get_disk_utilization(
    _: GetDiskUtilizationInput,
    _: Extension<Arc<State>>,
) -> Result<GetDiskUtilizationOutput, error::GetDiskUtilizationError> {
    Err(error::GetDiskUtilizationError::UnauthorizedException(
        error::UnauthorizedException {
            message: "Unauthorized".to_string(),
        },
    ))
}
