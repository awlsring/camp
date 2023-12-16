$version: "2.0"

namespace awlsring.camp.common

structure CoreSummary {
    @required
    identifier: Integer

    @required
    threadCount: Integer
}

list CoreList {
    member: CoreSummary
}
