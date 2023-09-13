use std::{sync::Arc, collections::HashMap};

use async_stream::stream;
use aws_smithy_http_server::Extension;
use containers::{Container, Port, ContainerProtocol, Volume, Network};
use camp_agent_server::{output::{StreamContainerLogsOutput, GetContainerOutput, ListContainersOutput, StreamContainerStatisticsOutput}, input::{StreamContainerLogsInput, ListContainersInput, GetContainerInput, StreamContainerStatisticsInput}, error::{self, ResourceNotFoundException}, model::{Logs, LogLine, ContainerSummary, ContainerState, ContainerPortBinding, ContainerPortProtocol, ContainerVolume, ContainerNetwork, ContainerStatistics, ContainerType}};
use crate::server::http::State;

pub fn containers_to_summaries(conts: &HashMap<String, Container>) -> Vec<ContainerSummary> {
    let mut summaries = Vec::new();
    for (_, c) in conts {
        let sum = container_to_summary(c);
        summaries.push(sum);
    }

    summaries
}

pub fn container_to_summary(container: &Container) -> ContainerSummary {
    let id = container.id().to_string();
    let name = container.name().to_string();
    let image = container.image().to_string();
    let created = container.created().timestamp() as i64;
    let state = match container.state() {
        containers::ContainerState::Running => ContainerState::Running,
        containers::ContainerState::Created => ContainerState::Created,
        containers::ContainerState::Restarting => ContainerState::Restarting,
        containers::ContainerState::Removing => ContainerState::Removing,
        containers::ContainerState::Paused => ContainerState::Paused,
        containers::ContainerState::Stopped => ContainerState::Stopped,
        containers::ContainerState::Dead => ContainerState::Dead,
        containers::ContainerState::Empty => ContainerState::Empty,
        _ => ContainerState::Unknown,
    };
    let ports = ports_to_port_summaries(container.ports());
    let volumes = volume_to_volume_summaries(container.volumes());
    let networks = network_to_network_summaries(container.networks());
    let labels = container.labels().to_owned();
    let command = match container.command() {
        Some(c) => Some(c.to_string()),
        None => None,
    };
    let environment = match container.environment() {
        Some(e) => Some(e.to_owned()),
        None => None,
    };
    let started = match container.started() {
        Some(s) => optional_timestamp(s.timestamp()),
        None => None,
    };
    let finished = match container.finished() {
        Some(f) => optional_timestamp(f.timestamp()),
        None => None,
    };
    let statistics = match container.statistics() {
        Some(s) => Some(container_stats_to_summary(s)),
        None => None,
    };
    let container_type = match container.type_() {
        containers::ContainerType::Docker => ContainerType::Docker,
        _ => ContainerType::Unknown,
    };


    ContainerSummary {
        id,
        name,
        image,
        created,
        state,
        command,
        environment,
        started,
        finished,
        statistics,
        container_type,
        ports: match ports.len() {
            0 => None,
            _ => Some(ports),
        },
        volumes: match volumes.len() {
            0 => None,
            _ => Some(volumes),
        },
        networks: match networks.len() {
            0 => None,
            _ => Some(networks),
        },
        labels: match labels {
            None => None,
            Some(l) => Some(l.to_owned()),
        },
    }
}

fn optional_timestamp(ts: i64) -> Option<i64> {
    if ts < 0 {
        None
    } else {
        Some(ts)
    }
}

fn container_stats_to_summary(stats: &containers::ContainerStatistics) -> ContainerStatistics {
    let cpu_utilization = Some(stats.cpu_utilization() as f32);
    let memory_utilization = Some(stats.memory_utilization() as f32);
    let memory_usage = Some(stats.memory_usage() as i64);
    let memory_limit = Some(stats.memory_limit() as i64);
    let network_rx_bytes = Some(stats.network_rx_bytes() as i64);
    let network_tx_bytes = Some(stats.network_tx_bytes() as i64);
    let block_read_bytes = Some(stats.block_read_bytes() as i64);
    let block_write_bytes = Some(stats.block_write_bytes() as i64);

    ContainerStatistics {
        cpu_utilization,
        memory_utilization,
        memory_usage,
        memory_limit,
        network_rx_bytes,
        network_tx_bytes,
        block_read_bytes,
        block_write_bytes,
    }
}

fn network_to_network_summaries(networks: &Vec<Network>) -> Vec<ContainerNetwork> {
    let mut summaries = Vec::new();
    for network in networks {
        let name = network.name().to_string();
        let network_id = network.network_id().to_string();
        let endpoint_id = network.endpoint_id().to_string();

        let summary = ContainerNetwork {
            name,
            network_id,
            endpoint_id,
        };

        summaries.push(summary);
    }

    summaries
}

fn volume_to_volume_summaries(volumes: &Vec<Volume>) -> Vec<ContainerVolume> {
    let mut summaries = Vec::new();
    for volume in volumes {
        let source = match volume.source() {
            Some(s) => Some(s.to_string()),
            None => None,
        };
        let destination = match volume.destination() {
            Some(d) => Some(d.to_string()),
            None => None,
        };
        let mode = match volume.mode() {
            Some(m) => Some(m.to_string()),
            None => None,
        };

        let summary = ContainerVolume {
            source,
            destination,
            mode,
        };

        summaries.push(summary);
    }

    summaries
}

