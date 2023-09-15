use camp_agent_server::{error, input, output};

pub async fn check_health(
    _input: input::HealthInput,
) -> Result<output::HealthOutput, error::HealthError> {
    Ok(output::HealthOutput { success: true })
}
