mod linux;

pub use linux::disk::Disk;
pub use linux::disk::DiskKind;
pub use linux::disk::DiskInterface;
pub use linux::disk::load_disks;

pub use linux::nic::NetworkInterface;
pub use linux::nic::load_nics;

pub fn add(left: usize, right: usize) -> usize {
    left + right
}

#[cfg(test)]
mod tests {
    use crate::linux::disk::load_disks;

    use super::*;

    #[test]
    fn it_works() {
        let result = add(2, 2);
        assert_eq!(result, 4);
    }

    #[test]
    fn disk() {
        load_disks();
    }
}
