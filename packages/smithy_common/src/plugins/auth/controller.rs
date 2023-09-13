use http::HeaderValue;
use log::debug;

#[derive(Clone, Debug)]
pub struct AuthController {
    allowed_keys: Vec<String>,
    no_auth_operations: Vec<String>,
}

impl AuthController {
    pub fn new(no_auth_operations: &Vec<String>, allowed_keys: &Vec<String>) -> Self {
        AuthController {
            allowed_keys: allowed_keys.clone(),
            no_auth_operations: no_auth_operations.clone(),
        }
    }

    pub async fn auth(&self, operation: &str, key: Option<&HeaderValue>) -> bool {
        debug!("Determining auth for operation: {}", operation);
        if self.no_auth_operations.contains(&operation.to_string()) {
            debug!("Operation {} does not require auth", operation);
            return true;
        }

        if let Some(auth_header) = key {
            let key = extract_api_key(auth_header);

            match key {
                Some(key) => {
                    if self.allowed_keys.contains(&key.to_string()) {
                        debug!("Key is in allowlist");
                        return true;
                    }
                },
                None => {
                    debug!("No key found in header");
                },
            }
        }

        debug!("authorization denied");
        false
    }
}

fn extract_api_key(header: &HeaderValue) -> Option<String> {
    let key = header.to_str().unwrap_or("none");
    match key {
        "none" => None,
        _ => {
            let parts = key.split(" ").collect::<Vec<&str>>();
            match parts.len() {
                2 => Some(parts[1].to_string()),
                _ => None,
            }
        },
    }
}