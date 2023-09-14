use config::{AgentConfig, ReportingConfig, ServerConfig};
use daemonize::Daemonize;
use reporting::start_heartbeat;
use std::env;
use std::error::Error;
use std::fs::File;
use std::sync::Arc;
use tokio::sync::Mutex;
use tokio::time::{sleep, Duration};

mod config;
mod reporting;
mod server;
mod stats;

use log::{debug, error, info};
use server::http::start_server;
use stats::controller::SystemController;

fn main() -> Result<(), Box<dyn Error>> {
    if env::var("RUST_LOG").is_err() {
        env::set_var("RUST_LOG", "info");
    }
    env_logger::init();

    let env = env::var("RUNTIME_ENV").unwrap_or("dev".to_string());

    if env != "dev" {
        let log = File::create("/opt/campd/campd.log").unwrap();

        let daemonize = Daemonize::new()
            .working_directory("/opt/campd")
            .user("campd")
            .group("campd")
            .umask(0o027)
            .stderr(log) // all goes to err
            .privileged_action(|| "Executed before drop privileges");

        match daemonize.start() {
            Ok(_) => debug!("Daemonized"),
            Err(e) => {
                error!("Error, {}", e);
                std::process::exit(1)
            }
        }
    }

    tokio_main()
}

#[tokio::main]
async fn tokio_main() -> Result<(), Box<dyn Error>> {
    info!("Initializing agent");
    let mut config = config::load_config();

    let ctl = Arc::new(Mutex::new(SystemController::new()));
    let machine_id = ctl.lock().await.system().machine_id().to_owned();
    let sctl = ctl.clone();

    info!("Starting agent loop");
    let reporting_config = config.take_reporting();

    tokio::spawn(agent_loop(
        ctl,
        reporting_config.clone(),
        config.take_agent(),
    ));

    info!("Starting heartbeat loop");
    tokio::spawn(heartbeat_loop(reporting_config, machine_id));

    info!("Starting server loop");
    server_loop(sctl, config.take_server()).await;

    Ok(())
}

async fn heartbeat_loop(config: ReportingConfig, machine_id: String) {
    start_heartbeat(config, machine_id).await;
}

async fn agent_loop(
    ctl: Arc<Mutex<SystemController>>,
    reporting_cfg: ReportingConfig,
    agent_cfg: AgentConfig,
) {
    loop {
        let mut lo = ctl.lock().await;

        lo.refresh().await;

        drop(lo);

        sleep(Duration::from_millis(agent_cfg.get_interval())).await;
    }
}

async fn server_loop(ctl: Arc<Mutex<SystemController>>, config: ServerConfig) {
    start_server(ctl, config).await;
}
