use std::collections::HashMap;
use std::net::IpAddr;

use log::info;
use hw_info::load_nics;
use sysinfo::{System, SystemExt, NetworkData, NetworkExt};

use network_interface::Addr;
use network_interface::NetworkInterface as NetIface;
use network_interface::NetworkInterfaceConfig;
use network_interface::V6IfAddr;

use hw_info::NetworkInterface as PhysNetworkInterface;

use super::util::handle_optional_string;

pub enum AddressKind {
    V4,
    V6,
    V6Local,
}

pub struct Address {
    version: AddressKind,
    address: IpAddr,
    broadcast: Option<IpAddr>,
    netmask: Option<IpAddr>,
}

impl Address {
    pub fn new(addr: &Addr) -> Address {
        let version = match addr {
            Addr::V4(_addr) => AddressKind::V4,
            Addr::V6(addr) => {
                if Address::is_local_ipv6(addr) {
                    AddressKind::V6Local
                } else {
                    AddressKind::V6
                }
            }
        };
        Address {
            version,
            address: addr.ip(),
            broadcast: addr.broadcast(),
            netmask: addr.netmask(),
        }
    }

    fn is_local_ipv6(addr: &V6IfAddr) -> bool {
        addr.ip.to_string().starts_with("fe80")
    }

    pub fn version(&self) -> &AddressKind {
        &self.version
    }

    pub fn address(&self) -> &IpAddr {
        &self.address
    }

    pub fn broadcast(&self) -> &Option<IpAddr> {
        &self.broadcast
    }

    pub fn netmask(&self) -> &Option<IpAddr> {
        &self.netmask
    }
}

pub struct NetworkInterfaceTraffic {
    transmitted: u64,
    recieved: u64
}

impl NetworkInterfaceTraffic {
    pub fn transmitted(&self) -> &u64 {
        &self.transmitted
    }

    pub fn recieved(&self) -> &u64 {
        &self.recieved
    }
}

pub struct NetworkInterface {
    name: String,
    addresses: Vec<Address>,
    bytes: NetworkInterfaceTraffic,
    packets: NetworkInterfaceTraffic,
    r#virtual: bool,
    mac: Option<String>,
    vendor: Option<String>,
    mtu: Option<u16>,
    duplex: Option<String>,
    speed: Option<u16>,
}

impl NetworkInterface {
    pub fn new(iface: &NetIface, phys: Option<&&PhysNetworkInterface>) -> NetworkInterface {
        let name = &iface.name;
        let mac = iface.mac_addr.to_owned();

        let mut addresses = Vec::new();
        for addr in &iface.addr {
            addresses.push(Address::new(addr));
        }

        let r#virtual = phys.is_none();

        let vendor: Option<String> = match phys {
            Some(phys) => Some(phys.vendor().to_owned()),
            None => None,
        };

        let mtu: Option<u16> = match phys {
            Some(phys) => Some(phys.mtu().to_owned()),
            None => None,
        };

        let duplex: Option<String> = match phys {
            Some(phys) => Some(phys.duplex().to_owned()),
            None => None,
        };

        let speed: Option<u16> = match phys {
            Some(phys) => Some(phys.speed().to_owned()),
            None => None,
        };
        
        NetworkInterface {
            name: String::from(name),
            addresses,
            bytes: NetworkInterfaceTraffic {
                transmitted: 0,
                recieved: 0
            },
            packets: NetworkInterfaceTraffic {
                transmitted: 0,
                recieved: 0
            },
            mac,
            vendor,
            mtu,
            duplex,
            speed,
            r#virtual: r#virtual,
        }
    }

    pub fn addresses(&self) -> &Vec<Address> {
        &self.addresses
    }

    pub fn name(&self) -> &String {
        &self.name
    }

    pub fn mac(&self) -> &Option<String> {
        &self.mac
    }

    pub fn bytes(&self) -> &NetworkInterfaceTraffic {
        &self.bytes
    }

    pub fn packets(&self) -> &NetworkInterfaceTraffic {
        &self.packets
    }

    pub fn vendor(&self) -> &Option<String> {
        &self.vendor
    }

    pub fn mtu(&self) -> &Option<u16> {
        &self.mtu
    }

    pub fn duplex(&self) -> &Option<String> {
        &self.duplex
    }

    pub fn speed(&self) -> &Option<u16> {
        &self.speed
    }

    pub fn is_virtual(&self) -> &bool {
        &self.r#virtual
    }

    pub fn update(&mut self, data: &NetworkData) {
        self.bytes.transmitted = data.total_transmitted();
        self.bytes.recieved = data.total_received();
        self.packets.transmitted = data.packets_transmitted();
        self.packets.recieved = data.packets_received();
    }

}

pub struct Network {
    interfaces: HashMap<String, NetworkInterface>
}

impl Network {
    pub fn new(system: &System) -> Network {
        let mut interfaces = HashMap::new();
        let mut phys_interfaces = HashMap::new();
        let network_interfaces = NetIface::show().unwrap_or_default();

        let binding = load_nics();
        for nic in binding.iter() {
            phys_interfaces.insert(nic.name().clone(), nic);
        }
        
        for interface in network_interfaces.iter() {
            let name = &interface.name;
            let phys = phys_interfaces.get(name);

            let network_interface = NetworkInterface::new(interface, phys);

            interfaces.insert(String::from(name), network_interface);
        }        

        let mut net = Network {
            interfaces
        };
        net.update(system);

        net
    }

    pub fn get_network_interface(&self, name: &str) -> Option<&NetworkInterface> {
        self.interfaces.get(name)
    }

    pub fn network_interfaces(&self) -> Vec<&NetworkInterface> {
        let mut nics = Vec::new();

        for nic in self.interfaces.values() {
            nics.push(nic);
        }

        nics
    }

    pub fn update(&mut self, system: &System) {
        for (name, iface_data) in system.networks() {

            let iface = self.interfaces.get_mut(name);

            if let Some(iface) = iface {
                iface.update(iface_data);
            }
        }
    }

}
