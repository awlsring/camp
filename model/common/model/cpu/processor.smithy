$version: "2.0"

namespace awlsring.camp.common

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
