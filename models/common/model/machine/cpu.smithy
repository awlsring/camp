$version: "2.0"

namespace awlsring.camp.common.machine

enum CpuArchitecture {
    X86 = "x86"
    ARM = "arm"
    UNKNOWN = "Unknown"
}

structure MachineCpuSummary {
    @required
    cores: Integer

    @required
    architecture: CpuArchitecture

    model: String

    vendor: String
}
