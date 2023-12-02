use daemonize::Daemonize;
use stats::system_description::SystemDescription;
use std::env;
use std::error::Error;
use std::fs::File;
use std::sync::Arc;
use tokio::select;
use tokio::signal::unix::{signal, SignalKind};
use tokio::sync::{watch, Mutex};

mod config;
mod reporting;
mod server;
mod stats;
mod tasks;

use log::{debug, error, info};
use server::http::start_server;
use stats::controller::SystemController;

use crate::reporting::client::ReportingClient;
use crate::tasks::heartbeat::heartbeat_loop;
use crate::tasks::system_collection::system_stats_loop;

fn main() -> Result<(), Box<dyn Error>> {
    initialize_logging();
    daemonize_process()?;
    tokio_main()
}

#[tokio::main]
async fn tokio_main() -> Result<(), Box<dyn Error>> {
    info!("Initializing agent");
    let (stop_tx, stop_rx) = watch::channel(());

    let mut config = config::load_config();
    let reporting_config = config.take_reporting();
    let snapshot_location = determine_snapshot_location();

    let ctl = Arc::new(Mutex::new(SystemController::new()));
    let description = ctl.lock().await.get_system_description();
    let machine_id = ctl.lock().await.system().machine_id().to_owned();
    let camp_local = Arc::new(ReportingClient::new(
        reporting_config.get_endpoint(),
        reporting_config.get_api_key(),
    ));
    setup_signal_handlers(stop_tx);

    info!("checking registration");
    check_registration(camp_local.clone(), description).await;

    info!("sending startup notification");
    camp_local.report_start(&machine_id).await;

    info!("Starting tasks");
    tokio::join!(
        heartbeat_loop(machine_id.as_str(), camp_local.clone(), 15, stop_rx.clone()),
        system_stats_loop(
            ctl.clone(),
            camp_local.clone(),
            snapshot_location.as_str(),
            config.take_agent().get_interval(),
            stop_rx.clone(),
        ),
        start_server(ctl, config.take_server(), stop_rx)
    );

    info!("all tasks closed, sending shutdown notification");
    camp_local.report_stop(&machine_id).await;

    Ok(())
}

async fn check_registration(client: Arc<ReportingClient>, description: SystemDescription) {
    match client
        .check_registration(description.machine_id.as_str())
        .await
    {
        Ok(registered) => {
            if !registered {
                debug!("machine is not registered, registering");
                client.register(description).await;
            }
        }
        Err(e) => {
            error!("failed to check registration: {}", e);
            std::process::exit(1);
        }
    }
}

fn initialize_logging() {
    if std::env::var("RUST_LOG").is_err() {
        std::env::set_var("RUST_LOG", "debug");
    }
    env_logger::init();
}

fn determine_snapshot_location() -> String {
    let env = env::var("RUNTIME_ENV").unwrap_or("dev".to_string());
    if env != "dev" {
        return "/opt/campd".to_string();
    }
    ".".to_string()
}

fn daemonize_process() -> Result<(), Box<dyn Error>> {
    let env = env::var("RUNTIME_ENV").unwrap_or("dev".to_string());
    if env != "dev" {
        let log = File::create("/opt/campd/campd.log")?;
        let daemonize = Daemonize::new()
            .working_directory("/opt/campd")
            .user("campd")
            .group("campd")
            .umask(0o027)
            .stderr(log)
            .privileged_action(|| "Executed before drop privileges");

        match daemonize.start() {
            Ok(_) => debug!("Daemonized"),
            Err(e) => {
                error!("Error, {}", e);
                std::process::exit(1);
            }
        }
    }

    Ok(())
}

fn setup_signal_handlers(stop_tx: watch::Sender<()>) {
    let mut sigterm = signal(SignalKind::terminate()).unwrap();
    let mut sigint = signal(SignalKind::interrupt()).unwrap();
    tokio::spawn(async move {
        loop {
            select! {
                _ = sigterm.recv() => {
                    info!("Received SIGTERM, starting shutdown of tasks.");
                    stop_tx.send(()).unwrap();
                }
                _ = sigint.recv() => {
                    info!("Received SIGINT, starting shutdown of tasks.");
                    stop_tx.send(()).unwrap();
                }
            };
        }
    });
}
