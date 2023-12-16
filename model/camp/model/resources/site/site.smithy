$version: "2.0"

namespace awlsring.camp.api

use awlsring.camp.common#Tags

@documentation("The machine's identifier.")
@length(min: 38, max: 38)
@pattern("^s-[a-zA-Z0-9-\b]{38}$")
string SiteIdentifier

structure SiteSummary {
    @required
    identifier: SiteIdentifier

    @required
    platform: String

    @required
    added: Timestamp

    @required
    updated: Timestamp

    @required
    tags: Tags
}