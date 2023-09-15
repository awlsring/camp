use std::sync::Arc;
use tokio::select;
use tokio::sync::watch;
use tokio::time::sleep;

use log::debug;

use crate::reporting::client::ReportingClient;

pub(crate) async fn heartbeat_loop(
    machine_id: &str,
    client: Arc<ReportingClient>,
    interval: u64,
    mut stop_rx: watch::Receiver<()>,
) {
    loop {
        select! {
            biased;

            _ = stop_rx.changed() => {
                debug!("stopping heartbeat process.");
                break;
            },
            _ = client.heartbeat(machine_id) => {
                sleep(tokio::time::Duration::from_secs(interval)).await;
            }
        };
    }
}