fn ports_to_port_summaries(ports: &Vec<Port>) -> Vec<ContainerPortBinding> {
    let mut bindings = Vec::new();
    for port in ports {
        // let p = port.to_owned();
        let hp = port.host();
        let host_port = match hp {
            Some(p) => Some(p.to_owned() as i32),
            None => None,
        };
        let container_port = port.container() as i32;
        let protocol = match port.protocol() {
            ContainerProtocol::TCP => ContainerPortProtocol::Tcp,
            ContainerProtocol::UDP => ContainerPortProtocol::Udp,
            ContainerProtocol::SCTP => ContainerPortProtocol::Sctp,
            _ => ContainerPortProtocol::Unknown,
        };
        let host_addresses = port.host_addresses().to_owned();

        let binding = ContainerPortBinding {
            host_port,
            container_port,
            protocol,
            host_addresses,
        };

        bindings.push(binding);
    }

    bindings
}

pub async fn get_container(input: GetContainerInput, state: Extension<Arc<State>>) -> Result<GetContainerOutput, error::GetContainerError> {
    let ctl = state.controller.lock().await;
    let containers = ctl.containers();

    let id = input.id.to_string();

    let container = match containers.get(&id) {
        Some(c) => {
            let summary = container_to_summary(c);
            let output = GetContainerOutput { summary };
            return Ok(output)
        },
        None => return Err(error::GetContainerError::ResourceNotFoundException(ResourceNotFoundException { message: format!("Container not found") })),
    };
}

pub async fn list_containers(input: ListContainersInput, state: Extension<Arc<State>>) -> Result<ListContainersOutput, error::ListContainersError> {
    let ctl = state.controller.lock().await;
    let containers = ctl.containers();

    let summaries = containers_to_summaries(containers);

    let output = ListContainersOutput { summaries };
    Ok(output)
}

pub async fn stream_container_statistics(input: StreamContainerStatisticsInput, state: Extension<Arc<State>>) -> Result<StreamContainerStatisticsOutput, error::StreamContainerStatisticsError> {
    // let ctl = state.controller.lock().await;
    // let containers = ctl.containers();

    // let id = input.id.to_string();

    // do this better
    // for container in containers {
    //     if container.id() == &id {
    //         let stats = container.statistics();
    //         let summary = container_stats_to_summary(stats);
    //         let output = StreamContainerStatisticsOutput { stats };
    //         return Ok(output)
    //     }
    // }

    Err(error::StreamContainerStatisticsError::ResourceNotFoundException(ResourceNotFoundException { message: format!("Container not found") }))

    // let container = match containers.get(id) {
    //     Some(c) => {
    //         let stats = c.stats();
    //         let output = StreamContainerStatisticsOutput { stats };
    //         return Ok(output)
    //     },
    //     None => return Err(error::StreamContainerStatisticsError::ResourceNotFoundException(ResourceNotFoundException { message: format!("Container not found") })),
    // };
}

pub async fn stream_container_logs(input: StreamContainerLogsInput, state: Extension<Arc<State>>) -> Result<StreamContainerLogsOutput, error::StreamContainerLogsError> {
    let ctl = state.controller.lock().await;
    let containers = ctl.containers();

    // let id = input.id.to_string();

    let output_stream = stream! {
        let mut i = 0;
        loop {
            yield Ok(Logs::Line(LogLine { message: Some("line".to_string()), timestamp: None }));

            if i == 10 {
                break;
            }
            i += 1;
        }
    };
    Ok(StreamContainerLogsOutput::builder()
        .logs(output_stream.into())
        .build()
        .unwrap())

}

// STREAM TEST CODE BELOW

// pub async fn stream_container_logs(input: StreamContainerLogsInput, state: Extension<Arc<State>>) -> Result<StreamContainerLogsOutput, error::StreamContainerLogsError> {
    
//     let (tx, mut rx) = mpsc::channel(16); // Adjust the channel capacity as per your requirements    // async move {
//     task::spawn(async move {
//         loop {
//             tx.send(Bytes::from_static(b"hello world!")).await.unwrap();
//             sleep(Duration::from_secs(1)).await;
//         }
//     });

//     let stream = async_stream::stream! {
//         while let Some(item) = rx.recv().await {
//             yield Ok::<_, hyper::Error>(item);
//         }
//     };

//     Ok(StreamContainerLogsOutput {
//         data: ByteStream::new(SdkBody::from(
//             Body::wrap_stream(stream),
//         )),
//     })

// }

// pub async fn stream_container_log_events(input: StreamContainerLogEventsInput, state: Extension<Arc<State>>) -> Result<StreamContainerLogEventsOutput, error::StreamContainerLogEventsError> {

//     let is_true = true;

//     let output_stream = stream! {
//         let mut i = 0;
//         loop {
//             debug!("sending line");
//             // sleep(Duration::from_secs(1)).await;

//             yield Ok(Logs::Line(LogLine { line: Some("line".to_string()) }));

//             // match is_true {
//             //     true => {
//             //         // yield Result::<Logs, MessageStreamError>::Ok(Logs::Line(LogLine { line: Some("line".to_string()) }));
//             //     },
//             //     false => {
//             //         break;
//             //     }
//             // }
//             // Result::<Logs, MessageStreamError>::Ok(Logs::Line(LogLine { line: Some("line".to_string()) }));

//             if i == 10 {
//                 break;
//             }
//             i += 1;
//         }
//     };
//     Ok(StreamContainerLogEventsOutput::builder()
//         .logs(output_stream.into())
//         .build()
//         .unwrap())
//     // Err(error::StreamContainerLogEventsError::ResourceNotFoundException(error::ResourceNotFoundException { message: format!("no events") }))
// }