use bollard::{Docker, API_DEFAULT_VERSION, container::StatsOptions};

use futures_util::stream::StreamExt;
mod common;
mod docker;
mod containers;

pub use containers::Containers;
pub use common::Container;
pub use common::ContainerStatistics;
pub use common::ContainerState;
pub use common::Port;
pub use common::Volume;
pub use common::ContainerLogLine;
pub use common::ContainerProtocol;
pub use common::Network;
pub use common::ContainerType;


#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use bollard::container::{ListContainersOptions, LogsOptions};

    use super::*;

    #[tokio::test]
    async fn containers() {
        let co = containers::Containers::new().unwrap();
        let containers = co.list_containers().await;
        for c in containers {
            println!("{:?}", c);
        }
    }

    #[tokio::test]
    async fn stat_stream() {
        let co = containers::Containers::new().unwrap();

        let cons = co.list_containers().await;

        match cons.first() {
            Some(c) => {
                println!("Container: {:?}", c);
                let mut s = co.stream_container_stats(&c.id).await;
                loop {
                    match s.next().await {
                        Some(Ok(stat)) => {
                            println!("{:?}", stat);
                        },
                        Some(Err(e)) => {
                            println!("Error: {}", e);
                        },
                        None => {
                            println!("Stream ended");
                            break;
                        },
                    }
                }
            },
            None => {
                println!("No containers found");
                return;
            },
        }
    }

    #[tokio::test]
    async fn log_stream() {
        let co = containers::Containers::new().unwrap();

        let cons = co.list_containers().await;

        match cons.first() {
            Some(c) => {
                println!("Container: {:?}", c);
                let mut s = co.stream_container_logs(&c.id, false, None).await;

                loop {
                    match s.next().await {
                        Some(Ok(stat)) => {
                            println!("{:?}", stat);
                        },
                        Some(Err(e)) => {
                            println!("Error: {}", e);
                        },
                        None => {
                            println!("Stream ended");
                            break;
                        },
                    }
                }
            },
            None => {
                println!("No containers found");
                return;
            },
        }
    }
}
