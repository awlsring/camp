use std::sync::Arc;

use crate::server::http::State;
use aws_smithy_http_server::Extension;
use camp_agent_server::{
    error, input::GetVolumeUtilizationInput, output::GetVolumeUtilizationOutput,
};

pub async fn get_volume_utilization(
    _: GetVolumeUtilizationInput,
    _: Extension<Arc<State>>,
) -> Result<GetVolumeUtilizationOutput, error::GetVolumeUtilizationError> {
    Err(error::GetVolumeUtilizationError::UnauthorizedException(
        error::UnauthorizedException {
            message: "Unauthorized".to_string(),
        },
    ))
}
