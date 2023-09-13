pub fn handle_optional_string(optional: Option<String>) -> String {
    match optional {
        Some(s) => s,
        None => "Unknown".to_string()
    }
}

pub fn handle_optional_str(optional: Option<&str>) -> String {
    match optional {
        Some(s) => String::from(s),
        None => "Unknown".to_string()
    }
}

pub fn u8_as_string(u: &[u8]) -> String {
    let str = std::str::from_utf8(u);
    match str {
        Ok(str) => str.to_string().to_uppercase(),
        Err(_e) => "Unknown".to_string()
    }
}

pub fn handle_optional_usize(optional: Option<usize>) -> usize {
    optional.unwrap_or(0)
}