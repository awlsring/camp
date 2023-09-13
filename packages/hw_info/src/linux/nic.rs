use std::fs;


#[derive(Debug)]
/// Represents a network interface
pub struct NetworkInterface {
    /// The name of the interface
    name: String,
    /// The vendor of the interface
    vendor: String,
    /// The vendor code of the interface
    vendor_code: String,
    /// The MAC address of the interface
    mac_address: String,
    /// The speed of the interface
    speed: u16,
    /// The duplex of the interface
    duplex: String,
    /// The MTU of the interface
    mtu: u16,
    /// The broadcast of the interface
    broadcast: String,
    /// Whether the interface is active
    active: bool,
    /// The device ID of the interface
    device_id: String,
    /// The subsystem ID of the interface
    subsystem_id: String,
}

impl NetworkInterface {
    pub fn name(&self) -> &String {
        &self.name
    }

    pub fn vendor(&self) -> &String {
        &self.vendor
    }

    pub fn vendor_code(&self) -> &String {
        &self.vendor_code
    }

    pub fn mac_address(&self) -> &String {
        &self.mac_address
    }

    pub fn speed(&self) -> &u16 {
        &self.speed
    }

    pub fn duplex(&self) -> &String {
        &self.duplex
    }

    pub fn mtu(&self) -> &u16 {
        &self.mtu
    }

    pub fn broadcast(&self) -> &String {
        &self.broadcast
    }

    pub fn active(&self) -> &bool {
        &self.active
    }

    pub fn device_id(&self) -> &String {
        &self.device_id
    }

    pub fn subsystem_id(&self) -> &String {
        &self.subsystem_id
    }

}

pub fn load_nics() -> Vec<NetworkInterface> {
    let mut nics = Vec::new();

    let paths = fs::read_dir("/sys/class/net").unwrap();
    for path in paths {

        match path {
            Ok(_) => {
                let path = path.unwrap();
                if is_physical(&path) {
                    let nic = form_nic(&path);
                    nics.push(nic);
                }
            },
            Err(_) => {},
        }

    }

    nics
}

fn form_nic(dir: &fs::DirEntry) -> NetworkInterface {
    let name = get_name(dir);
    let vendor_code = get_vendor_code(dir);
    let vendor = vendor_code_to_vendor(&vendor_code);
    let mac_address = get_mac_address(dir);
    let speed = get_speed(dir);
    let duplex = get_duplex(dir);
    let mtu = get_mtu(dir);
    let broadcast = get_broadcast(dir);
    let active = is_active(dir);
    let device_id = get_device_id(dir);
    let subsystem_id = get_subsystem_id(dir);

    NetworkInterface {
        name,
        vendor_code,
        vendor,
        mac_address,
        speed,
        duplex,
        mtu,
        broadcast,
        active,
        device_id,
        subsystem_id,
    }
}

fn get_name(dir: &fs::DirEntry) -> String {
    let binding = dir.file_name();
    let name = binding.to_str();
    match name {
        Some(name) => name.to_string(),
        None => "".to_string(),
    }
}

fn get_vendor_code(dir: &fs::DirEntry) -> String {
    let vendor_code = fs::read_to_string(dir.path().join("device").join("vendor"));
    match vendor_code {
        Ok(vendor_code) => vendor_code.trim().to_owned(),
        Err(_) => "".to_string(),
    }
}

fn get_mac_address(dir: &fs::DirEntry) -> String {
    let mac_address = fs::read_to_string(dir.path().join("address"));
    match mac_address {
        Ok(mac_address) => mac_address.trim().to_owned(),
        Err(_) => "".to_string(),
    }
}

fn get_speed(dir: &fs::DirEntry) -> u16 {
    let speed = fs::read_to_string(dir.path().join("speed"));
    match speed {
        Ok(speed) => speed.trim().parse::<u16>().unwrap(),
        Err(_) => 0,
    }
}

fn is_active(dir: &fs::DirEntry) -> bool {
    let operstate = fs::read_to_string(dir.path().join("operstate"));

    let operstate = match operstate {
        Ok(operstate) => operstate,
        Err(_) => return false,
    };
    if operstate == "up" {
        return true;
    }
    false
}

fn get_duplex(dir: &fs::DirEntry) -> String {
    let duplex = fs::read_to_string(dir.path().join("duplex"));
    match duplex {
        Ok(duplex) => duplex.trim().to_string(),
        Err(_) => "".to_string(),
    }
}

fn get_broadcast(dir: &fs::DirEntry) -> String {
    let broadcast = fs::read_to_string(dir.path().join("broadcast"));
    match broadcast {
        Ok(broadcast) => broadcast.trim().to_string(),
        Err(_) => "".to_string(),
    }
}

fn get_device_id(dir: &fs::DirEntry) -> String {
    let device_id = fs::read_to_string(dir.path().join("device").join("device"));
    match device_id {
        Ok(device_id) => device_id.trim().to_string(),
        Err(_) => "".to_string(),
    }
}

fn get_subsystem_id(dir: &fs::DirEntry) -> String {
    let subsystem_id = fs::read_to_string(dir.path().join("device").join("subsystem_device"));
    let subsystem_id = match subsystem_id {
        Ok(subsystem_id) => subsystem_id.trim().to_string(),
        Err(_) => "".to_string(),
    };
    subsystem_id
}

fn get_mtu(dir: &fs::DirEntry) -> u16 {
    let mtu = fs::read_to_string(dir.path().join("mtu"));
    match mtu {
        Ok(mtu) => mtu.trim().parse::<u16>().unwrap(),
        Err(_) => 0,
    }
}

fn is_physical(dir: &fs::DirEntry) -> bool {
    let driver_dir = dir.path().join("device").join("driver");
    let driver_dir = fs::read_dir(driver_dir);
    match driver_dir {
        Ok(_) => {
            if dir.file_type().unwrap().is_symlink() {
                return true;
            }
            false
        },
        Err(_) => false,
    }
}

fn vendor_code_to_vendor(code: &str) -> String {
    match code {
        "0x8086"=>"Intel Corporation".to_string(),
        "0x8087"=>"Intel Corporation".to_string(),
        "0x14e4"=>"Broadcom Inc. and subsidiaries".to_string(),
        "0x10ec"=>"Realtek Semiconductor Co., Ltd.".to_string(),
        "0x1af4"=>"Red Hat, Inc.".to_string(),
        "0x168C"=>"Qualcomm Atheros".to_string(),
        "0x15b3"=>" Mellanox Technologies".to_string(),
        _=>"Unknown".to_string(),
    }
}

#[cfg(test)]
mod tests {
    // only works on linux
    use super::load_nics;

    #[test]
    fn nics() {
        let nics = load_nics();
        assert!(!nics.is_empty());
        println!("{:?}", nics)
    }
}