use std::fs;


#[derive(Debug)]
/// Represents the kind of disk (HDD, SSD, NVME, etc.)
pub enum DiskKind {
    HDD,
    SSD,
    NVME,
    Unknown(String),
}

#[derive(Debug)]
/// Represents the interface of a disk (SATA, SCSI, etc.)
pub enum DiskInterface {
    SATA,
    SCSI,
    PCI_E,
    Unknown(String),
}

#[derive(Debug)]
/// Represents a physical disk
pub struct Disk {
    device: String,
    /// The model of the disk
    model: String,
    /// The vendor of the disk
    vendor: String,
    /// The interface of the disk
    interface: DiskInterface,
    /// The serial number of the disk
    serial: String,
    /// The kind of disk
    kind: DiskKind,
    /// The sector size of the disk
    sector_size: i16,
    /// The raw size of the disk
    size_raw: i64,
    /// The actual size of the disk
    size_actual: i64,
}

impl Disk {
    /// Returns the device name of the disk
    pub fn get_device(&self) -> &String {
        &self.device
    }

    /// Returns the model of the disk
    pub fn get_model(&self) -> &String {
        &self.model
    }

    /// Returns the vendor of the disk
    pub fn get_vendor(&self) -> &String {
        &self.vendor
    }

    /// Returns the interface of the disk
    pub fn get_interface(&self) -> &DiskInterface {
        &self.interface
    }

    /// Returns the serial number of the disk
    pub fn get_serial(&self) -> &String {
        &self.serial
    }

    /// Returns the kind of disk
    pub fn get_kind(&self) -> &DiskKind {
        &self.kind
    }

    /// Returns the sector size of the disk
    pub fn get_sector_size(&self) -> &i16 {
        &self.sector_size
    }

    /// Returns the raw size of the disk in bytes
    pub fn get_size_raw(&self) -> &i64 {
        &self.size_raw
    }

    /// Returns the actual size of the disk in bytes (sector size * raw size)
    pub fn get_size_actual(&self) -> &i64 {
        &self.size_actual
    }
}

/// Loads all phyiscal disks on the system and returns them as a vector
pub fn load_disks() -> Vec<Disk> {
    let mut disks = Vec::new();

    let ls_dir =  fs::read_dir("/sys/block");
    let hits = match ls_dir {
        Ok(_) => ls_dir.unwrap(),
        Err(_) => { 
            return disks;
        },
    };

    for dir in hits {
        match dir {
            Ok(dir) => {
                let f_type = dir.file_type().unwrap();
                if (f_type.is_dir() || f_type.is_symlink()) && is_device(&dir) {
                    let disk = form_disk(&dir);
                    disks.push(disk);
                }
            },
            Err(_) => {
                continue;
            },
        }
    }

    disks

}

fn is_device(dir: &fs::DirEntry) -> bool {
    let device_dir = dir.path().join("device");
    let device_dir = fs::read_dir(device_dir);
    device_dir.is_ok()
}

fn form_disk(dir: &fs::DirEntry) -> Disk {
    let device = get_device(dir);
    let model = get_model(dir);
    let vendor = get_vendor(dir);
    let serial = get_serial(dir);
    let interface = get_interface(dir);
    let sector_size = get_sector_size(dir);
    let size_raw = get_size_raw(dir);
    let size_actual = determine_actual_size(sector_size, size_raw);
    let kind = determine_kind(dir);

    Disk {
        device,
        model,
        vendor,
        serial,
        interface,
        sector_size,
        size_raw,
        size_actual,
        kind,
    }
}

fn get_device(dir: &fs::DirEntry) -> String {
    let file_name = dir.file_name();
    let file_name = file_name.to_str();
    match file_name {
        Some(file_name) => file_name.to_string(),
        None => String::from(""),
    }
}

fn get_model(dir: &fs::DirEntry) -> String {
    let model_file = dir.path().join("device").join("model");
    let model_file = fs::read_to_string(model_file);
    match model_file {
        Ok(model) => model.trim().to_string(),
        Err(_) => String::from(""),
    }
}

fn get_vendor(dir: &fs::DirEntry) -> String {
    let vendor_file = dir.path().join("device").join("vendor");
    let vendor_file = fs::read_to_string(vendor_file);
    match vendor_file {
        Ok(vendor) => vendor.trim().to_string(),
        Err(_) => String::from("Unknown"),
    }
}

fn get_sector_size(dir: &fs::DirEntry) -> i16 {
    let sector_size_file = dir.path().join("queue").join("logical_block_size");
    let sector_size_file = fs::read_to_string(sector_size_file);
    match sector_size_file {
        Ok(sector_size) => sector_size.trim().parse::<i16>().unwrap(),
        Err(_) => 0,
    }
}

fn get_size_raw(dir: &fs::DirEntry) -> i64 {
    let size_file = dir.path().join("size");
    let size_file = fs::read_to_string(size_file);
    match size_file {
        Ok(size) => size.trim().parse::<i64>().unwrap(),
        Err(_) => 0,
    }
}

fn get_serial(dir: &fs::DirEntry) -> String {
    let serial_file = dir.path().join("device").join("serial");
    let serial_file = fs::read_to_string(serial_file);
    match serial_file {
        Ok(serial) => serial.trim().to_string(),
        Err(_) => String::from("???"),
    }
}

fn get_interface(dir: &fs::DirEntry) -> DiskInterface {
    if dir.file_name().to_str().unwrap_or_else(|| "").contains("nvme") {
        return DiskInterface::PCI_E;
    }
    let interface_file = dir.path().join("device").join("type");
    let interface_file = fs::read_to_string(interface_file);
    match interface_file {
        Ok(interface) => {
            if interface.contains('0') {
                DiskInterface::SATA
            } else if interface.contains('1') {
                return DiskInterface::SCSI;
            } else {
                return DiskInterface::Unknown(String::from(interface.trim()));
            }
        },
        Err(_) => DiskInterface::Unknown(String::from("Unknown")),
    }
}

fn determine_actual_size(sector: i16, raw: i64) -> i64 {
    let sector_size = sector as i64;
    let raw_size = raw;
    raw_size * sector_size
}

fn determine_kind(dir: &fs::DirEntry) -> DiskKind {
    if dir.file_name().to_str().unwrap().contains("nvme") {
        return DiskKind::NVME;
    }

    let rotational_file = dir.path().join("queue").join("rotational");
    let rotational_file = fs::read_to_string(rotational_file);
    match rotational_file {
        Ok(rotational) => {
            if rotational.contains('1') {
                DiskKind::HDD
            } else {
                DiskKind::SSD
            }
        },
        Err(_) => DiskKind::Unknown(String::from("Other")),
    }
}

#[cfg(test)]
mod tests {
    // only works on linux
    use super::load_disks;

    #[test]
    fn disk() {
        let disks = load_disks();
        assert!(!disks.is_empty());
        println!("{:?}", disks)
    }
}