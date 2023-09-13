use std::sync::Arc;

use aws_smithy_http_server::Extension;
use camp_agent_server::{output::GetMemoryOutput, model::{MemorySummary, MemoryTypeSummary}, input::GetMemoryInput, error};

use crate::{server::http::State, stats::memory::Memory};


pub async fn get_memory(_input: GetMemoryInput, state: Extension<Arc<State>>) -> Result<GetMemoryOutput, error::GetMemoryError> {
    let ctl = state.controller.lock().await;
    let mem = ctl.memory();

    let sum = memory_to_summary(mem);

    let output = GetMemoryOutput {
        summary: sum
    };

    Ok(output)
}

pub fn memory_to_summary(mem: &Memory) -> MemorySummary {
    MemorySummary {
        memory: MemoryTypeSummary {
            total: *mem.memory().total() as i64,
            available: *mem.memory().available() as i64,
            used: *mem.memory().used() as i64,
        },
        swap: MemoryTypeSummary {
            total: *mem.swap().total() as i64,
            available: *mem.swap().available() as i64,
            used: *mem.swap().used() as i64,
        }
        
    }
}