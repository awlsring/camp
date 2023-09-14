use std::sync::Arc;

use crate::server::http::State;
use aws_smithy_http_server::Extension;
use camp_agent_server::{
    error, input::GetNetworkInterfaceUtilizationInput, output::GetNetworkInterfaceUtilizationOutput,
};

pub async fn get_network_interface_utilization(
    _: GetNetworkInterfaceUtilizationInput,
    _: Extension<Arc<State>>,
) -> Result<GetNetworkInterfaceUtilizationOutput, error::GetNetworkInterfaceUtilizationError> {
    Err(
        error::GetNetworkInterfaceUtilizationError::UnauthorizedException(
            error::UnauthorizedException {
                message: "Unauthorized".to_string(),
            },
        ),
    )
}
