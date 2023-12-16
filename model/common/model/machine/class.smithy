$version: "2.0"

namespace awlsring.camp.common

@documentation("The class of a machine.")
enum MachineClass {
    BARE_METAL = "BareMetal"
    VIRTUAL_MACHINE = "VirtualMachine"
    HYPERVISOR = "Hypervisor"
    UNKNOWN = "Unknown"
}
