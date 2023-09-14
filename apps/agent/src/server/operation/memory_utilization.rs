use std::sync::Arc;

use crate::server::http::State;
use aws_smithy_http_server::Extension;
use camp_agent_server::{
    error, input::GetMemoryUtilizationInput, output::GetMemoryUtilizationOutput,
};

pub async fn get_memory_utilization(
    _: GetMemoryUtilizationInput,
    _: Extension<Arc<State>>,
) -> Result<GetMemoryUtilizationOutput, error::GetMemoryUtilizationError> {
    Err(error::GetMemoryUtilizationError::UnauthorizedException(
        error::UnauthorizedException {
            message: "Unauthorized".to_string(),
        },
    ))
}
