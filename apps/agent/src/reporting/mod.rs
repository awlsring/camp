use aws_smithy_client::erase::{DynConnector, DynMiddleware};
use aws_smithy_http::operation::Request;
use camp_local_rs::{config::AuthApiKey, Builder, Client, Config};
use http::{
    uri::{Authority, Scheme},
    Uri,
};
use log::{debug, error, info};
use std::str::FromStr;

use crate::config::ReportingConfig;

pub fn rewrite_base_url(
    scheme: Scheme,
    authority: Authority,
) -> impl Fn(Request) -> Request + Clone {
    move |mut req| {
        let http_req = req.http_mut();
        let mut uri_parts = http_req.uri().clone().into_parts();
        uri_parts.authority = Some(authority.clone());
        uri_parts.scheme = Some(scheme.clone());
        *http_req.uri_mut() = Uri::from_parts(uri_parts).expect("failed to create uri from parts");
        req
    }
}

fn make_reporting_client(
    config: ReportingConfig,
) -> Client<DynConnector, DynMiddleware<DynConnector>> {
    let authority = Authority::from_str(config.get_endpoint()).expect("failed to parse authority");
    let raw_client = Builder::new()
        .rustls_connector(Default::default())
        .middleware_fn(rewrite_base_url(Scheme::HTTP, authority))
        .build_dyn();
    let config = Config::builder()
        .api_key(AuthApiKey::from(config.get_api_key()))
        .build();
    Client::with_config(raw_client, config)
}

pub async fn start_heartbeat(config: ReportingConfig, machine_id: String) {
    let endpoint = config.get_endpoint().to_owned();
    let client = make_reporting_client(config);
    loop {
        let id = machine_id.to_owned();
        let resp = client
            .heartbeat()
            .set_internal_identifier(Some(id))
            .send()
            .await;

        match resp {
            Ok(_) => debug!(
                "Send heartbeat to Local server at {:?}",
                endpoint.to_string()
            ),
            Err(e) => error!("Error sending heartbeat: {:?}", e),
        }
        tokio::time::sleep(tokio::time::Duration::from_secs(30)).await;
    }
}
