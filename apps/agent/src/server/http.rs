use std::{net::SocketAddr, sync::Arc};

use aws_smithy_http_server::{
    extension::OperationExtensionExt,
    instrumentation::InstrumentExt,
    plugin::{IdentityPlugin, PluginPipeline},
    request::request_id::ServerRequestIdProviderLayer,
    AddExtensionLayer,
};

use camp_agent_client::operation::get_volume_utilization;
use log::{error, info};
use tokio::sync::Mutex;

use crate::{
    config::ServerConfig,
    server::operation::{
        cpu_utilization::get_cpu_utilization, disk_utilization::get_disk_utilization,
        memory_utilization::get_memory_utilization,
        network_utilization::get_network_interface_utilization, overview::get_overview,
        uptime::get_uptime, volume_utilization::get_volume_utilization,
    },
    stats::controller::SystemController,
};

use smithy_common::auth::controller::AuthController;
use smithy_common::auth::plugin::AuthExtension;
use smithy_common::print::plugin::PrintExt;

use camp_agent_server::CampAgent;
use camp_agent_server::{error, input, output};

pub const DEFAULT_ADDRESS: &str = "0.0.0.0";

#[derive(Clone)]
struct Config;

pub struct State {
    pub controller: Arc<Mutex<SystemController>>,
}

impl State {
    pub fn new(ctl: Arc<Mutex<SystemController>>) -> State {
        State { controller: ctl }
    }
}

pub async fn check_health(
    _input: input::HealthInput,
) -> Result<output::HealthOutput, error::HealthError> {
    Ok(output::HealthOutput { success: true })
}

pub async fn start_server(ctl: Arc<Mutex<SystemController>>, config: ServerConfig) {
    // TODO: Add config where keys can be stored and retrived
    let auth_controller = AuthController::new(config.no_auth_operations(), config.allowed_keys());

    let plugins = PluginPipeline::new()
        .print()
        .auth(auth_controller.into(), Config)
        .insert_operation_extension()
        .instrument();

    let app = CampAgent::builder_with_plugins(plugins, IdentityPlugin)
        .health(check_health)
        .get_overview(get_overview)
        .get_cpu_utilization(get_cpu_utilization)
        .get_memory_utilization(get_memory_utilization)
        .get_disk_utilization(get_disk_utilization)
        .get_network_interface_utilization(get_network_interface_utilization)
        .get_volume_utilization(get_volume_utilization)
        .get_uptime(get_uptime)
        .build()
        .expect("failed to build an instance of GethAgent");

    // create state to add to request
    let state = State::new(ctl);
    let app = app
        .layer(&AddExtensionLayer::new(Arc::new(state)))
        .layer(&ServerRequestIdProviderLayer::new());

    let make_app = app.into_make_service_with_connect_info::<SocketAddr>();

    info!(
        "Starting server on: {}:{}",
        DEFAULT_ADDRESS,
        config.get_server_port()
    );
    let bind: SocketAddr = format!("{}:{}", DEFAULT_ADDRESS, config.get_server_port())
        .parse()
        .expect("unable to parse the server bind address and port");
    let server = hyper::Server::bind(&bind).serve(make_app);

    // Run forever-ish...
    if let Err(err) = server.await {
        error!("server error: {}", err);
    }
}
