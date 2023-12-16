$version: "2.0"

namespace awlsring.camp.common

@pattern("^[a-zA-Z0-9_]+( [a-zA-Z0-9_]+){0,127}$")
@length(min: 1, max: 50)
string TagKey

@pattern("^[a-zA-Z0-9_]+( [a-zA-Z0-9_]+){0,127}$")
@length(min: 1, max: 128)
string TagValue

structure Tag {
    @required
    key: TagKey

    @required
    value: TagValue
}

list Tags {
    member: Tag
}
