$version: "2.0"

namespace awlsring.camp.common.tags

@pattern("^[a-zA-Z0-9_]+( [a-zA-Z0-9_]+){0,127}$")
string TagString

structure Tag {
    @required
    key: TagString

    @required
    value: TagString
}

list Tags {
    member: Tag
}
