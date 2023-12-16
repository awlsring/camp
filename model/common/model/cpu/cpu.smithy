$version: "2.0"

namespace awlsring.camp.common

structure CpuSummary {
    @required
    totalCores: Integer

    @required
    totalThreads: Integer

    @required
    architecture: Architecture

    processors: ProcessorList

    vendor: String

    model: String
}

structure ProcessorSummary {
    @required
    identifier: Integer

    @required
    coreCount: Integer

    @required
    threadCount: Integer

    vendor: String

    model: String

    cores: CoreList
}

list ProcessorList {
    member: ProcessorSummary
}

structure CoreSummary {
    @required
    identifier: Integer

    @required
    threadCount: Integer
}

list CoreList {
    member: CoreSummary
}
