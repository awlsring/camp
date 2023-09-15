use std::fs::File;
use std::io::Write;
use std::path;
use std::sync::Arc;
use tokio::select;
use tokio::sync::{watch, Mutex};
use tokio::time::sleep;

use crate::stats::controller::SystemController;
use log::{debug, error};

use crate::reporting::client::ReportingClient;

async fn write_snapshot(path: &str, snapshot_data: String) {
    debug!("Writing snapshot to file");
    let file_name = format!("system_snapshot.json");
    let file_path = path::Path::new(path).join(file_name);
    if let Ok(mut file) = File::create(file_path) {
        if let Err(e) = file.write_all(snapshot_data.as_bytes()) {
            error!("Failed to write snapshot to file: {}", e);
        }
    } else {
        error!("Failed to create snapshot file");
    }
}

async fn update(
    ctl: &Arc<Mutex<SystemController>>,
    client: &Arc<ReportingClient>,
    snapshot_location: &str,
) {
    let mut lo = ctl.lock().await;
    let updated = lo.refresh().await;
    let description = lo.get_system_description();
    drop(lo);

    if updated {
        let serialized_description = serde_json::to_string(&description).unwrap();
        write_snapshot(snapshot_location, serialized_description).await;
        client.report_system_change(description).await;
    }
}

pub(crate) async fn system_stats_loop(
    ctl: Arc<Mutex<SystemController>>,
    reporting_client: Arc<ReportingClient>,
    change_location: &str,
    interval: u64,
    mut stop_rx: watch::Receiver<()>,
) {
    loop {
        select! {
            biased;

            _ = stop_rx.changed() => {
                debug!("stopping system stats process.");
                break;
            },
            _ = update(&ctl, &reporting_client, change_location) => {
                sleep(tokio::time::Duration::from_secs(interval)).await;
            }
        };
    }
}
