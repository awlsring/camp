use config::{AgentConfig, Config, ServerConfig};
use daemonize::Daemonize;
use std::env;
use std::error::Error;
use std::fs::File;
use std::sync::Arc;
use tokio::sync::Mutex;
use tokio::time::{sleep, Duration};

mod config;
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

    let config = config::load_config();

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

    tokio_main(config)
}

#[tokio::main]
async fn tokio_main(config: Config) -> Result<(), Box<dyn Error>> {
    info!("Initializing agent");

    let ctl = Arc::new(Mutex::new(SystemController::new()));
    let sctl = ctl.clone();

    info!("Starting agent loop");
    tokio::spawn(agent_loop(ctl, config.get_agent().clone()));

    info!("Starting server loop");
    server_loop(sctl, config.get_server().clone()).await;

    Ok(())
}

async fn agent_loop(ctl: Arc<Mutex<SystemController>>, config: AgentConfig) {
    loop {
        let mut lo = ctl.lock().await;

        lo.refresh().await;

        drop(lo);

        sleep(Duration::from_millis(config.get_interval())).await;
    }
}

async fn server_loop(ctl: Arc<Mutex<SystemController>>, config: ServerConfig) {
    start_server(ctl, config).await;
}
