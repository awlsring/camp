use bollard::{Docker, container::{StatsOptions, LogOutput, LogsOptions, Stats}, errors::Error};
use futures_util::{StreamExt, Stream};
use log::{warn, debug};

use crate::common::{Container, ContainerStatistics, ContainerLogLine};

pub struct Containers {
    docker: Docker,
}

impl Containers {
    pub fn new() -> Option<Containers> {
        let c = Docker::connect_with_local_defaults();
        match c {
            Ok(c) => Some(Containers {
                docker: c,
            }),
            Err(e) => {
                warn!("Error connecting to docker: {}", e);
                None
            }
        }
    }

    pub async fn list_containers(&self) -> Vec<Container> {
        let mut result: Vec<Container> = Vec::new();
        let containers = self.docker.list_containers::<String>(None).await;

        match containers {
            Ok(containers) => {
                debug!("Found {} containers", containers.len());
                for c in containers {
                    match c.id {
                        Some(ref id) => {
                            let c = self.get_container(id).await;
                            match c {
                                Some(c) => result.push(c),
                                None => {
                                    debug!("Found container with no data");
                                },
                            }
                        },
                        None => { 
                            debug!("Found container with no id");
                        },
                    }
                }
            },
            Err(e) => warn!("Error listing containers: {}", e),
        }
        result
    }

    pub async fn get_container(&self, id: &str) -> Option<Container> {
        let result = self.docker.inspect_container(id, None).await;
        match result {
            Ok(container) => {
                let stats = self.get_container_stats(id).await;
                let container = Container::new_from_docker(container, stats);
                Some(container)
            },
            Err(e) => {
                warn!("Error getting container {}: {}", id, e);
                None
            }
        }
    }

    pub async fn stream_container_logs(&self, id: &str, follow: bool, lines: Option<i32>) -> impl Stream<Item = Result<ContainerLogLine, Error>> {
        let mut options = Some(LogsOptions::<String>{
            follow,
            stdout: true,
            timestamps: true,
            ..Default::default()
        });
        if let Some(lines) = lines {
            options.as_mut().unwrap().tail = lines.to_string();
        }

        let logs = self.docker.logs(id, options);
        let l = logs.map(|l| {
            l.map(|op| {
                let log = ContainerLogLine::new_from_docker(op);
                log
            })
        });
        l
    }

    pub async fn stream_container_stats(&self, id: &str) -> impl Stream<Item = Result<ContainerStatistics, Error>> {
        let options = Some(StatsOptions {
            stream: true,
            ..Default::default()
        });
        let stats = self.docker.stats(id, options);
        let statistics = stats.map(|s| {
            s.map(|op| {
                let stats = ContainerStatistics::new_from_docker(op);
                stats
            })
        });
        statistics
    }

    pub async fn get_container_stats(&self, id: &str) -> Option<ContainerStatistics> {
        let options = Some(StatsOptions {
            stream: false,
            one_shot: true,
            ..Default::default()
        });
        let mut result = self.docker.stats(id, options);
        match result.next().await {
            Some(Ok(stats)) => {
                let stats = ContainerStatistics::new_from_docker(stats);
                Some(stats)
            },
            Some(Err(e)) => {
                warn!("Error getting stats for container {}: {}", id, e);
                None
            },
            None => {
                warn!("No stats for container {}", id);
                None
            }
        }
    }

}